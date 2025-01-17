use postgres::{Client, Error};
use std::fs;  // Synchronous file handling from std



fn get_migration_file_path(file_name:  &mut String) -> String {
    format!("saraswati-knowledge-schema/{}.sql", file_name)  // Returns the full path to the SQL file
}


fn apply_single_migration(client: &mut Client, file_name:  &mut String)-> Result<(), Error> {

    let file_apth = get_migration_file_path(file_name);

    let sql = fs::read_to_string(file_apth).expect("Failed to read migration file");
    client.batch_execute(&sql)?;
    Ok(())
    

}

pub fn apply_migration(client: &mut Client, file_name: &mut String) -> Result<(), Error> {
    apply_single_migration(client, file_name)
}