BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS urls(
                                         id SERIAL PRIMARY KEY,
                                         short_url VARCHAR(255) NOT NULL,
                                         original_url TEXT NOT NULL,
                                         user_id VARCHAR(255) NOT NULL,
                                         deleted BOOLEAN default FALSE
);
CREATE UNIQUE INDEX idx_unique_short_url ON urls(short_url) WHERE deleted = FALSE;
CREATE UNIQUE INDEX idx_unique_original_url ON urls(original_url) WHERE deleted = FALSE;

COMMIT;