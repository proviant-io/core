package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/behind-the-fridge/product/internal/db"
	"gitlab.com/behind-the-fridge/product/internal/pkg/category"
	"gitlab.com/behind-the-fridge/product/internal/pkg/list"
	"gitlab.com/behind-the-fridge/product/internal/pkg/product"
	"gitlab.com/behind-the-fridge/product/internal/pkg/service"
	"gitlab.com/behind-the-fridge/product/internal/pkg/stock"
	"strconv"
)

var SqliteLocation = "pantry.db"

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func main() {

	d, err := db.NewSQLite(SqliteLocation)

	if err != nil {
		panic(err)
	}

	productRepo, err := product.Setup(d)

	if err != nil {
		panic(err)
	}

	stockRepo, err := stock.Setup(d)

	if err != nil {
		panic(err)
	}

	categoryRepo, err := category.Setup(d)

	if err != nil {
		panic(err)
	}

	listRepo, err := list.Setup(d)

	if err != nil {
		panic(err)
	}

	relationService := service.NewRelationService(productRepo, listRepo, categoryRepo)

	r := gin.Default()

	// product
	r.GET("/api/v1/product/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {

			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		p, err := relationService.GetProduct(id)

		if err != nil {

			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
			Data:   p,
		}

		c.JSON(response.Status, response)
	})

	r.GET("/api/v1/product/", func(c *gin.Context) {

		p := productRepo.GetAll()

		var models []product.DTO

		for _, model := range p {
			models = append(models, product.ModelToDTO(model))
		}

		response := Response{
			Status: 200,
			Data:   models,
		}

		c.JSON(response.Status, response)
	})

	r.DELETE("/api/v1/product/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		err = productRepo.Delete(id)

		if err != nil {
			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/product/", func(c *gin.Context) {

		dto := product.DTO{}

		err := c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		productRepo.Create(dto)

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/product/:id/", func(c *gin.Context) {

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := product.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		err = productRepo.Update(id, dto)

		if err != nil {
			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	// stock
	r.GET("/api/v1/product/:id/stock/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		s := stockRepo.GetAllByProductId(id)

		var models []stock.DTO

		for _, model := range s {
			models = append(models, stock.ModelToDTO(model))
		}

		response := Response{
			Status: 200,
			Data:   models,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/product/:id/add/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := stock.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto.ProductId = id

		stockRepo.Add(dto)

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/product/:id/consume/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := stock.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto.ProductId = id

		stockRepo.Consume(dto)

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	// category

	r.GET("/api/v1/category/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		model, err := categoryRepo.Get(id)

		if err != nil {
			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
			Data:   category.ModelToDTO(model),
		}

		c.JSON(response.Status, response)

	})

	r.GET("/api/v1/category/", func(c *gin.Context) {

		model := categoryRepo.GetAll()

		var models []category.DTO

		for _, model := range model {
			models = append(models, category.ModelToDTO(model))
		}

		response := Response{
			Status: 200,
			Data:   models,
		}

		c.JSON(response.Status, response)

	})

	r.DELETE("/api/v1/category/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		err = categoryRepo.Delete(id)

		if err != nil {
			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/category/", func(c *gin.Context) {

		dto := category.DTO{}

		err := c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		categoryRepo.Create(dto)

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/category/:id/", func(c *gin.Context) {

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := category.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		err = categoryRepo.Update(id, dto)

		if err != nil {
			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	// list

	r.GET("/api/v1/list/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		model, err := listRepo.Get(id)

		if err != nil {

			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
			Data:   list.ModelToDTO(model),
		}

		c.JSON(response.Status, response)
	})

	r.GET("/api/v1/list/", func(c *gin.Context) {

		model := listRepo.GetAll()

		var models []list.DTO

		for _, model := range model {
			models = append(models, list.ModelToDTO(model))
		}

		response := Response{
			Status: 200,
			Data:   models,
		}

		c.JSON(response.Status, response)

	})

	r.DELETE("/api/v1/list/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		err = listRepo.Delete(id)

		if err != nil {
			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/list/", func(c *gin.Context) {

		dto := list.DTO{}

		err := c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		listRepo.Create(dto)

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/list/:id/", func(c *gin.Context) {

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := list.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: 500,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		err = listRepo.Update(id, dto)

		if err != nil {
			response := Response{
				Status: 404,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: 200,
		}

		c.JSON(response.Status, response)
	})

	r.Run("0.0.0.0:80")
}
