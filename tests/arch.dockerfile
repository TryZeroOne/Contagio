FROM archlinux:latest

WORKDIR /instlr

RUN pacman -Fy --noconfirm
RUN pacman -Syu --noconfirm
RUN pacman -S sudo --noconfirm
COPY ./scripts/installer.sh /instlr/
