
CREATE TYPE status AS ENUM (
    'active',
    'draft',
    'disabled',
    'deleted'
);

CREATE TABLE realms (
    key         UUID PRIMARY KEY,
    id          UUID NOT NULL,
    name        VARCHAR(50) NOT NULL,
    description TEXT,
    status      status  NOT NULL,
    created_at  TIMESTAMP   NOT NULL,
    updated_at  TIMESTAMP   NOT NULL,
    deleted_at  TIMESTAMP
);

