DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS words CASCADE;
DROP TABLE IF EXISTS userword CASCADE;
DROP TABLE IF EXISTS wordstate CASCADE;
DROP TABLE IF EXISTS archivedwords CASCADE;

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
    PRIMARY KEY (id),
    UNIQUE (word)
);

CREATE TABLE userword(
    user_id UUID                                                                            NOT NULL,
    word_id INTEGER                                                                         NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (word_id) REFERENCES words(id),
    PRIMARY KEY(user_id, word_id)
);

CREATE TABLE wordstate(
    user_id UUID                                                                            NOT NULL,
    word_id INTEGER                                                                         NOT NULL,
    state   INTEGER                                                                         NOT NULL,
    FOREIGN KEY(user_id, word_id) REFERENCES userword(user_id, word_id) ON DELETE CASCADE
);

CREATE TABLE archivedwords(
    user_id UUID                                                                            NOT NULL,
    word_id INTEGER                                                                         NOT NULL,
    FOREIGN KEY(word_id) REFERENCES words(id),
    UNIQUE (user_id, word_id)
);
