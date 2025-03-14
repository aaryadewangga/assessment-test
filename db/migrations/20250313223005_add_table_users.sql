-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE _users
(
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "NAME" VARCHAR(255) NOT NULL,
    "USERNAME" VARCHAR(255) UNIQUE NOT NULL,
    "PASSWORD" VARCHAR(255) NOT NULL,
    "ROLE" VARCHAR(255) NOT NULL,
    "CREATED_AT" TIMESTAMPTZ DEFAULT NOW()
);

-- migrate:down
DROP TABLE IF EXISTS _users;
DROP EXTENSION IF EXISTS "uuid-ossp";