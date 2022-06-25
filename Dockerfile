# syntax=docker/dockerfile:1
FROM golang:1.17.5-alpine

RUN mkdir /opt/error-count

WORKDIR /opt/error-count

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ec .

CMD [ "./ec" ]