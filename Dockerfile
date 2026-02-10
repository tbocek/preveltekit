FROM tinygo/tinygo:latest

USER root

# Install entr (file watcher) and curl
RUN apt-get update && apt-get install -y --no-install-recommends \
    entr \
    curl \
    gpg \
    && rm -rf /var/lib/apt/lists/*

# Install Caddy from official repo
RUN curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg \
    && curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' > /etc/apt/sources.list.d/caddy-stable.list \
    && apt-get update && apt-get install -y --no-install-recommends caddy \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

EXPOSE 8080

CMD ["./dev.sh", "examples"]
