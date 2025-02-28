# AWS WAF IP Sets Management

## Setup

1. Clone the repository.
2. Create a `config.env` file at the root of the project and add your AWS and DB credentials as shown in the `config.env` example.
3. Install dependencies using `go mod tidy`.
4. Run the server using `go run backend/main.go`.

## Project Structure
```
aws-waf-ip-sets-management/
│
├── backend/
│   ├── backend (compiled binary)
│   ├── main.go
│   ├── config/
│   │   └── common.go
│   ├── routes/
│   │   ├── createIPSet.go
│   │   ├── addIPAddress.go
│   │   ├── removeIPAddress.go
│   │   ├── deleteIPSet.go
│   │   └── listIPSets.go
│   └── utils/
│       └── utils.go
│
├── frontend/
│   ├── index.html
│   ├── script.js
│   └── styles.css
│
├── go.mod
├── go.sum
├── config.env
└── Dockerfile

```

#### create db and Tables
    ```
    USE wafdb;

    CREATE TABLE actions (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        type VARCHAR(50) NOT NULL,
        action TEXT NOT NULL,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    ```

#### create complier
    ```
    go build -o backend main.go
    ```

#### startup DB in local
```
docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=wafdb -p 3306:3306 -d mysql:latest
```
