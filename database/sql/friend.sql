CREATE TABLE friendship (
id SERIAL NOT NULL PRIMARY KEY ,
firendship_id VARCHAR(255),
user_id INT,
friend_id INT,
friendship_type VARCHAR(255),
group_id INT,
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
)

CREATE TABLE friendRequest (
id SERIAL NOT NULL PRIMARY KEY ,
sent_by INT,
sent_to INT,
status VARCHAR(255),
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
modified_at TIMESTAMP    
)