FROM golang:alpine

WORKDIR /build

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /main

EXPOSE 9000

CMD ["/main"]