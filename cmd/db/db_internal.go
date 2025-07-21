package db

import (
	"darius/models"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database interface {
	CreateReport(entry, res, resp, requestKey string, amount float64) error
	GetByRequestKey(requestKey string) (*models.LLMCallReport, error)
}

type db struct {
	DB *gorm.DB
}

func NewDatabase() (Database, error) {
	db := &db{}
	err := db.connectDatabase()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}
	db.DB.AutoMigrate(&models.LLMCallReport{})
	return db, nil
}

func (d *db) connectDatabase() error {
	log.Print("Connecting to database...")

	log.Print("Environment variables:")
	log.Printf("DB_USER: %s", viper.GetString("DB_USER"))
	log.Printf("DB_PASSWORD: %s", viper.GetString("DB_PASSWORD"))
	log.Printf("DB_HOST: %s", viper.GetString("DB_HOST"))
	log.Printf("DB_PORT: %s", viper.GetString("DB_PORT"))
	log.Printf("DB_NAME: %s", viper.GetString("DB_NAME"))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return err
	}

	log.Print("Connected to database!")
	d.DB = database
	return nil
}

func (d *db) CreateReport(entry, res, resp, requestKey string, amount float64) error {
	report := models.LLMCallReport{
		Entry:      entry,
		Res:        res,
		Resp:       resp,
		Amount:     amount,
		RequestKey: requestKey,
	}

	result := d.DB.Create(&report)
	return result.Error
}

func (d *db) GetByRequestKey(requestKey string) (*models.LLMCallReport, error) {
	var report models.LLMCallReport
	result := d.DB.Where("request_key = ?", requestKey).First(&report)
	if result.Error != nil {
		return nil, result.Error
	}
	return &report, nil
}
