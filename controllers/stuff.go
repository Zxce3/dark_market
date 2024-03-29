package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/crownss/dark_market/config"
	"github.com/crownss/dark_market/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAllStuff(c *gin.Context) {
	var allstuf []models.Stuff
	getall := config.DB.Find(&allstuf).RowsAffected
	if getall != 0 {
		c.JSON(http.StatusOK, models.StuffResponseMany{
			Code:    http.StatusOK,
			Message: "Data Found !",
			Status:  "success",
			Data:    allstuf,
		})
		return
	}
	c.JSON(http.StatusNotFound, models.Response{
		Code:    http.StatusNotFound,
		Message: "Not Found !",
		Status:  "error",
	})
}

func GetStuffCode(c *gin.Context) {
	var getstuffcode models.Stuff
	if err := config.DB.Where("code = ?", c.Param("code")).First(&getstuffcode).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Code:    http.StatusNotFound,
			Message: "Data not found !",
			Status:  "error",
		})
		return
	}
	c.JSON(http.StatusOK, models.StuffResponseAny{
		Code:    http.StatusOK,
		Message: "Data Found !",
		Status:  "success",
		Data:    getstuffcode,
	})
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func _RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func PostStuff(c *gin.Context) {
	var UserInDB models.Users
	session := sessions.Default(c)
	user := session.Get(userkey)
	check_admin := config.DB.Where("username = ?", user).Where("is_admin = ?", true).Find(&UserInDB).RowsAffected
	fmt.Println(user, check_admin)
	check_user := config.DB.Where("username = ?", user).Where("is_active = ?", true).Where("is_admin = ?", false).Find(&UserInDB).RowsAffected
	if check_admin != 0 {
		var input models.Stuff
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Status:  "error",
			})
			return
		}
		rand.Seed(time.Now().UnixNano())
		res := models.Stuff{
			Code:  _RandomString(5),
			Img:   input.Img,
			Title: input.Title,
			Desc:  input.Desc,
			Stock: input.Stock,
			Price: input.Price,
		}
		if e := config.DB.Create(&res).Error; e != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    http.StatusInternalServerError,
				Message: e.Error(),
				Status:  "error",
			})
			return
		}
		c.JSON(http.StatusOK, &res)
		return
	}
	if check_user != 0 {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    http.StatusForbidden,
			Message: "You not have access",
			Status:  "error",
		})
		return
	} else if check_user == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "You are not logged in",
			Status:  "error",
		})
		return
	}

}

func UpdateStuff(c *gin.Context) {
	var UserInDB models.Users
	session := sessions.Default(c)
	user := session.Get(userkey)
	check_admin := config.DB.Where("username = ?", user).Where("is_admin = ?", true).Find(&UserInDB).RowsAffected
	check_user := config.DB.Where("username = ?", user).Where("is_active = ?", true).Where("is_admin = ?", false).Find(&UserInDB).RowsAffected
	if check_admin != 0 {
		var GetCode models.Stuff
		if err := config.DB.Where("code = ?", c.Param("code")).First(&GetCode).Error; err != nil {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    http.StatusNotFound,
				Message: err.Error(),
				Status:  "error",
			})
			return
		}
		var input models.Stuff
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Status:  "error",
			})
			return
		}
		update := make(map[string]interface{})
		update["img"] = input.Img
		update["title"] = input.Title
		update["desc"] = input.Desc
		update["stock"] = input.Stock
		update["price"] = input.Price
		if err := config.DB.Where("code = ?", c.Param("code")).First(&GetCode).Updates(update).Error; err != nil {
			c.JSON(http.StatusNotFound, &models.Response{
				Code:    http.StatusNotFound,
				Message: err.Error(),
				Status:  "error",
			})
			return
		}
		c.JSON(http.StatusAccepted, &models.Response{
			Code:    http.StatusAccepted,
			Message: "Succesfuly Change",
			Status:  "success",
		})
		return
	}
	if check_user != 0 {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    http.StatusForbidden,
			Message: "You not have access",
			Status:  "error",
		})
		return
	} else if check_user == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "You are not logged in",
			Status:  "error",
		})
		return
	}
}

func DeleteStuff(c *gin.Context) {
	var UserInDB models.Users
	session := sessions.Default(c)
	user := session.Get(userkey)
	check_admin := config.DB.Where("username = ?", user).Where("is_admin = ?", true).Find(&UserInDB).RowsAffected
	check_user := config.DB.Where("username = ?", user).Where("is_active = ?", true).Where("is_admin = ?", false).Find(&UserInDB).RowsAffected
	if check_admin != 0 {
		var GetCode models.Stuff
		if err := config.DB.Where("code = ?", c.Param("code")).First(&GetCode).Error; err != nil {
			c.JSON(http.StatusNotFound, models.Response{
				Code:    http.StatusNotFound,
				Message: err.Error(),
				Status:  "error",
			})
			return
		}
		config.DB.Delete(&GetCode)
		c.JSON(http.StatusOK, models.Response{
			Code:    http.StatusOK,
			Message: "Succesfuly Delete",
			Status:  "success",
		})
		return
	}
	if check_user != 0 {
		c.JSON(http.StatusForbidden, models.Response{
			Code:    http.StatusForbidden,
			Message: "You not have access",
			Status:  "error",
		})
		return
	} else if check_user == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "You are not logged in",
			Status:  "error",
		})
		return
	}
}
