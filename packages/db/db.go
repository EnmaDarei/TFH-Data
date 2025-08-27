package database

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var pgUser string
var pgPassword string
var pgHost string
var pgPort string

var valid_tfh_characters = []string{"arizona", "oleander", "paprika", "pom", "shanty", "stronghoof", "tianhuo", "velvet", "texas", "nidra", "baihe"}

func Initialize_Database() {
	pgUser = os.Getenv("PG_USER")
	pgPassword = os.Getenv("PG_PASSWD")
	pgHost = os.Getenv("PG_HOST")
	pgPort = os.Getenv("PG_PORT")
	connection_str := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", pgUser, "tfh", pgPassword, pgHost, pgPort)
	db, err := sql.Open("postgres", connection_str)
	if err != nil {
		fmt.Println("Error initializing database connection pool:", err)
		return
	}
	DB = db
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
}

func Check_Valid_Character(character string) (string, error) {
	if !slices.Contains(valid_tfh_characters, strings.ToLower(character)) {
		return "", fmt.Errorf("%s is not a real tfh character, deer", character)
	}
	return character, nil
}
