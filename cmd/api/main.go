package main

import (
	"net/http"
	"tax-calculator/internal/entity"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/order", func(c *gin.Context) {
		order, _ := entity.NewOrder("1", 10, 1)
		err := order.CalculateFinalPrice()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	})
	r.Run()

}
