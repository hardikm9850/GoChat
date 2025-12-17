CREATE TABLE messages
(
    id              CHAR(36) PRIMARY KEY,
    conversation_id CHAR(36) NOT NULL,
    sender_id       CHAR(36) NOT NULL,
    content         TEXT     NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_message_conversation
        FOREIGN KEY (conversation_id)
            REFERENCES conversations (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_message_sender
        FOREIGN KEY (sender_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);
