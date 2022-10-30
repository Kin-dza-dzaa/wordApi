CREATE ROLE wordapi WITH PASSWORD '12345' LOGIN;

CREATE DATABASE wordapi WITH OWNER wordapi;

\c wordapi wordapi;

CREATE TABLE words(
    word                TEXT                                                                NOT NULL CHECK(word != ''),
    trans_data          JSONB                                                               NOT NULL,
    PRIMARY KEY (word)
);

CREATE TABLE user_collection(
    user_id                                     UUID                                        NOT NULL,
    word                                        TEXT                                        NOT NULL,
    state                                       INTEGER                                     NOT NULL,
    collection_name                             TEXT                                        NOT NULL,
    time_of_last_repeating                      TIMESTAMP                                   NOT NULL,
    FOREIGN KEY (word) REFERENCES words(word),
    UNIQUE(user_id, word, collection_name)
);