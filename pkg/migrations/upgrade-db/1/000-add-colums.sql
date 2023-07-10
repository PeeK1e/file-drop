ALTER TABLE file
    ADD COLUMN is_encrypted BOOLEAN DEFAULT False,
    ADD COLUMN secret_sha VARCHAR(128) DEFAULT NULL,
    ADD COLUMN max_downloads INT DEFAULT -1
;