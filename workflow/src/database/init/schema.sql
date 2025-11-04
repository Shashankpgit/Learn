-- Enable the UUID extension (run once per database)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create table if not exists
-----------------------------------

CREATE TABLE IF NOT EXISTS actors (
    actor_id UUID PRIMARY KEY,
    did UUID NOT NULL UNIQUE,
    email VARCHAR NOT NULL UNIQUE,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    phone_number VARCHAR UNIQUE,
    master_public_key VARCHAR NOT NULL UNIQUE,
    entity_type VARCHAR NOT NULL,
    nationality VARCHAR,
    country_of_residence VARCHAR,
    country_of_incorporation VARCHAR,
    verification_level VARCHAR NOT NULL DEFAULT 'Tier0_Unverified',
    created_at TIMESTAMPTZ NOT NULL
);

-----------------------------------

CREATE TABLE IF NOT EXISTS identifiers (
    identifier VARCHAR PRIMARY KEY, 
    entity_type VARCHAR NOT NULL,
    entity_id UUID NOT NULL
);

-----------------------------------

CREATE TABLE IF NOT EXISTS actor_integrations (
    actor_id UUID PRIMARY KEY,
    provider VARCHAR NOT NULL,
    external_user_id VARCHAR NOT NULL,
    linked_at TIMESTAMPTZ NOT NULL
);

-----------------------------------

CREATE TABLE IF NOT EXISTS public.tokens (
	token_id uuid PRIMARY KEY,
	account_id uuid NOT NULL,
	token_type varchar,
	issuer_did varchar,
	token_standard varchar,
	status varchar,
	metadata jsonb,
	created_at timestamptz DEFAULT now()
);

-----------------------------------

CREATE TABLE IF NOT EXISTS documents (
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_name   VARCHAR(255)    NOT NULL,
    document_path   TEXT            NOT NULL UNIQUE,
    created_at      TIMESTAMP       DEFAULT CURRENT_TIMESTAMP  NOT NULL
);

