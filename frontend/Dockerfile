FROM node:24-slim

RUN npm install -g pnpm@11.5.2

WORKDIR /app
COPY package.json pnpm-*.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .

CMD ["pnpm", "dev"]
