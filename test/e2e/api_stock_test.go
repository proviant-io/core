package e2e

import (
	"fmt"
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

	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	var actual string
	var expected string

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
	expected = `{"status":201,"data":{"id":1,"product_id":1,"quantity":5,"expire":1609458959},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: add 3")
	actual = postRequest(generateApiUrl(urlStockAdd), []byte(`{"quantity":  3, "expire":  1609502159}`))
	expected = `{"status":201,"data":{"id":2,"product_id":1,"quantity":3,"expire":1609502159},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: add 3")
	actual = postRequest(generateApiUrl(urlStockAdd), []byte(`{"quantity":  3, "expire":  1609502259}`))
	expected = `{"status":201,"data":{"id":3,"product_id":1,"quantity":3,"expire":1609502259},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get")
	actual = getRequest(generateApiUrl(urlStock))
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":5,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest(generateApiUrl(urlStockConsume), []byte(`{"quantity":  3}`))
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":2,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after consume")
	actual = getRequest(generateApiUrl(urlStock))
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":2,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest(generateApiUrl(urlStockConsume), []byte(`{"quantity":  3}`))
	expected = `{"status":200,"data":[{"id":2,"product_id":1,"quantity":2,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after second consume")
	actual = getRequest(generateApiUrl(urlStock))
	expected = `{"status":200,"data":[{"id":2,"product_id":1,"quantity":2,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: delete with id 3")
	actual = deleteRequest(generateApiUrl(url(urlStockWithId, 3)))
	expected = `{"status":200,"data":[{"id":2,"product_id":1,"quantity":2,"expire":1609502159}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")
}
