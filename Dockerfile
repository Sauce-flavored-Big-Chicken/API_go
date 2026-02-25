FROM golang:1.24-alpine AS api-builder

WORKDIR /src

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/main.go

FROM node:22-alpine AS web-builder

WORKDIR /workspace

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY web/package.json ./web/package.json
RUN pnpm install --frozen-lockfile --filter web...

COPY web ./web

ARG VITE_API_BASE_URL=/
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}

RUN pnpm --filter web build

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata nginx tini

COPY --from=api-builder /out/server ./server
COPY --from=api-builder /src/profile ./profile
COPY --from=web-builder /workspace/web/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/http.d/default.conf
COPY docker-start.sh /app/docker-start.sh

RUN chmod +x /app/docker-start.sh \
    && mkdir -p /app/data /app/profile/upload/image /app/profile/upload/file /app/profile/upload/thumb

ENV SERVER_PORT=8080
ENV DB_PATH=/app/data/data.db

EXPOSE 80

VOLUME ["/app/data", "/app/profile/upload"]

CMD ["tini", "--", "/app/docker-start.sh"]
