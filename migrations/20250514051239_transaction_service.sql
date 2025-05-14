-- +goose Up

CREATE TABLE IF NOT EXISTS public.transactions (
  transaction_id     SERIAL    PRIMARY KEY,                    -- auto‑incrementing ID
  ref_no             VARCHAR(12)         NOT NULL UNIQUE                 -- your 12‑digit code
      CHECK (char_length(ref_no) = 12),
  sender_id          INT         NOT NULL,                       -- matches OpenedAccounts UUID
  receiver_id        INT,                                   -- nullable
  transaction_type   TEXT         NOT NULL                       -- use CHECK or SQL enum if desired
      CHECK (transaction_type IN ('Transfer','Payment','Deposit','Withdraw','Loan')),
  amount             NUMERIC(10,2) NOT NULL,
  transaction_status TEXT         NOT NULL DEFAULT 'Pending'      -- likewise validate with CHECK
      CHECK (transaction_status IN ('Pending','Completed','Failed')),
  date_transaction   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),         -- server time
  transaction_fee    NUMERIC(10,2) NOT NULL DEFAULT 0.00,
  notes              TEXT
); 


-- +goose Down
DROP TABLE IF EXISTS public.transactions;
