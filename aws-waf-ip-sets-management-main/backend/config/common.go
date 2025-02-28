package config

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/wafv2"
)

var (
    DB    *sql.DB
    WAFv2 *wafv2.WAFV2
)

func init() {
    // Load environment variables from .env file
    if err := godotenv.Load("../config.env"); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initialize AWS WAFv2
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION")),
        Credentials: credentials.NewStaticCredentials(
            os.Getenv("AWS_ACCESS_KEY_ID"),
            os.Getenv("AWS_SECRET_ACCESS_KEY"),
            ""),
    }))
    WAFv2 = wafv2.New(sess)

    // Initialize MySQL database connection
    dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME")
    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Error connecting to database:", err)
    }
    if err := DB.Ping(); err != nil {
        log.Fatal("Error pinging database:", err)
    }
}
