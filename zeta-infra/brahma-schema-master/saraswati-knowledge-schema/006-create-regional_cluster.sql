create table regional_clsuters (
    id SERIAL PRIMARY KEY,
    name varchar(255) not null,
    topics varchar(255) not null,
    type varchar(255) not null,
    network varchar(255) not null,
    connection varchar(255) not null,
    configuration JSON NOT null,
    storage varchar(255) not null
)