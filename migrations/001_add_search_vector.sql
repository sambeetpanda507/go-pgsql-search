-- For full text search or lexical search
ALTER TABLE news
ADD COLUMN IF NOT EXISTS search_vector tsvector GENERATED ALWAYS AS
(
    to_tsvector('english', coalesce(title, '') || ' ' || coalesce(description, ''))
) STORED;

CREATE INDEX IF NOT EXISTS idx_news_search_vector ON news USING GIN(search_vector);