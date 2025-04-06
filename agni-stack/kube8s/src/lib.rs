mod k8s;

use k8s::deploy::deployk8s;
use k8s::network::create_ingress;
use k8s::service::create_service;


pub async fn deploy_app(
    appname: &str,
    image: &str,
) -> Result<(), Box<dyn std::error::Error + Send + Sync + 'static>> {
    // Deploy the application
    println!("Hello, world!");
    
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
                        Ok(_) => {
                            println!("Ingress created successfully");
                            Ok(())
                        },
                        Err(e) => Err(e as Box<dyn std::error::Error + Send + Sync + 'static>),
                    }
                }
                Err(e) =>Err(e as Box<dyn std::error::Error + Send + Sync + 'static>),
            }
        },
        Err(e) => Err(e as Box<dyn std::error::Error + Send + Sync + 'static>),

    }
}      