# syntax=docker/dockerfile:1



FROM golang:1.18.2-alpine



WORKDIR /employee-api
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./
RUN go mod download
RUN go build -o ./employee-api

EXPOSE 5050
CMD [ "./employee-api" ]
