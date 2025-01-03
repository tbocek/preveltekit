FROM node:22-alpine AS base
RUN apk add --no-cache libc6-compat brotli gzip zstd parallel
WORKDIR /app
RUN npm install -g pnpm
COPY package.json pnpm-lock.yaml ./
RUN pnpm install
COPY src ./src
COPY public ./public
COPY rsbuild.config.ts ssr.mjs tsconfig.json ./
RUN pnpm build

FROM caddy:2-alpine
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=base /app/dist/ /var/www/html

# run with:
# docker build . -t tag
# docker run -p3000:3000 tag