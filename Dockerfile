# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /goira && chmod a+x /goira

##
## Deploy
##
FROM gcr.io/distroless/base-debian10
#FROM ubuntu

WORKDIR /

COPY --from=build /goira /goira

EXPOSE 1323

USER nonroot:nonroot

CMD ["/goira"]