package e2e

import (
	"fmt"
	"github.com/proviant-io/core/internal/http"
	"github.com/proviant-io/core/internal/pkg/category"
	"github.com/proviant-io/core/internal/pkg/list"
	"github.com/proviant-io/core/internal/pkg/product"
	"github.com/shopspring/decimal"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestApiProduct(t *testing.T) {

	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	fmt.Print("product: check db empty")
	actual := getRequest(generateApiUrl("/api/v1/product/"))
	expected := http.Response{
		Status: 200,
		Data:   []interface{}{},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: create list")
	actual = postRequest(generateApiUrl("/api/v1/list/"), []byte(`{"title": "Fridge"}`))
	expected = http.Response{
		Status: 201,
		Data:   list.DTO{
			Id:    1,
			Title: "Fridge",
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: create category")
	actual = postRequest(generateApiUrl("/api/v1/category/"), []byte(`{"title": "Drinks"}`))
	expected = http.Response{
		Status: 201,
		Data:   category.DTO{
			Id:    1,
			Title: "Drinks",
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: create")
	actual = postRequest(generateApiUrl("/api/v1/product/"),
		[]byte(`{"title":"Milk Shake", "description":  "Milk Shake", "link":  "https://test.com/test", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = http.Response{
		Status: 201,
		Data:   product.DTO{
			Id:          1,
			Title:       "Milk Shake",
			Description: "Milk Shake",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{{
				Id: 1,
				Title: "Drinks",
			}},
			ListId:      1,
			List:        list.DTO{
				Id:    1,
				Title: "Fridge",
			},
			Stock:       0,
			Price:       decimal.New(0,0),
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: get")
	actual = getRequest(generateApiUrl("/api/v1/product/1/"))
	expected = http.Response{
		Status: 200,
		Data:   product.DTO{
			Id:          1,
			Title:       "Milk Shake",
			Description: "Milk Shake",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{{
				Id: 1,
				Title: "Drinks",
			}},
			ListId:      1,
			List:        list.DTO{
				Id:    1,
				Title: "Fridge",
			},
			Stock:       0,
			Price:       decimal.New(0,0),
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: update")
	actual = putRequest(generateApiUrl("/api/v1/product/1/"),
		[]byte(`{"title":"Milk Shake 2", "description":  "Milk Shake 2", "link":  "https://test.com/test", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = http.Response{
		Status: 200,
		Data:   product.DTO{
			Id:          1,
			Title:       "Milk Shake 2",
			Description: "Milk Shake 2",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{{
				Id: 1,
				Title: "Drinks",
			}},
			ListId:      1,
			List:        list.DTO{
				Id:    1,
				Title: "Fridge",
			},
			Stock:       0,
			Price:       decimal.New(0,0),
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: get updated")
	actual = getRequest(generateApiUrl("/api/v1/product/1/"))
	expected = http.Response{
		Status: 200,
		Data:   product.DTO{
			Id:          1,
			Title:       "Milk Shake 2",
			Description: "Milk Shake 2",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{{
				Id: 1,
				Title: "Drinks",
			}},
			ListId:      1,
			List:        list.DTO{
				Id:    1,
				Title: "Fridge",
			},
			Stock:       0,
			Price:       decimal.New(0,0),
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: get all")
	actual = getRequest(generateApiUrl("/api/v1/product/"))
	expected = http.Response{
		Status: 200,
		Data:   []product.DTO{{
			Id:          1,
			Title:       "Milk Shake 2",
			Description: "Milk Shake 2",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{},
			ListId:      1,
			List:        nil,
			Stock:       0,
			Price:       decimal.New(0,0),
		}},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: delete")
	actual = deleteRequest(generateApiUrl("/api/v1/product/1/"))
	expected = http.Response{
		Status: 200,
		Data:   nil,
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: deleted not found")
	actual = getRequest(generateApiUrl("/api/v1/product/1/"))
	expected = http.Response{
		Status: 404,
		Data:   nil,
		Error:  "product with id 1 not found",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")
}

func TestApiProductFilter(t *testing.T){
	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	fmt.Print("product: check db empty")
	actual := getRequest(generateApiUrl("/api/v1/product/"))
	expected := http.Response{
		Status: 200,
		Data:   []interface{}{},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: create list")
	actual = postRequest(generateApiUrl("/api/v1/list/"), []byte(`{"title": "Fridge"}`))
	expected = http.Response{
		Status: 201,
		Data:   list.DTO{
			Id:    1,
			Title: "Fridge",
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: create category")
	actual = postRequest(generateApiUrl("/api/v1/category/"), []byte(`{"title": "Drinks"}`))
	expected = http.Response{
		Status: 201,
		Data:   category.DTO{
			Id:    1,
			Title: "Drinks",
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: create")
	actual = postRequest(generateApiUrl("/api/v1/product/"),
		[]byte(`{"title":"Milk Shake", "description":  "Milk Shake", "link":  "https://test.com/test", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = http.Response{
		Status: 201,
		Data:   product.DTO{
			Id:          1,
			Title:       "Milk Shake",
			Description: "Milk Shake",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{{
				Id: 1,
				Title: "Drinks",
			}},
			ListId:      1,
			List:        list.DTO{
				Id:    1,
				Title: "Fridge",
			},
			Stock:       0,
			Price:       decimal.New(0,0),
		},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?list=1")
	actual = getRequest(generateApiUrl("/api/v1/product/?list=1"))
	expected = http.Response{
		Status: 200,
		Data:   []product.DTO{{
			Id:          1,
			Title:       "Milk Shake",
			Description: "Milk Shake",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{},
			ListId:      1,
			List:        nil,
			Stock:       0,
			Price:       decimal.New(0,0),
		}},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?list=2")
	actual = getRequest(generateApiUrl("/api/v1/product/?list=2"))
	expected = http.Response{
		Status: 200,
		Data:   []interface{}{},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?category=1")
	actual = getRequest(generateApiUrl("/api/v1/product/?category=1"))
	expected = http.Response{
		Status: 200,
		Data:   []product.DTO{{
			Id:          1,
			Title:       "Milk Shake",
			Description: "Milk Shake",
			Link:        "https://test.com/test",
			Image:       "",
			Barcode:     "1234567890Z",
			CategoryIds: []int{1},
			Categories:  []category.DTO{},
			ListId:      1,
			List:        nil,
			Stock:       0,
			Price:       decimal.New(0,0),
		}},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?list=2")
	actual = getRequest(generateApiUrl("/api/v1/product/?list=2"))
	expected = http.Response{
		Status: 200,
		Data:   []interface{}{},
		Error:  "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")
}

