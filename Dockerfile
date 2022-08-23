## Build
FROM golang:1.19-bullseye AS build

WORKDIR /app
COPY . ./
RUN go mod download
WORKDIR /app/cmd/backend
RUN go build -o /backend


## Deploy
FROM ubuntu:22.04
ENV TZ=Asia/Tokyo
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /opt/
COPY --from=build /backend /opt/backend
RUN apt-get update && apt install -y ssh tzdata git && apt-get clean && rm -rf /var/lib/apt/lists/*

CMD ["/opt/backend", "start", "cron", "--config", "/opt/config/config.json", "--template", "/opt/config/template.json"]
