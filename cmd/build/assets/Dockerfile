FROM tinygo/tinygo:latest AS build

USER root

RUN apt-get update && apt-get install -y --no-install-recommends \
    wabt \
    brotli \
    zopfli \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

RUN ./build.sh --release .

FROM caddy:latest

COPY --from=build /app/dist/ /srv/
COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 8080
