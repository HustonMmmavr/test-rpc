create table users(
    uuid varchar(36) UNIQUE, 
    login text UNIQUE,
    created TIMESTAMP
 );