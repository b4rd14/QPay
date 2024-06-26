package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"qpay/database"
	"qpay/models"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

var users []models.User
var service UserInterface
var serverDoOnce sync.Once

func startServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/signup", RegisterHandler(UserInterfaceService{}))
	e.GET("/home", Authentication, AuthMiddleware)
	err := e.Start("localhost:6060")
	if err != nil {
		return
	}
}

func setupServer() {
	serverDoOnce.Do(func() {
		go startServer()
		time.Sleep(300 * time.Millisecond)
	})
}

func TestDBConnection(t *testing.T) {
	db := database.NewGormPostgres()
	err := db.Exec("SELECT 1").Error
	assert.NoError(t, err)

}

func TestCreateUser(t *testing.T) {
	user := models.User{
		IsCompany: false,
		Name:      "Qpay",
		Email:     "Qpayfake@gmail.com",
		Password:  "password",
	}
	err := service.RegisterUser(user)
	assert.NoError(t, err)
}

func TestCreateHandler2(t *testing.T) {
	user := models.User{
		IsCompany: false,
		Name:      "Qpay",
		Email:     "Qpay@gmail.com",
		Password:  "password",
	}
	err := service.RegisterUser(user)
	assert.Error(t, err)
}

func TestCreateHandler(t *testing.T) {
	setupServer()
	user := models.User{
		IsCompany: false,
		Name:      "arshia1235",
		Email:     "arshia1235@example.com",
		Password:  "password",
	}
	body, _ := json.Marshal(user)

	resp, err := http.DefaultClient.Post("http://127.0.0.1:6060/signup", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
}
