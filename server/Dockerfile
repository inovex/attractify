FROM golang:1.17-alpine AS builder
RUN mkdir /app
ADD ./server/ /app/
WORKDIR /app
RUN go mod download
WORKDIR /app/cmd/server
RUN go build -o server .
WORKDIR /app/cmd/cron
RUN go build -o cron .
WORKDIR /app/cmd/consumer
RUN go build -o consumer .

FROM alpine
LABEL org.opencontainers.image.source https://github.com/inovex/attractify
RUN apk add tzdata
RUN mkdir /app
WORKDIR /app
ADD frontend/dist /app/dist/
COPY --from=builder /app/cmd/server .
COPY --from=builder /app/cmd/cron .
COPY --from=builder /app/cmd/consumer .
COPY /server/config.sample.json /app/config.json
CMD ["sh", "-c", "/app/server config.json & /app/cron config.json & /app/consumer config.json"]