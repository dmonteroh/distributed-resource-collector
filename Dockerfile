FROM golang:1.17.2 as development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/cespare/reflex@latest

EXPOSE 8080

CMD export GIN_MODE=release
CMD reflex -g '*.go' go run cmd/main.go --start-service