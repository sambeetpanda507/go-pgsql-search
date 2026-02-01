-- pg_trgm is used to find similarity using levenstine algorithm
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create GIST index for fuzzy search on title and description
CREATE INDEX IF NOT EXISTS IDX_NEWS_FUZZY_SEARCH ON NEWS USING GIST (
	(
		COALESCE(TITLE, '') || ' ' || COALESCE(DESCRIPTION, '')
	) GIST_TRGM_OPS
);