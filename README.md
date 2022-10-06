# HTTP Processing service

## Description
The idea of this web service is to create a particular execution plan of tasks (bash commands), taking into concideration the dependencies between them. As an additional step one can choose to get the sequence of tasks (commands) in the form of a bash script

## Configuration
The server uses the following external dependencies, which should be installed:
### Production
* `go` - preferably versions above `1.19.0`
* `github.com/gin-gonic/gin` - used for the implementation of the REST API
* `github.com/pkg/errors` - used for easier creation of errors
* `gorm.io/gorm` - used for mapping models (go structs) to sql tables

### Testing
* `github.com/golang/mock` - used for mocking external dependencies
* `github.com/onsi/ginkgo` - used as the main testing framework
* `github.com/onsi/gomega` - used for assertions

The following environment variables must be set:
### Server configuration
* `HOST` - env variable, containing the host name, on which the server will be running
* `PORT` - env variable, containing the port number, which the server will run on

## Installation
### Manual
```bash
# Clone repo
git clone https://github.com/danielpenchev98/ProcessingService.git

# Build web-server - several dependencies will be downloaded
go build ./...

# Go to cmd dir, containg the server startup file
cd cmd

# Start the server
go run server.go
```
### Via Docker
```bash
# By default the server will start listining on port 8080 in the container and will start in release mode
# To modify this behaviour modify the Dockerfile in the root directory

docker build -t <image_name> .

docker run -p8080:8080 <image_name> 
```

## Running tests
```bash
# Execute it in web-server directory
ginkgo ./...
```

## API endpoints

|api endpoint | payload | usage | result |
|--|--|--|--|
|`GET /v1/healthcheck` |-| check of the availability of the service | 200 status code if the service is reachable |
|`POST /v1/processing-plan`|`JSON object`containing array of tasks with their dependencies|find the execution order of the tasks|order of execution of tasks|
|`POST /v1/bash`|`JSON object` containing array of tasks with their dependencies|find the execution order of the tasks|bash script|

## Request payload example - for both POST API endpoints

```json
{
    "tasks": [
        {
            "name": "task-1",
            "command": "touch /tmp/file1"
        },
        {
            "name": "task-2",
            "command": "cat /tmp/file1",
            "requires": [
                "task-3"
            ]
        },
        {
            "name": "task-3",
            "command": "echo 'Hello World!' > /tmp/file1",
            "requires": [
                "task-1"
            ]
        },
        {
            "name": "task-4",
            "command": "rm /tmp/file1",
            "requires": [
                "task-2",
                "task-3"
            ]
        }
    ]
}
```

## Response Payload example

### `POST /v1/processing-plan`
```json
[
    {
        "name":"task-1",
        "command":"touch /tmp/file1"
    },
    {
        "name":"task-3",
        "command":"echo 'Hello World!' > /tmp/file1"
    },
    {
        "name":"task-2",
        "command":"cat /tmp/file1"
    },
    {
        "name":"task-4",
        "command":"rm /tmp/file1"
    }
]
```

### `POST /v1/bash`
```bash
#!/usr/bin/env bash

touch /tmp/file1
echo "Hello World!" > /tmp/file1
cat /tmp/file1
rm /tmp/file1
```