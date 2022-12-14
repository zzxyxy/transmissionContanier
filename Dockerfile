FROM debian:bookworm

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install -y transmission-daemon

COPY entry.sh /

ENTRYPOINT /entry.sh
