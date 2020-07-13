package utils

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sajallimbu/go_securing_api/models"

	//postgres drivers
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// It is not recommended to store the database credentials in the code
// I have done this for convenience purpose
// Use something like joho/godotenv library to store the credentials inside a .env file in your project directory
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "secureapi"
)

// ConnectDB ... our database connection function. It return a pointer to a gorm Database
func ConnectDB() *gorm.DB {
	//Define Database connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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
