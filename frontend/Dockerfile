FROM node:25-slim

WORKDIR /app

RUN npm install -g pnpm
COPY package.json pnpm-*.yaml ./
RUN pnpm install --frozen-lockfile

COPY . .

CMD ["pnpm", "run", "dev"]
