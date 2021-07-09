package e2e

import (
	"fmt"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestApiStock(t *testing.T) {

	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	fmt.Print("stock: check db empty")
	actual := getRequest("http://localhost:8081/api/v1/product/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: create list")
	actual = postRequest("http://localhost:8081/api/v1/list/", []byte(`{"title": "Fridge"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: create category")
	actual = postRequest("http://localhost:8081/api/v1/category/", []byte(`{"title": "Drinks"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: create product")
	actual = postRequest("http://localhost:8081/api/v1/product/",
		[]byte(`{"title":"Milk Shake", "description":  "Milk Shake", "link":  "https://test.com/test", "image":  "https://inage.com/1.jpg", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = `{"status":201,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"},"stock":0},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: add 5")
	actual = postRequest("http://localhost:8081/api/v1/product/1/add/", []byte(`{"quantity":  5, "expire":  1609458959}`))
	expected = `{"status":201,"data":{"id":1,"product_id":1,"quantity":5,"expire":1609458959},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: add 3")
	actual = postRequest("http://localhost:8081/api/v1/product/1/add/", []byte(`{"quantity":  3, "expire":  1609502159}`))
	expected = `{"status":201,"data":{"id":2,"product_id":1,"quantity":3,"expire":1609502159},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: add 3")
	actual = postRequest("http://localhost:8081/api/v1/product/1/add/", []byte(`{"quantity":  3, "expire":  1609502259}`))
	expected = `{"status":201,"data":{"id":3,"product_id":1,"quantity":3,"expire":1609502259},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get")
	actual = getRequest("http://localhost:8081/api/v1/product/1/stock/")
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":5,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest("http://localhost:8081/api/v1/product/1/consume/", []byte(`{"quantity":  3}`))
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":2,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after consume")
	actual = getRequest("http://localhost:8081/api/v1/product/1/stock/")
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":2,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest("http://localhost:8081/api/v1/product/1/consume/", []byte(`{"quantity":  3}`))
	expected = `{"status":200,"data":[{"id":2,"product_id":1,"quantity":2,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after second consume")
	actual = getRequest("http://localhost:8081/api/v1/product/1/stock/")
	expected = `{"status":200,"data":[{"id":2,"product_id":1,"quantity":2,"expire":1609502159},{"id":3,"product_id":1,"quantity":3,"expire":1609502259}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: delete with id 3")
	actual = deleteRequest("http://localhost:8081/api/v1/product/1/stock/3/")
	expected = `{"status":200,"data":[{"id":2,"product_id":1,"quantity":2,"expire":1609502159}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")
}


