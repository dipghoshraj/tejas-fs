use k8s_openapi::api::networking::v1::Ingress;
use kube::{
    api::{Api, PostParams},
    Client,
};
use serde_json::json;
