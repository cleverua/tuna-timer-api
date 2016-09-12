FROM scratch

ADD slack-time-linux-amd64 /
ADD data/migrations /data/migrations
ADD config.example.yml /config.yml

WORKDIR /

EXPOSE 8080

CMD ["/slack-time-linux-amd64"]

#GOOS=linux GOARCH=amd64 go build -o slack-time-linux-amd64 .
#docker build -t slack-time .
#docker run -p 8080:8080 -e DATABASE_USER= -e DATABASE_PASS= -e DATABASE_HOST=192.168.0.83 -e DATABASE_NAME=slack_time_dev slack-time