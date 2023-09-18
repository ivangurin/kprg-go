FROM golang:alpine as build
WORKDIR /go/src/app 
COPY . .
RUN go mod download
RUN go build -o ./app ./cmd/grapher/main.go

FROM alpine
WORKDIR .
COPY --from=build /go/src/app/app ./app
COPY --from=build /go/src/app/.env ./.env
CMD ["./app"]