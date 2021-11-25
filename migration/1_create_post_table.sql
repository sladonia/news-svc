CREATE TABLE post
(
    id VARCHAR(20) NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);

CREATE INDEX created_at_idx on post using btree(created_at);
