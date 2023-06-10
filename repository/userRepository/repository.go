package userRepository

import (
	"bankai/models"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByUserId(userId uint) (*models.User, error)
	DeleteUser(user *models.User) error
}

type userGormRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() UserRepository {
	return &userGormRepository{
		db: getDbConnection(),
	}
}

func (ur *userGormRepository) CreateUser(user *models.User) error {
	return ur.db.Create(user).Error
}

func (ur *userGormRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userGormRepository) DeleteUser(user *models.User) error {
	return ur.db.Delete(user).Error
}

func (ur *userGormRepository) GetUserByUserId(userId uint) (*models.User, error) {
	var user models.User
	result := ur.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func getDbConnection() *gorm.DB {
	// Replace the values below with your actual Postgres database credentials
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	fmt.Println(dbURI)
	// Connect to the database
	db, err := gorm.Open(postgres.Open(dbURI))
	if err != nil {
		panic(err)
	}

	// Set up connection pool and other configuration options
	if db, err := db.DB(); err != nil {
		panic(err)
	} else {
		db.SetMaxOpenConns(100)
		db.SetMaxIdleConns(10)
	}
	// Migrate the User model to the database (if necessary)
	db.AutoMigrate(&models.User{})

	// Use the db instance to interact with the database in your application
	return db
}
