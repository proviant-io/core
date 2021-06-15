package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/behind-the-fridge/product/pkg"
	"strconv"
)


func main() {

	productRepo, err := pkg.Setup()

	if err != nil{
		panic(err)
	}

	r := gin.Default()
	r.GET("/product/:id", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(200, gin.H{
				"err": err,
			})
		}

		p := productRepo.Get(id)

		c.JSON(200, p)
	})


	r.POST("/product/", func(c *gin.Context) {

		dto := pkg.ProductDTO{}

		err := c.BindJSON(&dto)

		if err != nil {
			c.JSON(200, gin.H{
				"err": err,
			})
		}

		productRepo.Create(dto)

		c.JSON(200, gin.H{
			"ok": true,
		})
	})

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
