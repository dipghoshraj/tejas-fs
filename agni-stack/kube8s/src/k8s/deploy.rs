

use kube::{
    api::{Api, PostParams},
    Client
};

use k8s_openapi::api::apps::v1::Deployment;
use serde_json::json;

pub async fn deployk8s (appname: &str, image: &str) -> Result<(), Box<dyn std::error::Error>> {
    let client = Client::try_default().await.unwrap();
    let deployments: Api<Deployment> = Api::namespaced(client, "default");

    let deployment = json!({
        "apiVersion": "apps/v1",
        "kind": "Deployment",
        "metadata": {
            "name": appname,
        },
        "spec": {
            "replicas": 1,
            "selector": {
                "matchLabels": { "app": appname }
            },
            "template": {
                "metadata": {
                    "labels": { "app": appname }
                },
                "spec": {
                    "containers": [{
                        "name": appname,
                        "image": image,
                        "ports": [{"containerPort": 80}]
                    }]
                }
            }
        }
    });

    let deploy: Deployment = serde_json::from_value(deployment)?;
    deployments.create(&PostParams::default(), &deploy).await?;

    println!("✅ Successfully deployed {}", appname);
    Ok(())
}