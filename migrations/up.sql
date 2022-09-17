CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE  users(
    id UUID                                                                                 NOT NULL,
    user_name         TEXT                                                                  NOT NULL CHECk(user_name != ''),
    email             TEXT                                                                  NOT NULL CHECK(email != ''),
    password          TEXT                                                                  NOT NULL CHECK(password != ''),
    registration_date TIMESTAMP                                                             NOT NULL, 
    PRIMARY KEY (id),
    UNIQUE(user_name),
    UNIQUE (email)
);

CREATE TABLE words(
    id                  SERIAL                                                              NOT NULL,
    word                TEXT                                                                NOT NULL CHECK(word != ''),
    trans_data JSONB                                                                        NOT NULL,
    PRIMARY KEY (id),   
    UNIQUE (word)
);

CREATE TABLE user_collection(
    user_id UUID                                                                            NOT NULL,
    word_id INTEGER                                                                         NOT NULL,
    state INTEGER                                                                           NOT NULL,
    collection_name TEXT                                                                    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (word_id) REFERENCES words(id),
    PRIMARY KEY(user_id, word_id, collection_name)
);