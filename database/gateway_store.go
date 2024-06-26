package database

import (
	"math/rand"
	"net/http"
	"qpay/models"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// handle returns
func PostGateway(gateway models.Gateway) (int, error) {
	db := NewGormPostgres()
	gateway.CreatedAt = time.Now()
	gateway.UrlAddress = RandStringBytes(50)
	result := db.Create(&gateway)
	if result.Error != nil {
		return http.StatusBadRequest, result.Error
	}
	return http.StatusOK, nil
}

func PostCustomUrlGateway(getaway models.Gateway) (int, error) {
	db := NewGormPostgres()
	getaway.CreatedAt = time.Now()
	result := db.Create(&getaway)

	if result.Error != nil {
		return http.StatusBadRequest, result.Error
	}
	return http.StatusOK, nil
}

func UpdateGatewayAccount(gatewayId int, gatewayAccountNumber string) (int, error) {
	db := NewGormPostgres()
	err := db.Model(&models.Gateway{}).Where("id = ?", gatewayId).Update("account_id", gatewayAccountNumber)
	//err := db.Model(db.First(gatewayId).Update("account_id", gatewayAccountNumber))
	if err.Error != nil {
		return http.StatusBadRequest, err.Error
	}
	return http.StatusOK, nil
}
