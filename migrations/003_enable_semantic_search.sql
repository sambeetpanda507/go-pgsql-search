-- Create vector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Add index on embeddings column
CREATE INDEX IF NOT EXISTS idx_news_embedding_hnsw 
ON news USING hnsw (embedding vector_cosine_ops);