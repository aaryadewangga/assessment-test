-- migrate:up
CREATE INDEX idx_users_username ON _users ("USERNAME");
CREATE INDEX idx_products_name ON _products ("PRODUCT_NAME");
CREATE INDEX idx_transactions_user_id ON _transactions ("USER_ID");
CREATE INDEX idx_transaction_details_transaction_id ON _transaction_details ("TRANSACTION_ID");
CREATE INDEX idx_transaction_details_product_id ON _transaction_details ("PRODUCT_ID");

-- migrate:down
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP INDEX IF EXISTS idx_transaction_details_transaction_id;
DROP INDEX IF EXISTS idx_transaction_details_product_id;


