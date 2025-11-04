-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create tokens table
CREATE TABLE IF NOT EXISTS tokens (
    token_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    token_type VARCHAR,
    issuer_did VARCHAR,
    token_standard VARCHAR,
    status VARCHAR,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create index on account_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_tokens_account_id ON tokens(account_id);

-- Create index on issuer_did for faster lookups
CREATE INDEX IF NOT EXISTS idx_tokens_issuer_did ON tokens(issuer_did);

-- Create index on status for faster filtering
CREATE INDEX IF NOT EXISTS idx_tokens_status ON tokens(status);

-- Add comment to the table
COMMENT ON TABLE tokens IS 'Table for storing verifiable credentials and tokens';
