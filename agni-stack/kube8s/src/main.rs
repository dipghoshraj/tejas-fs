mod k8s;

use k8s::deploy::deployk8s;
use k8s::network::create_ingress;
use k8s::service::create_service;


#[tokio::main]
async fn main() {
    println!("Hello, world!");

    let appname = "myapp";
    let image = "dipghoshraj/omnia-ar-app:0.0.7";
    
    let result = deployk8s(appname, image).await;
    match result {
        Ok(_) => {
            println!("Deployment successful");
            let result = create_service(appname, 80).await;
            match result {
                Ok(_) => {
                    println!("Service created successfully");
                    let result = create_ingress(appname, "myapp.example.com").await;
                    match result {
                        Ok(_) => println!("Ingress created successfully"),
                        Err(e) => eprintln!("Error creating ingress: {}", e),
                    }
                }
                Err(e) => eprintln!("Error creating service: {}", e),
            }
        },
        Err(e) => eprintln!("Error deploying: {}", e),
    }
}
