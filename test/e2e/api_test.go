package e2e

import (
	"fmt"
	"gotest.tools/assert"
	"testing"
	"time"
)

func TestApiList(t *testing.T){

	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	// check that db is empty
	fmt.Print("list: check db empty")
	actual := getRequest("http://localhost:8081/api/v1/list/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// create list
	fmt.Print("list: create")
	actual = postRequest("http://localhost:8081/api/v1/list/", []byte(`{"title": "Fridge"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing list
	fmt.Print("list: get created")
	actual = getRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// update list
	fmt.Print("list: update")
	actual = putRequest("http://localhost:8081/api/v1/list/1/", []byte(`{"title": "Freezer"}`))
	expected = `{"status":200,"data":{"id":1,"title":"Freezer"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing list
	fmt.Print("list: get updated")
	actual = getRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Freezer"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get all lists
	fmt.Print("list: get all")
	actual = getRequest("http://localhost:8081/api/v1/list/")
	expected = `{"status":200,"data":[{"id":1,"title":"Freezer"}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// delete list 1
	fmt.Print("list: delete")
	actual = deleteRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// check that list does not exists
	fmt.Print("list: deleted not found")
	actual = getRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":404,"data":null,"error":"list with id 1 not found"}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")
}

func TestApiCategory(t *testing.T) {

	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	fmt.Print("category: check db empty")
	actual := getRequest("http://localhost:8081/api/v1/category/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("category: create")
	actual = postRequest("http://localhost:8081/api/v1/category/", []byte(`{"title": "Drinks"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("category: get created")
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("category: update")
	actual = putRequest("http://localhost:8081/api/v1/category/1/", []byte(`{"title": "Cold Drinks"}`))
	expected = `{"status":200,"data":{"id":1,"title":"Cold Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("category: get updated")
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Cold Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("category: get all")
	actual = getRequest("http://localhost:8081/api/v1/category/")
	expected = `{"status":200,"data":[{"id":1,"title":"Cold Drinks"}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("category: delete")
	actual = deleteRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("category: deleted not found")
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":404,"data":null,"error":"category with id 1 not found"}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

}

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
	expected = `{"status":201,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get")
	actual = getRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: update")
	actual = putRequest("http://localhost:8081/api/v1/product/1/",
		[]byte(`{"title":"Milk Shake 2", "description":  "Milk Shake 2", "link":  "https://test.com/test", "image":  "https://inage.com/1.jpg", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get updated")
	actual = getRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("product: get all")
	actual = getRequest("http://localhost:8081/api/v1/product/")
	expected = `{"status":200,"data":[{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[],"list_id":1,"list":null}],"error":""}`
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
	expected = `{"status":201,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
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

	fmt.Print("stock: get")
	actual = getRequest("http://localhost:8081/api/v1/product/1/stock/")
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":5,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest("http://localhost:8081/api/v1/product/1/consume/", []byte(`{"quantity":  3}`))
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after consume")
	actual = getRequest("http://localhost:8081/api/v1/product/1/stock/")
	expected = `{"status":200,"data":[{"id":1,"product_id":1,"quantity":2,"expire":1609458959},{"id":2,"product_id":1,"quantity":3,"expire":1609502159}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: consume 3")
	actual = postRequest("http://localhost:8081/api/v1/product/1/consume/", []byte(`{"quantity":  3}`))
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	fmt.Print("stock: get after second consume")
	actual = getRequest("http://localhost:8081/api/v1/product/1/stock/")
	expected = `{"status":200,"data":[{"id":2,"product_id":1,"quantity":2,"expire":1609502159}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")
}

