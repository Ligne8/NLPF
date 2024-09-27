FROM golang:1.23.1-alpine as builder

ARG GID=10001
ARG GROUP=ligne8
ARG UID=10001
ARG USER=ligne8

RUN apk add --no-cache gcc musl-dev
RUN addgroup -g ${GID} ${GROUP} && adduser -D -u ${UID} -G ${GROUP} ${USER}
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o tms-backend .

FROM alpine:latest

ENV GIN_MODE=release

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/tms-backend /app/tms-backend

USER ligne8:ligne8
EXPOSE 8080/tcp

ENTRYPOINT ["/app/tms-backend"]
