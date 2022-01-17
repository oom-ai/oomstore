use futures_util::{pin_mut, StreamExt};
use oomstore::{oomagent::EntityRow, Client};

macro_rules! svec { ($($part:expr),* $(,)?) => {{ vec![$($part.to_string(),)*] }} }

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let join_features = svec![
        "driver_stats.conv_rate",
        "driver_stats.acc_rate",
        "driver_stats.avg_daily_trips"
    ];
    let existed_features = svec![];

    let entity_rows = vec![
        EntityRow {
            entity_key: "1".into(),
            unix_milli: 3,
            values:     Vec::new(),
        },
        EntityRow {
            entity_key: "7".into(),
            unix_milli: 1,
            values:     Vec::new(),
        },
        EntityRow {
            entity_key: "7".into(),
            unix_milli: 0,
            values:     Vec::new(),
        },
    ];

    let (header, rows) = client
        .channel_join(join_features, existed_features, entity_rows.into_iter())
        .await?;

    println!("RESPONSE HEADER={:?}", header);

    pin_mut!(rows); // needed for iteration

    while let Some(row) = rows.next().await {
        println!("RESPONSE ROWS={:?}", row);
    }

    Ok(())
}
