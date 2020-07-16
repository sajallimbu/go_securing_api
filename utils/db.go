package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/jinzhu/gorm"
	"github.com/sajallimbu/go_securing_api/models"

	//postgres drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// ConnectDB ... our database connection function. It return a pointer to a gorm Database
func ConnectDB() *gorm.DB {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Map your database credentials to their respective variables
	host := os.Getenv("databaseHost")
	port := os.Getenv("databasePort")
	user := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	dbname := os.Getenv("databaseDbname")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("Error parsing string type of port into int")
	}

	//Define Database connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, portInt, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf(err.Error())
		panic("Failed to connect to the database")
	}
	// defer db.Close()

	// AutoMigrate automatically creates a table in the database with some additional columns for convenience
	db.AutoMigrate(&models.User{})

	fmt.Println("Successfully connected to the database:")
	return db
}
