package e2e

import (
	"fmt"
	"gotest.tools/assert"
	"testing"
	"time"
)

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
