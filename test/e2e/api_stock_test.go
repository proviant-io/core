package e2e

import (
	"fmt"
	"github.com/proviant-io/core/internal/http"
	"github.com/proviant-io/core/internal/pkg/consumption"
	"github.com/proviant-io/core/internal/pkg/stock"
	"gotest.tools/assert"
	"testing"
	"time"
)

const (
	urlProduct      = "/api/v1/product/"
	urlList         = "/api/v1/list/"
	urlCategory     = "/api/v1/category/"
	urlStock        = "/api/v1/product/1/stock/"
	urlStockAdd     = "/api/v1/product/1/add/"
	urlStockConsume = "/api/v1/product/1/consume/"
	urlStockWithId  = "/api/v1/product/1/stock/%d/"
)

func TestApiStock(t *testing.T) {
	t.Skip("Skipping stock tests as they are not polished yet (need regex comparison)")
	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	var actual string
	var expected http.Response

	execSuitStep(t, "stock: check db empty", urlProduct, GET, "", apiResponse{
		Status: 200,
		Data:   []interface{}{},
		Error:  "",
	})

	execSuitStep(t, "stock: create list", urlList, POST, `{"title": "Fridge"}`, apiResponse{
		Status: 201,
		Data: map[string]interface{}{
			"id":    1,
			"title": "Fridge",
		},
		Error: "",
	})

	fmt.Print("stock: create category")
	actual = postRequest(generateApiUrl(urlCategory), []byte(`{"title": "Drinks"}`))

	fmt.Print("stock: create product")
	actual = postRequest(generateApiUrl(urlProduct),
		[]byte(`{"title":"Milk Shake", "description":  "Milk Shake", "link":  "https://test.com/test", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))

	fmt.Print("stock: add 5")
	actual = postRequest(generateApiUrl(urlStockAdd), []byte(`{"quantity":  5, "expire":  1609458959}`))
	expected = http.Response{
		Status: 201,
		Data: stock.DTO{
			Id:        1,
			ProductId: 1,
			Quantity:  5,
			Expire:    1609458959,
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: add 3")
	actual = postRequest(generateApiUrl(urlStockAdd), []byte(`{"quantity":  3, "expire":  1609502159}`))
	expected = http.Response{
		Status: 201,
		Data: stock.DTO{
			Id:        2,
			ProductId: 1,
			Quantity:  3,
			Expire:    1609502159,
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: add 3")
	actual = postRequest(generateApiUrl(urlStockAdd), []byte(`{"quantity":  3, "expire":  1609502259}`))
	expected = http.Response{
		Status: 201,
		Data: stock.DTO{
			Id:        3,
			ProductId: 1,
			Quantity:  3,
			Expire:    1609502259,
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: get")
	actual = getRequest(generateApiUrl(urlStock))
	expected = http.Response{
		Status: 201,
		Data: []stock.DTO{
			{
				Id:        1,
				ProductId: 1,
				Quantity:  5,
				Expire:    1609458959,
			},
			{
				Id:        2,
				ProductId: 1,
				Quantity:  3,
				Expire:    1609502159,
			},
			{
				Id:        3,
				ProductId: 1,
				Quantity:  3,
				Expire:    1609502259,
			},
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest(generateApiUrl(urlStockConsume), []byte(`{"quantity":  3}`))
	expected = http.Response{
		Status: 201,
		Data: []stock.DTO{
			{
				Id:        1,
				ProductId: 1,
				Quantity:  2,
				Expire:    1609458959,
			},
			{
				Id:        2,
				ProductId: 1,
				Quantity:  3,
				Expire:    1609502159,
			},
			{
				Id:        3,
				ProductId: 1,
				Quantity:  3,
				Expire:    1609502259,
			},
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after consume")
	actual = getRequest(generateApiUrl(urlStock))
	expected = http.Response{
		Status: 201,
		Data: []stock.DTO{
			{
				Id:        1,
				ProductId: 1,
				Quantity:  2,
				Expire:    1609458959,
			},
			{
				Id:        2,
				ProductId: 1,
				Quantity:  3,
				Expire:    1609502159,
			},
			{
				Id:        3,
				ProductId: 1,
				Quantity:  3,
				Expire:    1609502259,
			},
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest(generateApiUrl(urlStockConsume), []byte(`{"quantity":  3}`))
	expected = http.Response{
		Status: 201,
		Data: struct {
			Stock           []stock.DTO     `json:"stock"`
			ConsumedLogItem consumption.DTO `json:"consumed_log_item"`
		}{
			Stock: []stock.DTO{
				{
					Id:        2,
					ProductId: 1,
					Quantity:  2,
					Expire:    1609502159,
				},
				{
					Id:        3,
					ProductId: 1,
					Quantity:  3,
					Expire:    1609502259,
				},
			},
			ConsumedLogItem: consumption.DTO{
				Id:         1,
				ProductId:  1,
				Quantity:   3,
				ConsumedAt: 0,
				UserId:     0,
				AccountId:  0,
			},
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after second consume")
	actual = getRequest(generateApiUrl(urlStock))
	expected = http.Response{
		Status: 201,
		Data: []stock.DTO{
			{
				Id:        2,
				ProductId: 1,
				Quantity:  2,
				Expire:    1609502159,
			},
			{
				Id:        3,
				ProductId: 1,
				Quantity:  3,
				Expire:    1609502259,
			},
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

	fmt.Print("stock: delete with id 3")
	actual = deleteRequest(generateApiUrl(url(urlStockWithId, 3)))
	expected = http.Response{
		Status: 200,
		Data: []stock.DTO{
			{
				Id:        2,
				ProductId: 1,
				Quantity:  2,
				Expire:    1609502159,
			},
		},
		Error: "",
	}
	assert.Equal(t, toJson(t, expected), actual)
	fmt.Println(" OK")

}
