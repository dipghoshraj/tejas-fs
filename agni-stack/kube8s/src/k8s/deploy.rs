

use kube::{
    api::{Api, Patch, PatchParams},
    Client
};

use k8s_openapi::api::apps::v1::Deployment;
use serde_json::json;

pub async fn deployk8s (appname: &str, image: &str) -> Result<(), Box<dyn std::error::Error + Send + Sync + 'static>> {
    let client = Client::try_default().await.unwrap();
    let deployments: Api<Deployment> = Api::namespaced(client, "default");
    let controller = format!("{}-deployer", appname);
    let patch_params  = PatchParams::apply(&controller).force();


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
    let patch = Patch::Apply(deploy);

    deployments.patch(appname, &patch_params, &patch).await?;

    println!("✅ Successfully deployed {}", appname);
    Ok(())
}