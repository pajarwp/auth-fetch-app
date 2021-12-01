#build
FROM golang:alpine AS build
WORKDIR /fetch-app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /fetch

##deploy
FROM alpine:3.13

WORKDIR /

COPY --from=build /fetch /fetch

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/fetch"]