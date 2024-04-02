FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weatherbycep cmd/server/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/weatherbycep .
ENTRYPOINT [ "./weatherbycep" ]