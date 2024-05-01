package db_authSvc

import (
	"database/sql"
	"fmt"
	"time"

	domain_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/domain"
	config_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/config"
	interface_hash_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/hash_password/interface"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(config *config_authSvc.DataBase, hashUtil interface_hash_authSvc.IhashPassword) (*gorm.DB, error) {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s port=%s", config.DBHost, config.DBUser, config.DBPassword, config.DBPort)
	sql, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("-------", err)
		return nil, err
	}

	rows, err := sql.Query("SELECT 1 FROM pg_database WHERE datname = '" + config.DBName + "'")
	if err != nil {
		fmt.Println("Error checking database existence:", err)
	}
	defer rows.Close()

	if rows.Next() {
		fmt.Println("Database" + config.DBName + " already exists.")
	} else {
		_, err = sql.Exec("CREATE DATABASE " + config.DBName)
		if err != nil {
			fmt.Println("Error creating database:", err)
		}
	}

	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", config.DBHost, config.DBUser, config.DBName, config.DBPort, config.DBPassword)
	DB, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC() // Set the timezone to UTC
		},
	})
	if dberr != nil {
		return DB, nil
	}

	// Table Creation
	if err := DB.AutoMigrate(&domain_authSvc.Admin{}); err != nil {
		return DB, err
	}
	if err := DB.AutoMigrate(&domain_authSvc.Users{}); err != nil {
		return DB, err
	}
	if err := DB.AutoMigrate(&domain_authSvc.OtpInfo{}); err != nil {
		return DB, err
	}

	CheckAndCreateAdmin(DB, hashUtil)
	return DB, nil
}
func CheckAndCreateAdmin(DB *gorm.DB, hashUtil interface_hash_authSvc.IhashPassword) {
	var count int
	var (
		Name     = "ciao"
		Email    = "ciao@gmail.com"
		Password = "ciaociao"
	)
	HashedPassword := hashUtil.HashPassword(Password)

	query := "SELECT COUNT(*) FROM admins"
	DB.Raw(query).Row().Scan(&count)
	if count <= 0 {
		query = "INSERT INTO admins(name, email, password) VALUES(?, ?, ?)"
		DB.Exec(query, Name, Email, HashedPassword).Row().Err()
	}
}
