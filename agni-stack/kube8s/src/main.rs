mod k8s;

use k8s::deploy::deployk8s;

fn main() {
    println!("Hello, world!");
    deployk8s();
}
