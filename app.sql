Account (
    Account_ID PK,
    Email VARCHAR(255) NOT NULL UNIQUE,
    First_Name VARCHAR(100),
    Last_Name VARCHAR(100),
    Surname VARCHAR(100),
    Phone_Number VARCHAR(20),
    Password VARCHAR(255) NOT NULL,   -- Store hashed passwords, not plain text
    Address TEXT,
    nationalIdNumber VARCHAR(20) UNIQUE,  -- Unique national ID number for each user
    Account_Number VARCHAR(20) UNIQUE,  -- Unique account number for each user
    Account_Type ENUM('Business', 'Personal'),
    Account_Status ENUM('Active', 'Suspended', 'Closed'),
    Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
)

OpenedAccount (
    OpenedAccount_ID INT PRIMARY KEY AUTO_INCREMENT,
    Account_ID INT,                           -- Foreign key to the user
    BankCard_ID INT,                     -- Foreign key to the bank card
    Account_Type ENUM('checking', 'credit', 'savings'),  -- Type of account #savings no bankcard
    Balance DECIMAL(10, 2),                -- Amount of money available
    Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

BankCard (
    BankCard_ID INT PRIMARY KEY AUTO_INCREMENT,
    Card_Number VARCHAR(16) NOT NULL UNIQUE,  -- Masked card number
    Expiry DATE,
    Account_ID INT,                         -- Foreign key to the Account table
    CVV CHAR(3),                            -- Securely store CVV or avoid storing it
    PinNumber CHAR(4),                      -- Securely store the PIN
    Card_Type ENUM('debit', 'credit'),      -- Card type (either debit or credit)
    Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

Transaction (
    Transaction_ID PK,                  -- Unique ID for each transaction
    Ref_No VARCHAR(20) NOT NULL UNIQUE,  -- Transaction reference number
    Sender_ID FK,                       -- Foreign key to the OpenedAccount (sender)
    Receiver_ID FK,                     -- Foreign key to the OpenedAccount (receiver)
    Transaction_Type ENUM('Withdraw', 'Deposit', 'Loan', 'Transfer', 'Payment'),  -- Type of transaction
    Amount INT NOT NULL,                 -- Store in cents for accuracy
    Transaction_Status ENUM('Completed', 'Failed'),
    Date_Transaction TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Transaction_Fee INT DEFAULT 0,      -- Any applicable transaction fees
    Notes TEXT                           -- Optional notes for special cases
)

Notification (
    Notif_ID PK,                      -- Primary Key
    Notif_Type ENUM('Transaction', 'Account', 'Security', 'System'),
    Notif_To_ID FK,                   -- Link to the user's account (foreign key)
    Notif_From_Name VARCHAR(255),      -- Sender of the notification (system/user name)
    Content TEXT,                      -- The content/message of the notification
    IsRead BOOLEAN DEFAULT FALSE,      -- Whether the notification has been read
    Redirect_Link VARCHAR(255),        -- Link for more details
    Date_Notified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Date_Read TIMESTAMP               -- Timestamp of when the user read the notification (nullable)
)
