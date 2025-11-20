CREATE TABLE friendship (
id SERIAL NOT NULL PRIMARY KEY ,
friendship_id VARCHAR(255),
username VARCHAR(255),
last_message VARCHAR(255),
friend_username VARCHAR(255),
friendship_type VARCHAR(255),
group_id INT,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
modified_at TIMESTAMP    
)

CREATE TABLE friendRequest (
id SERIAL NOT NULL PRIMARY KEY ,
sent_by VARCHAR(255),
sent_to VARCHAR(255),
status VARCHAR(255),
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
modified_at TIMESTAMP    
)