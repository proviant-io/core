package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	internalHttp "github.com/proviant-io/core/internal/http"
	"gotest.tools/assert"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

const (
	GET    = "get"
	POST   = "post"
	PUT    = "put"
	DELETE = "delete"
)

func toJson(t *testing.T, obj internalHttp.Response) string {
	expectedJson, err := json.Marshal(obj)
	assert.NilError(t, err)
	return string(expectedJson)
}

func url(t string, params ...interface{}) string {
	return fmt.Sprintf(t, params...)
}

func execSuitStep(t *testing.T, title string, uri string, requestType string, requestPayload string, e apiResponse) {
	fmt.Print(title)
	var actual string
	finalUrl := generateApiUrl(uri)
	switch requestType {
	case GET:
		actual = getRequest(finalUrl)
	case POST:
		actual = postRequest(finalUrl, []byte(requestPayload))
	}
	assert.Equal(t, generateExpResponse(e.Status, e.Data, e.Error), actual)
	fmt.Println(" OK")
}

type apiResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

func generateExpResponse(s int, d interface{}, e string) string {
	expected, err := json.Marshal(apiResponse{
		Status: s,
		Data:   d,
		Error:  e,
	})

	if err != nil {
		panic(err)
	}

	return string(expected)
}

func generateApiUrl(uri string) string {
	return fmt.Sprintf("http://localhost:8081%s", uri)
}

func getRequest(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return err.Error()
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	return string(body)
}

func deleteRequest(url string) string {
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

func postRequest(url string, json []byte) string {
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

func putRequest(url string, json []byte) string {
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

func runContainer(t *testing.T) string {
	id, err := createNewContainer("brushknight/proviant-core:e2e")
	if err != nil {
		fmt.Errorf(err.Error())
		t.Fail()
		return ""
	}
	return id
}

func stopContainer(t *testing.T, id string) {
	err := stopAndRemoveContainer(id)

	if err != nil {
		fmt.Errorf(err.Error())
		t.Fail()
	}
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
		}, nil, nil, "proviant-e2e")
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
