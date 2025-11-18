-- Drop indexes
DROP INDEX IF EXISTS idx_tokens_status;
DROP INDEX IF EXISTS idx_tokens_issuer_did;
DROP INDEX IF EXISTS idx_tokens_account_id;

-- Drop tokens table
DROP TABLE IF EXISTS tokens;
