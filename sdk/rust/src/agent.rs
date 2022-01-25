use signal_hook::{consts::signal::*, low_level::emulate_default_handler};
use signal_hook_tokio::{Handle, Signals};
use std::{
    env,
    net::{IpAddr, SocketAddr},
    path::Path,
    process::{Child, Command},
    sync::{Arc, Mutex},
};
use tokio::{fs, time, time::Duration};
use tokio_stream::StreamExt;

use crate::Result;

pub struct EmbeddedAgent {
    addr:   SocketAddr,
    handle: Handle,
    child:  Arc<Mutex<Child>>,
}

impl EmbeddedAgent {
    pub async fn new<P: AsRef<Path>>(
        cmd_path: impl Into<Option<P>>,
        cfg_path: impl Into<Option<P>>,
        port: impl Into<Option<u16>>,
    ) -> Result<Self> {
        let cmd_path = cmd_path.into().map(|p| p.as_ref().to_owned());
        let cfg_path = cfg_path.into().map(|p| p.as_ref().to_owned());
        let port = port.into().unwrap_or(0);

        let mut signals = Signals::new(&[SIGHUP, SIGTERM, SIGINT, SIGQUIT])?;
        let handle = signals.handle();

        let mut oomagent = Command::new(cmd_path.clone().unwrap_or_else(|| "oomagent".into()));
        oomagent.arg("--port").arg(port.to_string());
        if let Some(cfg_path) = cfg_path.clone() {
            oomagent.arg("--config").arg(cfg_path);
        }

        let child = oomagent.spawn()?;
        let child = Arc::new(Mutex::new(child));

        tokio::spawn({
            let child = Arc::clone(&child);
            async move {
                if let Some(signal) = signals.next().await {
                    child.lock().unwrap().kill().unwrap();
                    emulate_default_handler(signal).unwrap();
                }
            }
        });

        let pid = child.lock().unwrap().id();
        let addr = get_agent_address(pid).await?;
        Ok(Self { handle, child, addr })
    }

    pub fn ip(&self) -> IpAddr {
        self.addr.ip()
    }

    pub fn port(&self) -> u16 {
        self.addr.port()
    }

    pub fn address(&self) -> SocketAddr {
        self.addr
    }
}

impl Drop for EmbeddedAgent {
    fn drop(&mut self) {
        self.handle.close();
        self.child.lock().unwrap().kill().unwrap();
    }
}

async fn get_agent_address(pid: u32) -> Result<SocketAddr> {
    let mut path = env::temp_dir();
    path.push("oomagent");
    path.push(pid.to_string());
    path.push("address");
    let time = time::Instant::now();

    loop {
        let result = fs::read_to_string(&path)
            .await
            .map_err(|e| e.into())
            .and_then(|addr| Ok(addr.parse()?));
        if result.is_ok() || time.elapsed() > Duration::from_millis(3000) {
            return result;
        }
        time::sleep(Duration::from_millis(200)).await;
    }
}
