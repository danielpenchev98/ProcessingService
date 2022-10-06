FROM golang:1.19.2-alpine

# Server config
ENV HOST 0.0.0.0
ENV PORT 8080

# two modes are available - release and debug
ENV GIN_MODE release 

COPY . ./application/
WORKDIR ./application/
RUN go mod download
RUN GOOS=linux go build -o ./bin/service ./cmd

EXPOSE ${PORT}
ENTRYPOINT ["./bin/service"]