CREATE IF NOT EXISTS TABLE users (
    id serial,
    name VARCHAR(80) NOT NULL,
    login VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(200) NOT NULL,
    created_at TIMESTAMP current_timestamp,
    modifier_at TIMESTAMP NOT NULL,
    last_login TIMESTAMP DEFAULT current_timestamp,
    deleted bool NOT NULL DEFAULT false,

    PRIMARY KEY(id)
);