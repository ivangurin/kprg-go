FROM golang:alpine as build
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN apk add --no-cache gcc g++
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o ./app ./cmd/rester/main.go 

FROM alpine
WORKDIR /app
COPY --from=build /go/src/app/app ./app
COPY --from=build /go/src/app/.env ./.env
CMD ["./app"]