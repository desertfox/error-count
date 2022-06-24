# syntax=docker/dockerfile:1
FROM golang:1.17.5-alpine as builder

RUN mkdir /opt/error-count

WORKDIR /opt/error-count

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ec .

FROM scratch

WORKDIR /root

COPY --from=builder /opt/error-count ./

CMD [ "./ec" ]