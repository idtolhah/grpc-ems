FROM golang:1.17

WORKDIR /app
COPY filtration-query/go.mod .
COPY filtration-query/go.sum .
RUN go mod download

COPY filtration-query/. .

# RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
# CMD ["air"]

CMD ["go", "run", "main.go"]