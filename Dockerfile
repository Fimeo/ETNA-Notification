FROM golang:1.18-alpine as builder

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Move files
COPY cmd cmd
COPY internal internal
COPY .env .env
COPY pkg pkg

# Build
RUN go build -o bin /app/cmd/main.go

CMD [ "./bin" ]
