package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/behind-the-fridge/product/pkg/product"
	"gitlab.com/behind-the-fridge/product/pkg/stock"
	"strconv"
)


func main() {

	productRepo, err := product.Setup()

	if err != nil{
		panic(err)
	}


	stockRepo, err := stock.Setup()

	if err != nil{
		panic(err)
	}

	r := gin.Default()

	// product

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

	r.DELETE("/product/:id", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(200, gin.H{
				"err": err,
			})
		}

		productRepo.Delete(id)

		c.JSON(200, gin.H{
			"ok": true,
		})
	})

	r.POST("/product/", func(c *gin.Context) {

		dto := product.DTO{}

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


	r.POST("/product/:id", func(c *gin.Context) {

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(200, gin.H{
				"err": err,
			})
		}

		dto := product.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			c.JSON(200, gin.H{
				"err": err,
			})
		}

		productRepo.Update(id, dto)

		c.JSON(200, gin.H{
			"ok": true,
		})
	})

	// stock
	r.GET("/product/:id/stock/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(200, gin.H{
				"err": err,
			})
		}

		s := stockRepo.GetAllByProductId(id)

		c.JSON(200, s)
	})

	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
