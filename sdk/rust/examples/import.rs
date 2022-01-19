use std::fs;

use oomrpc::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://localhost:50051").await?;

    let contents = r#"
user,state,credit_score,account_age_days,has_2fa_installed
1,Nevada,530,242,true
2,South Carolina,520,268,false
3,New Jersey,655,84,false
4,Ohio,677,119,true
5,California,566,289,false
6,North Carolina,533,155,true
7,North Dakota,605,334,true
8,West Virginia,664,282,false
9,Alabama,577,150,true
10,Idaho,693,212,true
"#
    .trim_start();

    let path = "/tmp/demo-account.csv";
    fs::write(path, contents)?;
    client.import("account", None, None, path, None).await?;

    Ok(())
}
