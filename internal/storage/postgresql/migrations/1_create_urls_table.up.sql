BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS urls (
                                    id SERIAL PRIMARY KEY,
                                    short_url VARCHAR(255) NOT NULL,
                                    original_url TEXT NOT NULL,
                                    user_id VARCHAR(255) NOT NULL,
                                    deleted BOOLEAN DEFAULT FALSE,
                                    CONSTRAINT unique_short_url UNIQUE (short_url),
                                    CONSTRAINT unique_original_url UNIQUE (original_url)
);

COMMIT;


