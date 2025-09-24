FROM node:24-alpine AS base
RUN apk add --no-cache brotli zopfli zstd pnpm
WORKDIR /app
COPY pnpm-lock.yaml package.json tsconfig.json ./
RUN pnpm install --frozen-lockfile
COPY src ./src
RUN pnpm build

FROM caddy:2-alpine
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=base /app/dist/ /var/www/html

# run with:
# docker build . -t preveltekit
# docker run -p3000:3000 preveltekit
