package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	err error
)

type Student struct {
	Name string
	Age  int64
}

func connectDB() (*gorm.DB, error) {
	log.Println("Connecting to database")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		"5432",
	)
	log.Println("Connection STR : ", dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Println("DB connection error : ", err)
		return nil, err
	}

	if err := db.AutoMigrate(&Student{}); err != nil {
		log.Println("Error while migrating :", err)
		return nil, err
	}

	log.Println("DB conccetion has been established")
	return db, nil
}

func testRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("Request Build")
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	// create user
	user := Student{Name: "User", Age: int64(time.Now().Second())}
	result := db.Create(user)
	var response map[string]interface{}
	if result.Error != nil {
		response = map[string]interface{}{
			"status": false,
			"error":  "faild to create user",
		}

	} else {
		response = map[string]interface{}{
			"status":  true,
			"message": "User created successfully",
			"data":    user,
		}
	}

	json.NewEncoder(w).Encode(response)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	var users []Student
	var response map[string]interface{}
	if db == nil {
		connectDB()
	}
	result := db.Session(&gorm.Session{}).Find(&users)
	if result.Error != nil {
		response = map[string]interface{}{
			"status": false,
			"error":  "faild to create user",
		}
	} else {
		response = map[string]interface{}{
			"status":  true,
			"message": "User created successfully",
			"data":    users,
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error while loading env file : ", err)
	}

}
func main() {
	fmt.Println("Program run success")
	loadEnv()
	_, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	http.HandleFunc("/user", testRoute)
	http.HandleFunc("/get-user", getUser)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("Error : ", err)
	}
}
