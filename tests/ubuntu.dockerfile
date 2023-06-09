FROM ubuntu:latest

WORKDIR /instlr
RUN apt-get update -y
RUN apt-get install sudo -y

COPY ./scripts/installer.sh /instlr/
