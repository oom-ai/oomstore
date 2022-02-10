use nix::{
    sys::signal::{self, Signal},
    unistd::Pid,
};
use signal_hook::{consts::signal::*, low_level::emulate_default_handler};
use signal_hook_tokio::{Handle, Signals};
use std::{
    env,
    net::{IpAddr, SocketAddr},
    path::Path,
};
use tokio::{fs, process::Command, time, time::Duration};
use tokio_stream::StreamExt;

use crate::Result;

#[derive(Debug, Clone)]
pub struct ServerWrapper {
    addr:   SocketAddr,
    handle: Handle,
    pid:    u32,
}

// TODO: Add a Builder to create the server wrapper
impl ServerWrapper {
    pub async fn new<P1, P2>(bin_path: Option<P1>, cfg_path: Option<P2>, port: Option<u16>) -> Result<Self>
    where
        P1: AsRef<Path>,
        P2: AsRef<Path>,
    {
        let bin_path = bin_path.map(|p| p.as_ref().to_owned());
        let cfg_path = cfg_path.map(|p| p.as_ref().to_owned());
        let port = port.unwrap_or(0);

        let mut signals = Signals::new(&[SIGHUP, SIGTERM, SIGINT, SIGQUIT])?;
        let handle = signals.handle();

        let mut oomagent = Command::new(bin_path.clone().unwrap_or_else(|| "oomagent".into()));
        oomagent.arg("--port").arg(port.to_string());
        if let Some(cfg_path) = cfg_path.clone() {
            oomagent.arg("--config").arg(cfg_path);
        }

        let child = oomagent.spawn()?;
        let pid = child.id().expect("failed to get child pid");

        tokio::spawn({
            async move {
                while let Some(signal) = signals.next().await {
                    graceful_kill(pid, signal);
                    emulate_default_handler(signal).expect("failed to emulate default signal handler");
                }
            }
        });

        let addr = get_agent_address(pid).await?;
        Ok(Self { addr, handle, pid })
    }

    pub async fn default() -> Result<Self> {
        Self::new(None::<String>, None::<String>, None).await
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

impl Drop for ServerWrapper {
    fn drop(&mut self) {
        self.handle.close();
        graceful_kill(self.pid, SIGTERM);
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

fn graceful_kill(pid: u32, signal: i32) {
    let signal = match signal {
        SIGHUP => Signal::SIGHUP,
        SIGTERM => Signal::SIGTERM,
        SIGINT => Signal::SIGINT,
        SIGQUIT => Signal::SIGQUIT,
        _ => panic!("unexpected signal: {}", signal),
    };
    signal::kill(Pid::from_raw(pid as i32), signal).expect("failed to send SIGTERM signal to child process");
}
