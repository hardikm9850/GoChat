CREATE TABLE conversation_participants
(
    conversation_id CHAR(36) NOT NULL,
    user_id         CHAR(36) NOT NULL,
    joined_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (conversation_id, user_id),

    CONSTRAINT fk_cp_conversation
        FOREIGN KEY (conversation_id)
            REFERENCES conversations (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_cp_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE
);
