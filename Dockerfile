FROM node:22-alpine AS base
RUN apk add --no-cache libc6-compat brotli zopfli zstd pnpm
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install
COPY . ./
RUN pnpm build

FROM caddy:2-alpine
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=base /app/dist/ /var/www/html

# run with:
# docker build . -t preveltekit
# docker run -p3000:3000 preveltekit
