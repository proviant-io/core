package main

import (
	"github.com/brushknight/proviant/internal/db"
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/brushknight/proviant/internal/pkg/product_category"
	"github.com/brushknight/proviant/internal/pkg/service"
	"github.com/brushknight/proviant/internal/pkg/stock"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
)

var SqliteLocation = "pantry.db"

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

const (
	ResponseCodeOk      = 200
	ResponseCodeCreated = 201
	BadRequest          = 400
	InternalServerError = 500
)

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

	productCategoryRepo, err := product_category.Setup(d)

	if err != nil {
		panic(err)
	}

	relationService := service.NewRelationService(productRepo, listRepo, categoryRepo, stockRepo, productCategoryRepo)

	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTION"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	r.Use(cors.New(config))

	// product
	r.GET("/api/v1/product/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {

			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		p, customErr := relationService.GetProduct(id)

		if customErr != nil {

			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   p,
		}

		c.JSON(response.Status, response)
	})

	r.GET("/api/v1/product/", func(c *gin.Context) {

		queryRaw := c.Request.URL.Query()

		var query *product.Query

		if len(queryRaw) > 0 {
			query = &product.Query{}

			if listFilterRaw, ok := queryRaw["list"]; ok {
				if len(listFilterRaw) > 0 {
					listFilter, err := strconv.Atoi(listFilterRaw[0])

					if err != nil {
						response := Response{
							Status: BadRequest,
							Error:  err.Error(),
						}

						c.JSON(response.Status, response)
						return
					}

					query.List = listFilter
				}
			}
			if categoryFilterRaw, ok := queryRaw["category"]; ok {
				if len(categoryFilterRaw) > 0 {
					categoryFilter, err := strconv.Atoi(categoryFilterRaw[0])

					if err != nil {
						response := Response{
							Status: BadRequest,
							Error:  err.Error(),
						}

						c.JSON(response.Status, response)
						return
					}

					query.Category = categoryFilter
				}
			}
		}

		dtos := relationService.GetAllProducts(query)

		response := Response{
			Status: ResponseCodeOk,
			Data:   dtos,
		}

		c.JSON(response.Status, response)
	})

	r.DELETE("/api/v1/product/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		customErr := productRepo.Delete(id)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/product/", func(c *gin.Context) {

		dto := product.DTO{}

		err := c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		productDto, customErr := relationService.CreateProduct(dto)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeCreated,
			Data:   productDto,
		}

		c.JSON(response.Status, response)
	})

	r.PUT("/api/v1/product/:id/", func(c *gin.Context) {

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := product.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto.Id = id

		productDTO, customErr := relationService.UpdateProduct(dto)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   productDTO,
		}

		c.JSON(response.Status, response)
	})

	// stock
	r.GET("/api/v1/product/:id/stock/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
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
			Status: ResponseCodeOk,
			Data:   models,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/product/:id/add/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := stock.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto.ProductId = id

		model, customErr := relationService.AddStock(dto)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeCreated,
			Data:   stock.ModelToDTO(model),
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/product/:id/consume/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := stock.ConsumeDTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto.ProductId = id

		customErr := relationService.ConsumeStock(dto)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
		}

		c.JSON(response.Status, response)
	})

	// category

	r.GET("/api/v1/category/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		model, customErr := categoryRepo.Get(id)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   category.ModelToDTO(model),
		}

		c.JSON(response.Status, response)

	})

	r.GET("/api/v1/category/", func(c *gin.Context) {

		model := categoryRepo.GetAll()

		models := []category.DTO{}

		for _, model := range model {
			models = append(models, category.ModelToDTO(model))
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   models,
		}

		c.JSON(response.Status, response)

	})

	r.DELETE("/api/v1/category/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		customErr := categoryRepo.Delete(id)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/category/", func(c *gin.Context) {

		dto := category.DTO{}

		err := c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		categoryModel := categoryRepo.Create(dto)

		response := Response{
			Status: ResponseCodeCreated,
			Data:   category.ModelToDTO(categoryModel),
		}

		c.JSON(response.Status, response)
	})

	r.PUT("/api/v1/category/:id/", func(c *gin.Context) {

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := category.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		categoryModel, customErr := categoryRepo.Update(id, dto)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   category.ModelToDTO(categoryModel),
		}

		c.JSON(response.Status, response)
	})

	// list

	r.GET("/api/v1/list/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		model, customErr := listRepo.Get(id)

		if customErr != nil {

			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   list.ModelToDTO(model),
		}

		c.JSON(response.Status, response)
	})

	r.GET("/api/v1/list/", func(c *gin.Context) {

		model := listRepo.GetAll()

		models := []list.DTO{}

		for _, model := range model {
			models = append(models, list.ModelToDTO(model))
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   models,
		}

		c.JSON(response.Status, response)

	})

	r.DELETE("/api/v1/list/:id/", func(c *gin.Context) {
		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		customErr := listRepo.Delete(id)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
		}

		c.JSON(response.Status, response)
	})

	r.POST("/api/v1/list/", func(c *gin.Context) {

		dto := list.DTO{}

		err := c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		listModel := listRepo.Create(dto)

		response := Response{
			Status: ResponseCodeCreated,
			Data:   list.ModelToDTO(listModel),
		}

		c.JSON(response.Status, response)
	})

	r.PUT("/api/v1/list/:id/", func(c *gin.Context) {

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		dto := list.DTO{}

		err = c.BindJSON(&dto)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		listModel, customErr := listRepo.Update(id, dto)

		if customErr != nil {
			response := Response{
				Status: customErr.Code(),
				Error:  customErr.Error(),
			}

			c.JSON(response.Status, response)
			return
		}

		response := Response{
			Status: ResponseCodeOk,
			Data:   list.ModelToDTO(listModel),
		}

		c.JSON(response.Status, response)
	})

	// static
	r.Static("/static/", "./public/")
	r.StaticFile("/", "./public/index.html")
	//r.StaticFile("/product*", "./public/index.html")

	r.Run("0.0.0.0:80")
}
