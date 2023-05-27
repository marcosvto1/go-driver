CREATE TABLE files (
    id serial,
    folder_id SERIAL,
    owner_id serial NOT NULL,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL,
    path VARCHAR(250) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp,
    modifier_at TIMESTAMP NOT NULL,
    deleted bool NOT NULL DEFAULT false,

    PRIMARY KEY(id),
    CONSTRAINT fk_folders FOREIGN KEY(folder_id) REFERENCES folders(id),
    CONSTRAINT fk_owner FOREIGN KEY(owner_id) REFERENCES users(id)
);