-- +goose Up
CREATE TYPE account_type AS ENUM ('checking', 'credit', 'savings');
CREATE TYPE account_status AS ENUM ('Active', 'Suspended', 'Closed');

CREATE TABLE OpenedAccount (
    OpenedAccount_ID SERIAL PRIMARY KEY,
    Account_ID INTEGER NOT NULL,
    BankCard_ID INTEGER NOT NULL,
    Account_Type account_type NOT NULL,
    Balance DECIMAL(10, 2) DEFAULT 0.00,
    OpenedAccount_Status account_status DEFAULT 'Active',
    Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE OpenedAccount;
DROP TYPE account_type;
DROP TYPE account_status;
