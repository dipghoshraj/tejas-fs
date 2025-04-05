use k8s_openapi::api::core::v1::Service;
use kube::{
    api::{Api, PostParams},
    Client,
};
use serde_json::json;

pub async fn create_service(
    appname: &str,
    target_port: i32,
) -> Result<(), Box<dyn std::error::Error + Send + Sync + 'static>> {
    let client = Client::try_default().await?;
    let services: Api<Service> = Api::namespaced(client, "default");

    let service = json!({
        "apiVersion": "v1",
        "kind": "Service",
        "metadata": {
            "name": appname,
        },
        "spec": {
            "selector": { "app": appname },
            "ports": [{
                "protocol": "TCP",
                "port": 80,
                "targetPort": target_port
            }],
            "type": "ClusterIP"
        }
    });

    let service: Service = serde_json::from_value(service)?;
    services.create(&PostParams::default(), &service).await?;

    println!("✅ Successfully created service for {}", appname);
    Ok(())
}