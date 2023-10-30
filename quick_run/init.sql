SET TIME ZONE '+8';

CREATE TABLE certs (
    sn TEXT PRIMARY KEY NOT NULL,
    key TEXT,
    note TEXT
);

CREATE TABLE temporary_permits (
    key TEXT PRIMARY KEY NOT NULL,
    expiration TIMESTAMP WITH TIME ZONE NOT NULL
);