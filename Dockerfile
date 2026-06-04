FROM node:22.20.0-alpine3.22 AS node_modules_deps_builder
WORKDIR /web
COPY web/package.json .
RUN #npm config set registry http://registry.npmmirror.com
RUN npm install

FROM node:22.20.0-alpine3.22 AS web_builder
RUN apk add --no-cache git
WORKDIR /
COPY . .
COPY --from=node_modules_deps_builder /web/node_modules /web/node_modules
RUN cd web && npm install && npm rebuild && npm run build

FROM golang:1.25.1-alpine3.22 AS builder
WORKDIR /app
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -o ./bin/server .

FROM alpine:3.22
COPY --from=web_builder /web/dist/ /assets/web/v1/
COPY --from=builder /app/bin/server /server
RUN chmod +x server
EXPOSE 1423
CMD ["./server"]
