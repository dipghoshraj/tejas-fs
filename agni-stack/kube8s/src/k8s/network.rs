use k8s_openapi::api::networking::v1::Ingress;
use kube::{
    api::{Api, PostParams},
    Client,
};
use serde_json::json;


pub async fn create_ingress(
    appname: &str,
    host: &str,
) -> Result<(), Box<dyn std::error::Error + Send + Sync + 'static>> {
    let client = Client::try_default().await.unwrap();
    let ingresses: Api<Ingress> = Api::namespaced(client, "default");

    let ingress = json!({
        "apiVersion": "networking.k8s.io/v1",
        "kind": "Ingress",
        "metadata": {
            "name": appname,
        },
        "spec": {
            "rules": [{
                "host": host,
                "http": {
                    "paths": [{
                        "path": "/",
                        "pathType": "Prefix",
                        "backend": {
                            "service": {
                                "name": appname,
                                "port": {
                                    "number": 80
                                }
                            }
                        }
                    }]
                }
            }]
        }
    });

    let ingress: Ingress = serde_json::from_value(ingress)?;
    ingresses.create(&PostParams::default(), &ingress).await?;

    println!("✅ Successfully created ingress for {}", appname);
    Ok(())
}