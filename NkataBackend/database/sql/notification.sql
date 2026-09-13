CREATE TABLE message(
    id SERIAL NOT NULL PRIMARY KEY,
    user_id VARCHAR(100),
    username VARCHAR(100),
    friendship_id VARCHAR(100),
    title VARCHAR(100),
    message VARCHAR(500),
    url_img VARCHAR(255),
    created_at TIMESTAMP  WITH TIME ZONE DEFAULT NOW() NOT NULL,
)