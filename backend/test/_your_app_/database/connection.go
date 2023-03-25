package database

import (
	"fmt"

	organisations "github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/models/organisation"
	users "github.com/HousewareHQ/backend-engineering-octernship/test/_your_app_/models/user"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var EnvMap map[string]string
var DB *gorm.DB

func init() {
	fmt.Println("Initializing DB connection")

	EnvMap = map[string]string{
		"SECRET_KEY":     "kittens-are-cute",
		"NGINX_PORT":     "8080",
		"DB_USERNAME":    "root",
		"DB_PASSWORD":    "mypassword",
		"DB_HOST":        "localhost",
		"DB_PORT":        "3306",
		"DB_NAME":        "mydb",
		"ADMIN_ORG_NAME": "Admin",
		"ADMIN_ORG_HEAD": "Rohit-1",
		"ADMIN_USERNAME": "Rohit-1",
		"ADMIN_PASSWORD": "pass-1",
	}

	// ! connect to DB
	connectionURI := EnvMap["DB_USERNAME"] + ":" + EnvMap["DB_PASSWORD"] + "@tcp(" + EnvMap["DB_HOST"] + ":" + EnvMap["DB_PORT"] + ")/" + EnvMap["DB_NAME"] + "?charset=utf8mb4&parseTime=True&loc=Local"
	connection, err := gorm.Open(mysql.Open(connectionURI), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to DB")
	}

	// ! save connection
	DB = connection

	// ! migrate models
	DB.AutoMigrate(&organisations.Organisation{})
	DB.AutoMigrate(&users.User{})
}

func init() {
	fmt.Println("Creating Admin Account")
	// ! check if admin account exists
	var foundUser users.User
	DB.Where("is_admin = ?", true).First(&foundUser)

	// ! if not, create admin account
	if foundUser.Id == 0 {
		// ! create admin org
		adminOrg := organisations.Organisation{
			Name: EnvMap["ADMIN_ORG_NAME"],
			Head: EnvMap["ADMIN_ORG_HEAD"],
		}

		// ! save to db
		DB.Create(&adminOrg)

		// ! hash password
		password := EnvMap["ADMIN_PASSWORD"]
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 8)
		if err != nil {
			fmt.Println("Failed to hash password")
			return
		}
		// ! create admin user
		user := users.User{
			Username: EnvMap["ADMIN_USERNAME"],
			Password: hashedPass,
			IsAdmin:  true,
			OrgId:    adminOrg.Id,
		}
		DB.Create(&user)
	}
}
