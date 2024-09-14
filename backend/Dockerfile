FROM golang:1.15-alpine as builder

ARG GID=10001
ARG GROUP=ligne8
ARG UID=10001
ARG USER=ligne8

# Creates application user and group.
RUN addgroup -g ${GID} ${GROUP} && adduser -D -u ${UID} -G ${GROUP} ${USER}
WORKDIR /app
# Copy the project files inside the container.
COPY . .
# Download project dependencies.
RUN go mod download
# Build the binary totally static
# Source: http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .


FROM scratch

ENV GIN_MODE=release
# Copy user and group from builder container.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group  /etc/group
# Copy the binary built.
COPY --from=builder /app/tms-backend /app/tms-backend
# Use non-privileged user in container.
USER ligne8:ligne8
EXPOSE 8080/tcp
ENTRYPOINT [ "/app/tms-backend" ]
