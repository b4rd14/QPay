package services

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"qpay/database"
	"qpay/models"
)

type AdminInterface interface {
	BlockAllGateWay(user models.User, BlockType string) error
	BlockOneGateWay(gateway models.Gateway, BlockType string) error
	UnblockGateWay(gateway models.Gateway) error
}

type AdminInterfaceService struct{}

func (a *AdminInterfaceService) BlockAllGateWay(user models.User, BlockType string) error {
	db := database.NewGormPostgres()
	var gateways []models.Gateway
	err := db.Where("user_id = ?", user.ID).Find(&gateways).Error
	if err != nil {
		return err
	}
	for _, gateway := range gateways {
		if BlockType == "block" {
			if gateway.Blocked == true {
				return echo.NewHTTPError(http.StatusBadRequest, "gateway is already blocked")
			}
			err := db.Model(&gateway).Update("blocked", true).Error
			if err != nil {
				return err
			}
		} else if BlockType == "alwaysBlock" {
			if gateway.AlwaysBlocked == true {
				return echo.NewHTTPError(http.StatusBadRequest, "gateway is already always blocked")
			}
			err := db.Model(&gateway).Update("always_blocked", true).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *AdminInterfaceService) BlockOneGateWay(gateway models.Gateway, BlockType string) error {
	db := database.NewGormPostgres()
	if BlockType == "block" {
		if gateway.Blocked == true {
			return echo.NewHTTPError(http.StatusBadRequest, "gateway is already blocked")
		}
		err := db.Model(&gateway).Update("blocked", true).Error
		if err != nil {
			return err
		}
	} else if BlockType == "alwaysBlock" {
		if gateway.AlwaysBlocked == true {
			return echo.NewHTTPError(http.StatusBadRequest, "gateway is already always blocked")
		}
		err := db.Model(&gateway).Update("always_blocked", true).Error
		if err != nil {
			return err
		}
	}
	return nil

}

func (a *AdminInterfaceService) UnblockGateWay(gateway models.Gateway) error {
	db := database.NewGormPostgres()
	if gateway.Blocked == false {
		return echo.NewHTTPError(http.StatusBadRequest, "gateway is not blocked")
	}
	err := db.Model(&gateway).Update("blocked", false).Error
	if err != nil {
		return err
	}
	return nil
}

func BlockOneGateWayHandler(service AdminInterfaceService) echo.HandlerFunc {
	return func(c echo.Context) error {
		BlockType := c.QueryParam("blockType")
		var gateway models.Gateway
		err := c.Bind(&gateway)
		if err != nil {
			return err
		}
		err = service.BlockOneGateWay(gateway, BlockType)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "gateway blocked successfully ",
		})
	}
}

func BlockAllGateWayHandler(service AdminInterfaceService) echo.HandlerFunc {
	return func(c echo.Context) error {
		BlockType := c.QueryParam("blockType")
		var user models.User
		err := c.Bind(&user)
		if err != nil {
			return err
		}
		err = service.BlockAllGateWay(user, BlockType)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"message": "all gateways blocked successfully ",
		})
	}
}

func AdminRoutes(server *echo.Echo) {
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.POST("/admin/login", LoginAdminHandler(AdminInterfaceService{}))
	server.POST("/admin/blockOneGateway", BlockOneGateWayHandler(AdminInterfaceService{}))
	server.POST("/admin/blockAllGateways", BlockAllGateWayHandler(AdminInterfaceService{}))
	server.GET("/admin", Authentication, AuthMiddleware)
}
