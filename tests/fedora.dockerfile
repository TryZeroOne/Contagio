FROM fedora:latest

WORKDIR /instlr

COPY ./scripts/installer.sh /instlr/

# COPY contagio /contagio/contagio
# COPY sqlite /contagio/sqlite
# COPY config.toml /contagio/
# COPY main.go /contagio/
# COPY themes /contagio/themes

# CMD ["sudo", "bash", "installer.sh"]
