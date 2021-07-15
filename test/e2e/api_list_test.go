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
	actual := getRequest(generateApiUrl("/api/v1/list/"))
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// create list
	fmt.Print("list: create")
	actual = postRequest(generateApiUrl("/api/v1/list/"), []byte(`{"title": "Fridge"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing list
	fmt.Print("list: get created")
	actual = getRequest(generateApiUrl("/api/v1/list/1/"))
	expected = `{"status":200,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// update list
	fmt.Print("list: update")
	actual = putRequest(generateApiUrl("/api/v1/list/1/"), []byte(`{"title": "Freezer"}`))
	expected = `{"status":200,"data":{"id":1,"title":"Freezer"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing list
	fmt.Print("list: get updated")
	actual = getRequest(generateApiUrl("/api/v1/list/1/"))
	expected = `{"status":200,"data":{"id":1,"title":"Freezer"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get all lists
	fmt.Print("list: get all")
	actual = getRequest(generateApiUrl("/api/v1/list/"))
	expected = `{"status":200,"data":[{"id":1,"title":"Freezer"}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// delete list 1
	fmt.Print("list: delete")
	actual = deleteRequest(generateApiUrl("/api/v1/list/1/"))
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// check that list does not exists
	fmt.Print("list: deleted not found")
	actual = getRequest(generateApiUrl("/api/v1/list/1/"))
	expected = `{"status":404,"data":null,"error":"list with id 1 not found"}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")
}

