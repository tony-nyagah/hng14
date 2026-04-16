FROM oven/bun:latest

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build-only

EXPOSE 4173

CMD ["bun", "run", "preview", "--", "--host", "--port", "4173"]
