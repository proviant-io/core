package e2e

import (
	"fmt"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestApiProduct(t *testing.T) {

	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	fmt.Print("product: check db empty")
	actual := getRequest("http://localhost:8081/api/v1/product/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: create list")
	actual = postRequest("http://localhost:8081/api/v1/list/", []byte(`{"title": "Fridge"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: create category")
	actual = postRequest("http://localhost:8081/api/v1/category/", []byte(`{"title": "Drinks"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: create")
	actual = postRequest("http://localhost:8081/api/v1/product/",
		[]byte(`{"title":"Milk Shake", "description":  "Milk Shake", "link":  "https://test.com/test", "image":  "https://inage.com/1.jpg", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = `{"status":201,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"},"stock":0},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get")
	actual = getRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"},"stock":0},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: update")
	actual = putRequest("http://localhost:8081/api/v1/product/1/",
		[]byte(`{"title":"Milk Shake 2", "description":  "Milk Shake 2", "link":  "https://test.com/test", "image":  "https://inage.com/1.jpg", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"},"stock":0},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get updated")
	actual = getRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"},"stock":0},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get all")
	actual = getRequest("http://localhost:8081/api/v1/product/")
	expected = `{"status":200,"data":[{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[],"list_id":1,"list":null,"stock":0}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: delete")
	actual = deleteRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: deleted not found")
	actual = getRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":404,"data":null,"error":"product with id 1 not found"}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")
}

func TestApiProductFilter(t *testing.T){
	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	fmt.Print("product: check db empty")
	actual := getRequest("http://localhost:8081/api/v1/product/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: create list")
	actual = postRequest("http://localhost:8081/api/v1/list/", []byte(`{"title": "Fridge"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: create category")
	actual = postRequest("http://localhost:8081/api/v1/category/", []byte(`{"title": "Drinks"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: create")
	actual = postRequest("http://localhost:8081/api/v1/product/",
		[]byte(`{"title":"Milk Shake", "description":  "Milk Shake", "link":  "https://test.com/test", "image":  "https://inage.com/1.jpg", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = `{"status":201,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"},"stock":0},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?list=1")
	actual = getRequest("http://localhost:8081/api/v1/product/?list=1")
	expected = `{"status":200,"data":[{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[],"list_id":1,"list":null,"stock":0}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?list=2")
	actual = getRequest("http://localhost:8081/api/v1/product/?list=2")
	expected = `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?category=1")
	actual = getRequest("http://localhost:8081/api/v1/product/?category=1")
	expected = `{"status":200,"data":[{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[],"list_id":1,"list":null,"stock":0}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get by ?category=2")
	actual = getRequest("http://localhost:8081/api/v1/product/?category=2")
	expected = `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")
}

