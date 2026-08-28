CREATE TABLE message(
    id SERIAL NOT NULL PRIMARY KEY,
    message_id VARCHAR(100),
    friendship_id VARCHAR(100),
    sender_username VARCHAR(100),
    message_type VARCHAR(100),
    text_content VARCHAR(255),
    media_url VARCHAR(200),
    media_type VARCHAR(100),
    created_at TIMESTAMP  WITH TIME ZONE DEFAULT NOW() NOT NULL,
modified_at TIMESTAMP   
)