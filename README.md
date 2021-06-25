# Proviant


## About The Project

This is yet another pantry organisation system, created with simplicity and automation in mind.
Here's why:
* We are tired of expired food in the pantry
* We are tired of manual ordering food
* We are tired of running out of stock for favorite ingredients and food

### Built With

* [GO](https://golang.org/)
* [GIN](https://github.com/gin-gonic/gin)
* [GORM](https://gorm.io/index.html)

## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* [Docker](https://docs.docker.com/get-docker/)
* [GNU make](https://www.gnu.org/software/make/)

### Dev installation

1. Clone the repo
2. Run app
   ```sh
   make docker/run
   ```
3. Open browser http://0.0.0.0:8080/api/v1/product
4. List of possible htp endpoints listed in this folder [http-calls](./http-calls)

### Installation

Preferable installation is via official docker container

```shell
docker run -d -p8080:80 --name="proviant" brushknight/proviant:latest
```

You can expose db file if you want to persist it via
```shell
-v /home/user/provian/db:/app/db
```