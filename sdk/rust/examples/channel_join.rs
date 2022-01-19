use futures_util::{pin_mut, StreamExt};
use oomrpc::{Client, EntityRow};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://localhost:50051").await?;

    let join_features = vec![
        "driver_stats.conv_rate".into(),
        "driver_stats.acc_rate".into(),
        "driver_stats.avg_daily_trips".into(),
    ];
    let existed_features = vec![];

    let entity_rows = tokio_stream::iter([
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
    ]);

    let (header, rows) = client
        .channel_join(join_features, existed_features, entity_rows)
        .await?;

    println!("RESPONSE HEADER={:?}", header);

    pin_mut!(rows); // needed for iteration

    while let Some(row) = rows.next().await {
        println!("RESPONSE ROWS={:?}", row?);
    }

    Ok(())
}
