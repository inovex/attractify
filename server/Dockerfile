FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY /server/go.mod .
COPY /server/go.sum .
RUN go mod download
ADD ./server/ .
WORKDIR /app/cmd/
RUN cd server && go build -o server .
RUN cd cron && go build -o cron .
RUN cd consumer && go build -o consumer .
RUN cd attractify && go build -o attractify .


FROM alpine
LABEL org.opencontainers.image.source https://github.com/inovex/attractify
RUN apk add tzdata
COPY --from=builder /app/cmd/attractify/attractify /bin/
WORKDIR /app
ADD frontend/dist /app/dist/
COPY --from=builder /app/cmd/server/server .
COPY --from=builder /app/cmd/cron/cron .
COPY --from=builder /app/cmd/consumer/consumer .
CMD ["sh", "-c", "/app/server config.json & /app/cron config.json & /app/consumer config.json"]
