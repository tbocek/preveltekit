FROM golang:1.25-alpine AS builder

RUN apk add --no-cache curl tar
RUN curl -sSL https://github.com/tinygo-org/tinygo/releases/download/v0.40.1/tinygo0.40.1.linux-amd64.tar.gz | tar -xz -C /usr/local
ENV PATH="/usr/local/tinygo/bin:${PATH}"

WORKDIR /app
COPY . .

RUN go build -o pregokit .

# Build counter example
WORKDIR /app/examples/counter
RUN /app/pregokit counter.goweb
RUN cd dist && go mod init counter && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/fetch
RUN /app/pregokit fetch.goweb
RUN cd dist && go mod init fetch && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/two-way
RUN /app/pregokit two-way.goweb
RUN cd dist && go mod init two-way && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/if
RUN /app/pregokit if.goweb
RUN cd dist && go mod init if && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/each
RUN /app/pregokit each.goweb
RUN cd dist && go mod init each && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/at-html
RUN /app/pregokit at-html.goweb
RUN cd dist && go mod init at-html && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/mod
RUN /app/pregokit mod.goweb
RUN cd dist && go mod init mod && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/class
RUN /app/pregokit class.goweb
RUN cd dist && go mod init class && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/nest
RUN /app/pregokit nest.goweb
RUN cd dist && go mod init nest && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

# Build fetch example
WORKDIR /app/examples/slot
RUN /app/pregokit slot.goweb
RUN cd dist && go mod init slot && tinygo build -o app.wasm -target wasm -no-debug -panic trap -scheduler asyncify -gc leaking .
RUN cp /usr/local/tinygo/targets/wasm_exec.js dist/

FROM caddy:alpine
COPY --from=builder /app/examples/counter/dist/ /srv/counter/
COPY --from=builder /app/examples/fetch/dist/ /srv/fetch/
COPY --from=builder /app/examples/two-way/dist/ /srv/two-way/
COPY --from=builder /app/examples/if/dist/ /srv/if/
COPY --from=builder /app/examples/each/dist/ /srv/each/
COPY --from=builder /app/examples/at-html/dist/ /srv/at-html/
COPY --from=builder /app/examples/mod/dist/ /srv/mod/
COPY --from=builder /app/examples/class/dist/ /srv/class/
COPY --from=builder /app/examples/nest/dist/ /srv/nest/
COPY --from=builder /app/examples/slot/dist/ /srv/slot/

RUN cat > /etc/caddy/Caddyfile <<EOF
:80 {
    file_server
    root * /srv
}
EOF

EXPOSE 80
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]
#CMD ["cp", "-r", "/srv",  "/out"]