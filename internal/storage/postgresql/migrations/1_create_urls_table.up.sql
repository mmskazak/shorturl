BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS urls(
                                         id SERIAL PRIMARY KEY,
                                         short_url VARCHAR(255) NOT NULL,
                                         original_url TEXT NOT NULL,
                                         user_id VARCHAR(255) NOT NULL,
                                         deleted BOOLEAN default FALSE
);
-- Создаем индекс, ограничивающий уникальность original_url только для строк, где deleted = FALSE
CREATE UNIQUE INDEX unique_original_url ON urls (original_url) WHERE deleted = FALSE;

COMMIT;
