# 1. Build frontend widget
FROM node:22-alpine AS widget-builder
WORKDIR /app/widget
RUN corepack enable && corepack prepare pnpm@latest --activate
COPY widget/package.json widget/pnpm-lock.yaml* ./
RUN pnpm install --ignore-scripts
COPY widget/ ./
RUN pnpm build

# 2. Build backend binary
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum* ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# 3. Final lightweight container
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=backend-builder /app/server /app/server
COPY --from=widget-builder /app/widget/dist /app/static
COPY --from=widget-builder /app/widget/index.html /app/static/demo.html

RUN mkdir -p /app/data
VOLUME /app/data

EXPOSE 8080
ENV PORT=8080
ENV DB_DRIVER=sqlite
ENV DB_DSN=/app/data/support.db
ENV STATIC_DIR=/app/static

CMD ["/app/server"]
