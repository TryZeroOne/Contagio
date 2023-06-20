FROM ubuntu:latest

WORKDIR /instlr
RUN apt-get update -y
RUN apt-get install sudo -y

COPY ./installer/installer.sh /instlr/
