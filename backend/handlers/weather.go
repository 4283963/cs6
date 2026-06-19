package handlers

import (
	"net/http"
	"streetlight-controller/weather"

	"github.com/gin-gonic/gin"
)

type UpdateWeatherRequest struct {
	Illuminance float64 `json:"illuminance" binding:"required,min=0"`
	Condition   string  `json:"condition" binding:"required"`
}

func GetWeather(c *gin.Context) {
	status := weather.GlobalWeatherService.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"data": status,
	})
}

func UpdateWeather(c *gin.Context) {
	var req UpdateWeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	status := weather.GlobalWeatherService.UpdateStatus(req.Illuminance, req.Condition)
	c.JSON(http.StatusOK, gin.H{
		"message": "天气状态更新成功",
		"data":    status,
	})
}

func SimulateStorm(c *gin.Context) {
	weather.GlobalWeatherService.SimulateStorm()
	c.JSON(http.StatusOK, gin.H{
		"message": "已模拟暴雨天气，光照度 15 lux",
		"data":    weather.GlobalWeatherService.GetStatus(),
	})
}

func SimulateNormal(c *gin.Context) {
	weather.GlobalWeatherService.SimulateNormal()
	c.JSON(http.StatusOK, gin.H{
		"message": "已恢复正常天气，光照度 1000 lux",
		"data":    weather.GlobalWeatherService.GetStatus(),
	})
}

func GetWeatherStream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	currentStatus := weather.GlobalWeatherService.GetStatus()
	c.SSEvent("message", currentStatus)
	flusher.Flush()

	listener := weather.GlobalWeatherService.GetListener()

	for {
		select {
		case status := <-listener:
			c.SSEvent("message", status)
			flusher.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
