CREATE TABLE otp(
id SERIAL NOT NULL PRIMARY KEY,
username VARCHAR(255),
token INT,
email VARCHAR(255),
purpose VARCHAR(255),
exp TIMESTAMP,
created_at TIMESTAMP ,
modified_at TIMESTAMP    
)