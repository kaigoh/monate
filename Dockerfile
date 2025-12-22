# --- UI build stage ---
FROM oven/bun:1 AS ui-build
WORKDIR /app
COPY gui ./gui
COPY graph/schema.graphqls ./graph/schema.graphqls
WORKDIR /app/gui
RUN bun install && bun run generate && bun run build

# --- Go build stage ---
FROM golang:1.24 AS go-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy the Bun-built UI bundle from the previous stage
RUN rm -rf web/dist && mkdir -p web/dist
COPY --from=ui-build /app/gui/build/ /app/web/dist/
RUN GOCACHE=/tmp/.gocache go build -o monate ./

# --- Runtime image ---
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=go-build /app/monate /app/monate
COPY --from=go-build /app/web/dist /app/web/dist
EXPOSE 8080
ENV PORT=8080
CMD ["/app/monate"]
