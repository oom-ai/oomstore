use std::fs;

use oomrpc::Client;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let contents = r#"
entity_key,unix_milli,age
1,0,20
1,3,44
1,4,37
2,3,28
3,3,48
4,4,22
5,2,41
5,2,43
5,3,47
5,4,59
5,5,55
6,1,27
6,3,46
7,0,42
7,4,40
7,4,49
8,4,39
9,3,35
10,4,54
10,4,33
"#
    .trim_start();

    let input = "/tmp/driver_stats_label.csv";
    let output = "/tmp/joined.csv";
    fs::write(input, contents)?;

    let features = vec![
        "driver_stats.conv_rate".into(),
        "driver_stats.acc_rate".into(),
        "driver_stats.avg_daily_trips".into(),
    ];

    client.join(features, input, output).await?;

    fs::copy(output, "/dev/stdout")?;

    Ok(())
}
