FROM node:22-alpine AS base
RUN apk add --no-cache libc6-compat brotli gzip zstd parallel
WORKDIR /app
RUN npm install -g pnpm
COPY package.json pnpm-lock.yaml rsbuild.config.ts ssr-prod.mjs ssr-server.mjs tsconfig.json ./
RUN pnpm install
COPY src ./src
COPY public ./public
RUN pnpm build
RUN find dist -type f \( -name "*.js" -o -name "*.css" -o -name "*.html" \) -print0 | parallel -0 -j+0 'gzip -9k {}; brotli -k {}; zstd -19k {}'

FROM caddy:2-alpine
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=base /app/dist/ /var/www/html