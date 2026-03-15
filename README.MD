# E-Wallet System

Tech Stack :
- ![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
- ![GORM](https://img.shields.io/badge/GORM-0D8ABC?style=for-the-badge&logo=go&logoColor=white)
- ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4479A1?style=for-the-badge&logo=postgresql&logoColor=white)


## Features:
- **User:** CRUD 
- **Wallet :**
    - Wallet Creation – create wallet per currency for a user
    - Top-up – add decimal money in wallet currency
    - Payment – spend decimal money from wallet
    - Transfer – move decimal money to another user’s wallet of same currency
    - Suspend Wallet – block operations
    - Query Wallet – view balance and status
    - Multi-Currency Awareness – prevent cross-currency operations
    - currency fixed 3 letter ISO
    - decimal max 2 digit
- **Ledger :**
    - Records every change to a wallet’s balance
    - Ledger entries are append-only
    - Ledger tracks currency and decimal amount per entry

## Set up local :

- Clone repo
- Rename file config.yml.example to config.yml
- Database use is Postgresql only
- Change the contents of the config.yml file according to the database name, username, password, host, port, redis address and redis password respectively.
  ```bash
  postgres:
  host: "127.0.0.1"
  port: "5432"
  username: "postgres"
  password: "Terserah123"
  dbname: "ewallet"
  idleconnect: 10
  maxconnect: 100
  lifeconnect: 300

  redis:
  host: "127.0.0.1"
  port: "6379"
  db: 0
  password:
  username:
  ```

- Run in terminal

  ```bash
  go run cmd/main.go
  ```

- URL Documentation API : https://documenter.getpostman.com/view/22397647/2sBXigMYhV
