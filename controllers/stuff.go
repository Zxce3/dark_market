package controllers

import (
	"net/http"

	"github.com/crownss/dark_market/config"
	"github.com/crownss/dark_market/models"
	"github.com/gin-gonic/gin"
)

func GetAllStuff(c *gin.Context){
	var allstuf []models.Stuff
	getall := config.DB.Find(&allstuf).RowsAffected
	if getall == 0{
		c.JSON(http.StatusNotFound, models.Response{
			Code: http.StatusNotFound,
			Message: "Data not found !",
			Status: "error",
		})
		return
	}else if(getall != 0){
		c.JSON(http.StatusOK, models.StuffResponseMany{
			Code: http.StatusOK,
			Message: "Data Found !",
			Status: "success",
			Data: allstuf,
		})
		return
	}
}

func GetStuffCode(c *gin.Context){
	var getstuffcode models.Stuff
	if err := config.DB.Where("code = ?", c.Param("code")).First(&getstuffcode).Error; err != nil{
		c.JSON(http.StatusNotFound, models.Response{
			Code: http.StatusNotFound,
			Message: "Data not found !",
			Status: "error",
		})
		return
	}
	c.JSON(http.StatusOK, models.StuffResponseAny{
		Code: http.StatusOK,
		Message: "Data Found !",
		Status: "success",
		Data: getstuffcode,
	})
}

func PostStuff(c *gin.Context){}

func UpdateStuff(c *gin.Context){}

func DeleteStuff(c *gin.Context){}