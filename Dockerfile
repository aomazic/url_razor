FROM ubuntu:latest
LABEL authors="kaiser"

ENTRYPOINT ["top", "-b"]