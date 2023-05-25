CREATE IF NOT EXISTS TABLE folders(
    id serial,
    name varchar(60) NOT NULL,
    parent_id INT,
    created_at TIMESTAMP current_timestamp,
    modifier_at TIMESTAMP NOT NULL,
    deleted bool NOT NULL DEFAULT false,

    PRIMARY KEY(id),
    CONSTRAINT fk_parent_folders
        FOREIGN KEY(parent_id) REFERENCES folders(id)
);