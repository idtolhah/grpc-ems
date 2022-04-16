FROM golang:1.17

WORKDIR /app
COPY user-query/go.mod .
COPY user-query/go.sum .
RUN go mod download

COPY user-query/. .

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
CMD ["air"]

# CMD ["go", "run", "main.go"]