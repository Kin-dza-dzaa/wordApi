CREATE ROLE wordapi WITH LOGIN;

CREATE DATABASE wordapi WITH OWNER wordapi;

\c wordapi wordapi

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
    time_of_last_repeating TIMESTAMP                                                        NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words(id),
    UNIQUE(user_id, word_id, collection_name)
);