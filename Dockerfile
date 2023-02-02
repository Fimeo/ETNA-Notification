FROM golang:1.18 AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 go build  -o /bin/go-app -v ./cmd/main.go

FROM scratch AS runtime
WORKDIR /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/go-app .
COPY --from=build /app/api/* ./api/

EXPOSE 8080

ENTRYPOINT ["./go-app"]
