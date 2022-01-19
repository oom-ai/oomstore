use oomrpc::{Client, Value};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = Client::connect("http://127.0.0.1:50051").await?;

    let kvs = vec![
        ("last_5_click_posts", Value::String("hello".into())),
        ("number_of_user_starred_posts", Value::Int64(28)),
    ];

    client.push("1", "user-click", kvs).await?;

    Ok(())
}
