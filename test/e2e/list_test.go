package e2e

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gotest.tools/assert"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
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

	// check that db is empty
	fmt.Print("category: check db empty")
	actual := getRequest("http://localhost:8081/api/v1/category/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// create category
	fmt.Print("category: create")
	actual = postRequest("http://localhost:8081/api/v1/category/", []byte(`{"title": "Drinks"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing category
	fmt.Print("category: get created")
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// update category
	fmt.Print("category: update")
	actual = putRequest("http://localhost:8081/api/v1/category/1/", []byte(`{"title": "Cold Drinks"}`))
	expected = `{"status":200,"data":{"id":1,"title":"Cold Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing category
	fmt.Print("category: get updated")
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Cold Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get all categories
	fmt.Print("category: get all")
	actual = getRequest("http://localhost:8081/api/v1/category/")
	expected = `{"status":200,"data":[{"id":1,"title":"Cold Drinks"}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// delete category 1
	fmt.Print("category: delete")
	actual = deleteRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// check that category does not exists
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

	// check that db is empty
	fmt.Print("product: check db empty")
	actual := getRequest("http://localhost:8081/api/v1/product/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// create list
	fmt.Print("product: create list")
	actual = postRequest("http://localhost:8081/api/v1/list/", []byte(`{"title": "Fridge"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// create category
	fmt.Print("product: create category")
	actual = postRequest("http://localhost:8081/api/v1/category/", []byte(`{"title": "Drinks"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// create product
	fmt.Print("product: create")
	actual = postRequest("http://localhost:8081/api/v1/product/",
		[]byte(`{"title":"Milk Shake", "description":  "Milk Shake", "link":  "https://test.com/test", "image":  "https://inage.com/1.jpg", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = `{"status":201,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing product
	fmt.Print("product: get")
	actual = getRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake","description":"Milk Shake","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// update product
	fmt.Print("product: update")
	actual = putRequest("http://localhost:8081/api/v1/product/1/",
		[]byte(`{"title":"Milk Shake 2", "description":  "Milk Shake 2", "link":  "https://test.com/test", "image":  "https://inage.com/1.jpg", "barcode":  "1234567890Z", "list_id": 1, "category_ids":  [1]}`))
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get existing product
	fmt.Print("product: get updated")
	actual = getRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[{"id":1,"title":"Drinks"}],"list_id":1,"list":{"id":1,"title":"Fridge"}},"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// get all products
	fmt.Print("product: get all")
	actual = getRequest("http://localhost:8081/api/v1/product/")
	expected = `{"status":200,"data":[{"id":1,"title":"Milk Shake 2","description":"Milk Shake 2","link":"https://test.com/test","image":"https://inage.com/1.jpg","barcode":"1234567890Z","category_ids":[1],"categories":[],"list_id":1,"list":null}],"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// delete product 1
	fmt.Print("product: delete")
	actual = deleteRequest("http://localhost:8081/api/v1/product/1/")
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)
	fmt.Println(" OK")

	// check that product does not exists
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
}

func getRequest(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return err.Error()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return string(body)
}

func deleteRequest(url string) string{
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func postRequest(url string, json []byte) string{
	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func putRequest(url string, json []byte) string{
	//var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func runContainer(t *testing.T) string{
	id, err := createNewContainer("pantry:latest")
	if err != nil{
		fmt.Errorf(err.Error())
		t.Fail()
		return ""
	}
	return id
}

func stopContainer(t *testing.T, id string){
	err := stopAndRemoveContainer(id)

	if err != nil{
		fmt.Errorf(err.Error())
		t.Fail()
	}
}

func execCommand(command, params string) {
	cmd := exec.Command(command, params)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}


func createNewContainer(image string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "8081",
	}
	containerPort, err := nat.NewPort("tcp", "80")
	if err != nil {
		panic("Unable to get the port")
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, nil, "test")
	if err != nil {
		panic(err)
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		panic("Unable to start container")
	}
	fmt.Printf("Container %s is started\n", cont.ID)
	return cont.ID, nil
}

func stopAndRemoveContainer(id string) error {
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	ctx := context.Background()

	if err := client.ContainerStop(ctx, id, nil); err != nil {
		log.Printf("Unable to stop container %s: %s", id, err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := client.ContainerRemove(ctx, id, removeOptions); err != nil {
		log.Printf("Unable to remove container: %s\n", err)
		return err
	}

	fmt.Printf("Container %s removed\n", id)

	return nil
}