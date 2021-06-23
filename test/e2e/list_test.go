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
	actual := getRequest("http://localhost:8081/api/v1/list/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)

	// create list
	actual = postRequest("http://localhost:8081/api/v1/list/", []byte(`{"title": "Fridge"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)

	// get existing list
	actual = getRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Fridge"},"error":""}`
	assert.Equal(t, expected, actual)

	// update list
	actual = putRequest("http://localhost:8081/api/v1/list/1/", []byte(`{"title": "Freezer"}`))
	expected = `{"status":200,"data":{"id":1,"title":"Freezer"},"error":""}`
	assert.Equal(t, expected, actual)

	// get existing list
	actual = getRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Freezer"},"error":""}`
	assert.Equal(t, expected, actual)

	// get all lists
	actual = getRequest("http://localhost:8081/api/v1/list/")
	expected = `{"status":200,"data":[{"id":1,"title":"Freezer"}],"error":""}`
	assert.Equal(t, expected, actual)

	// delete list 1
	actual = deleteRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)

	// check that list does not exists
	actual = getRequest("http://localhost:8081/api/v1/list/1/")
	expected = `{"status":404,"data":null,"error":"list with id 1 not found"}`
	assert.Equal(t, expected, actual)
}

func TestApiCategory(t *testing.T) {

	id := runContainer(t)

	defer stopContainer(t, id)

	time.Sleep(1 * time.Second)

	// check that db is empty
	actual := getRequest("http://localhost:8081/api/v1/category/")
	expected := `{"status":200,"data":[],"error":""}`
	assert.Equal(t, expected, actual)

	// create category
	actual = postRequest("http://localhost:8081/api/v1/category/", []byte(`{"title": "Drinks"}`))
	expected = `{"status":201,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)

	// get existing category
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Drinks"},"error":""}`
	assert.Equal(t, expected, actual)

	// update category
	actual = putRequest("http://localhost:8081/api/v1/category/1/", []byte(`{"title": "Cold Drinks"}`))
	expected = `{"status":200,"data":{"id":1,"title":"Cold Drinks"},"error":""}`
	assert.Equal(t, expected, actual)

	// get existing category
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":{"id":1,"title":"Cold Drinks"},"error":""}`
	assert.Equal(t, expected, actual)

	// get all categories
	actual = getRequest("http://localhost:8081/api/v1/category/")
	expected = `{"status":200,"data":[{"id":1,"title":"Cold Drinks"}],"error":""}`
	assert.Equal(t, expected, actual)

	// delete category 1
	actual = deleteRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":200,"data":null,"error":""}`
	assert.Equal(t, expected, actual)

	// check that category does not exists
	actual = getRequest("http://localhost:8081/api/v1/category/1/")
	expected = `{"status":404,"data":null,"error":"category with id 1 not found"}`
	assert.Equal(t, expected, actual)

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
	fmt.Printf("Container %s is started", cont.ID)
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
		log.Printf("Unable to remove container: %s", err)
		return err
	}

	return nil
}