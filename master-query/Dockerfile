FROM golang:1.17

WORKDIR /app
COPY master-query/go.mod .
COPY master-query/go.sum .
RUN go mod download

COPY master-query/. .

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
CMD ["air"]

# CMD ["go", "run", "main.go"]