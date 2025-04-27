package config

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

type connection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

var PostgresConn *sql.DB

func InitDB() {
	connInfo := connection{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASS"),
		DBName:   os.Getenv("PG_DBNAME"),
	}

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		connInfo.Host,
		connInfo.Port,
		connInfo.User,
		connInfo.Password,
		connInfo.DBName,
	)

	log.Println("DB connection string: ", connString)

	var db *sql.DB
	var err error

	maxRetries := 5

	for attempts := 1; attempts <= maxRetries; attempts++ {
		db, err = sql.Open("postgres", connString)
		if err != nil {
			log.Printf("Attempt %d: failed to open PostgreSQL connection: %v", attempts, err)
		} else {
			// Check the actual connection with Ping
			err = db.Ping()
			if err != nil {
				log.Printf("Attempt %d: failed to ping PostgreSQL: %v", attempts, err)
			} else {
				log.Println("Connected to PostgreSQL")
				break
			}
		}

		if attempts < maxRetries {
			log.Println("Retrying in 2 seconds...")
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL after %d attempts: %v", maxRetries, err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping PostgreSQL: %v", err)
		return
	}

	PostgresConn = db

	//creating table and inserting values if not created
	SetupSchemaAndSeedData()
}

func SetupSchemaAndSeedData() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL
	);`

	_, err := PostgresConn.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("failed to create products table: %v", err)
	}
	log.Println("Products table created or already exists.")

	// Check if table already has data
	var count int
	err = PostgresConn.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		log.Fatalf("failed to count products: %v", err)
	}

	if count == 0 {
		// Seed initial products
		seedQuery := `
		INSERT INTO products (name, price) VALUES
		('Wireless Mouse', 25.99),
		('Mechanical Keyboard', 89.50),
		('HD Monitor', 199.99),
		('USB-C Hub', 39.95),
		('Noise Cancelling Headphones', 129.00);
		`
		_, err = PostgresConn.Exec(seedQuery)
		if err != nil {
			log.Fatalf("failed to seed products: %v", err)
		}
		log.Println("Seeded 5 sample products.")
	} else {
		log.Printf("Products table already contains %d entries. Skipping seeding.\n", count)
	}
}

var RedisClient *redis.Client

func InitRedis() {
	host := os.Getenv("REDIS_HOST")     // e.g., localhost
	port := os.Getenv("REDIS_PORT")     // e.g., 6379
	password := os.Getenv("REDIS_PASS") // leave empty "" if no password

	addr := fmt.Sprintf("%s:%s", host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	ctx := context.Background()
	var err error
	maxRetries := 5

	for attempts := 1; attempts <= maxRetries; attempts++ {
		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			log.Println("Connected to Redis")
			break
		}

		log.Printf("Attempt %d: failed to connect to Redis: %v", attempts, err)

		if attempts < maxRetries {
			log.Println("Retrying Redis connection in 2 seconds...")
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("failed to connect to Redis after %d attempts: %v", maxRetries, err)
	}

	RedisClient = rdb
}
