use dotenv::dotenv;
use std::env;

pub fn load_config() {
    dotenv().ok();
    let db_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    println!("Loaded database URL: {}", db_url);
}