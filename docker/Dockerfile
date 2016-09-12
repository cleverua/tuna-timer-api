FROM alpine:latest

# Bash needed by `wait-for-it.sh` script
RUN apk add --no-cache bash

ADD slack-time-linux-amd64 /
ADD data/migrations /data/migrations
ADD config.example.yml /config.yml
ADD docker/wait-for-it.sh /wait-for-it.sh

WORKDIR /

EXPOSE 8080

CMD ["/slack-time-linux-amd64"]

# GOOS=linux GOARCH=amd64 go build -o slack-time-linux-amd64 .
# docker build -f docker/Dockerfile -t pavlikus/slack-time:latest .
# docker push pavlikus/slack-time:latest
