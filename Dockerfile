FROM golang:1.17.2-alpine as development

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .
RUN go build -o /build

EXPOSE $INTERNAL_PORT

CMD [ "/build" ]