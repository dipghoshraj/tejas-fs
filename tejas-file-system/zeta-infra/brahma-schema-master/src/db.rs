use postgres::{Client, NoTls, Error};

pub fn database_connection(db_url: &str) -> Result<Client, Error> {
 let db_client =  Client::connect(db_url, NoTls)?;
 Ok(db_client)
}

pub fn create_database(client: &mut Client, db_name: String) -> Result<(), Error> {
    let query = format!("CREATE DATABASE {}", db_name);
    client.batch_execute(&query)?;
    println!("INFO :: Creating database successful {}", db_name);
    Ok(())
}

pub fn drop_database(client: &mut Client, db_name: String) -> Result<(), Error> {
    terminate_connections(client, db_name.clone())?;

    let query = format!("DROP DATABASE {}", db_name);
    client.batch_execute(&query)?;
    println!("INFO :: Drop database successful {}",db_name);
    Ok(())
}


fn terminate_connections(client: &mut Client, db_name: String) -> Result<(), Error> {
    let terminate_query = format!(
        "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1 AND pid <> pg_backend_pid();"
    );

    client.execute(&terminate_query, &[&db_name])?;
    println!("Terminated all active connections to the database {}", db_name);
    Ok(())
}