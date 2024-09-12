FROM golang:1.23.1-alpine3.20 as build

WORKDIR /app
COPY . /app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build --ldflags '-extldflags=-static' -o github-users main.go

FROM alpine:3.20
COPY --from=build /app/github-users /app/
COPY /template /template
ENTRYPOINT ["/app/github-users"]
