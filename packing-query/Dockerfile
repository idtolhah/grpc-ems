FROM golang:1.17

WORKDIR /app
COPY packing-query/go.mod .
COPY packing-query/go.sum .
RUN go mod download

COPY packing-query/. .

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
CMD ["air"]

# CMD ["go", "run", "main.go"]