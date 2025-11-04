-- Create documents table
CREATE TABLE IF NOT EXISTS documents (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_name   VARCHAR(255)    NOT NULL,
    document_path   VARCHAR(255)    NOT NULL UNIQUE,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL
);

-- Create index on document_path for faster lookups
CREATE INDEX IF NOT EXISTS idx_documents_path ON documents(document_path);
