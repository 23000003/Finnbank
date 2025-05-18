-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
    CREATE TABLE BankCard (
        bankcard_ID SERIAL PRIMARY KEY,
        card_number VARCHAR(16) NOT NULL UNIQUE,
        expiry_date DATE NOT NULL,
        account_id UUID NOT NULL,
        cvv CHAR(3) NOT NULL,                         
        pin_number CHAR(4) NOT NULL,                        
        card_type VARCHAR(10) CHECK (Card_Type IN ('Debit', 'Credit')),
        date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
    DROP TABLE BankCard;
