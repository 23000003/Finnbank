# Online Banking System
IT 3101A - Microservices / QA Final Project

## Golang Commands

   ```bash
   go run ./services/<ur-assigned-service> . 
   go get -u ./<ur-assigned-service> # Import external packages
   ./run_all_services.bat # run all services
   go mod tidy # on ur specific service directory for package cleanup/update
   ```

## Note:
Utillize ``logger`` from finnbank/common/utils for terminal logs

```
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
    OpenedAccount_Status ENUM('Active', 'Suspended', 'Closed'),
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
    Notif_Type ENUM('Transaction', 'System'),
    Notif_To_ID FK,                   -- Link to the user's account (foreign key)
    Notif_From_Name VARCHAR(255),      -- Sender of the notification (system/user name)
    Content TEXT,                      -- The content/message of the notification
    IsRead BOOLEAN DEFAULT FALSE,      -- Whether the notification has been read
    Redirect_Link VARCHAR(255),        -- Link for more details
    Date_Notified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Date_Read TIMESTAMP               -- Timestamp of when the user read the notification (nullable)
)
```

# ToDo:
```
graphql-api
├── entities
│   └── Put your entities here
└── handlers
    └── Initialize the GraphQL handler here (no gRPC anymore, base it from accountHandler)
└── resolvers
    └── Initialize the GraphQL resolvers here with the respective resolver type
    └── Call the business logic here (the services)
└── services
    └── Put your business logics here
└── db
    └── Initialize your DB here
```
    
## Account
 - This only serves as Auth/Autho
 - ikaw na bahala saimo fn
 - No need complexity validation sa Register just input (unya ra R.T. validation if we are done with evrything)
 - (FE) When register, user has to pick a Account Type that will be called at OpenedAccount via Microservice to create that account type

## Opened Account (me)
 - User should can create an Account with the account type that he still doesnt have.
 - Account Types: **Credit** - Generates CreCard | **Checking** - Generates DebCard | **Savings** - No Card
 - Every Open It will call BankCard to Generate the card
 - Functions:
    - (GetAllOpenedAccountsByUserId, GetOpenedAccountOfUserById, OpenAnAccountByAccountType, UpdateOpenedAccountStatus)

## BankCard
  - Functions:
    - (GetBankCardOfUserById, CreateBankCardForUser, UpdateBankcardExpiryDateByUserId)
   
## Transaction
  - Can only execute a transaction if he has an opened account with that account Type (this will be done in frontend so no complex Logics)
  ![image](https://github.com/user-attachments/assets/544279c3-3077-4218-a57a-4fb8c2588fdc)
  - Functions:
    - (GetTransactionByUserId, GetTransactionByTimeStampByUserId, CreateTransactionByUserId)

## Statement
  - This only Generates a pdf, calls the GetTransactionByTimeStampByUserId to get data and generate and pass the pdf via buffer
  - FN: GenerateStatementByTimeStampForUser

## Notification
  - Notifies if every Transaction Trigger or System Trigger
  - Function:
    - GetAllNotificationByUserId, GetNotificationByUserId, GenerateNotification, ReadNotificationByUserId


## Service URL's
   ```bash
   # http ports (for data response test)
   http://localhost:8080/api # api-gateway 
   http://localhost:8081/api/bankcard # bankcard-api 
   http://localhost:8082/api/account # account-api
   http://localhost:8083/api/graphql # graphql-api (not http)
   http://localhost:8084/api/statement # statement-api
   http://localhost:8085/api/transaction # transaction-api
   http://localhost:8086/api/notification # notification-api


   # gRPC ports (main communication)
   http://localhost:9001/api/bankcard # bankcard-grpc
   http://localhost:9002/api/account # account-grpc
   http://localhost:9004/api/statement # statement-grpc 
   http://localhost:9005/api/transaction # transaction-grpc 
   http://localhost:9006/api/notification # notification-grpc
   ```

## Generate proto for gRPC communication
   - Only do this if you added a proto file in ``Protobuf Directory``
   - This is to create gRPC services/communication of your respective service
   ```bash
   # install first
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   $env:PATH="$PATH:$(go env GOPATH)/bin"
   protoc --version # check installation

   # if its unrecognized then try:
   winget install protobuf

   # update proto_gen.bat
   start protoc \
        --proto_path=protobuf "protobuf/<change-lines-with-this-format>.proto" \

   # then generate
   ./proto_gen.bat
   ```

## Microservice Architecture [Not Final]
![alt text](PROJECT.drawio%20(1).png)