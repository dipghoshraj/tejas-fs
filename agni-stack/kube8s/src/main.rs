mod k8s;

use k8s::deploy::deployk8s;


#[tokio::main]
async fn main() {
    println!("Hello, world!");

    let appname = "myapp";
    let image = "dipghoshraj/omnia-ar-app:0.0.7";
    let result = deployk8s(appname, image).await;
    match result {
        Ok(_) => println!("Deployment successful"),
        Err(e) => eprintln!("Error deploying: {}", e),
    }
}
