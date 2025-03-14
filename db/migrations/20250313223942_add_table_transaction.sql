-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE _transactions
(
    "ID" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "USER_ID" UUID NOT NULL,
    "TOTAL_AMOUNT" NUMERIC(15,2) NOT NULL,
    "CREATED_AT" TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY ("USER_ID") REFERENCES _users("ID") ON DELETE CASCADE
);

-- migrate:down
DROP TABLE IF EXISTS _transactions;
DROP EXTENSION IF EXISTS "uuid-ossp";