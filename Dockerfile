# Use the official Golang image as a builder.
FROM golang:1.20.6 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# Allows for caching of dependencies unless go.{mod,sum} change.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build each microservice and the API gateway.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o gateway ./app/gateway/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o user ./app/user/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o relation ./app/relation/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o video ./app/video/cmd

# Use the Alpine image for the runtime container.
FROM alpine:3.14

# Copy each compiled binary from the builder stage.
COPY --from=builder /app/gateway /app/gateway
COPY --from=builder /app/user /app/user
COPY --from=builder /app/relation /app/relation
COPY --from=builder /app/video /app/video


# Add a script to start services and the API gateway in the specified order.
COPY run.sh /app/run.sh
COPY config/configuration.yaml /app/configuration.yaml
COPY public/sensitiveDict.txt /app/sensitiveDict.txt
RUN chmod +x /app/run.sh


# 添加执行权限
RUN chmod +x /app/video
RUN chmod +x /app/user
RUN chmod +x /app/relation
RUN chmod +x /app/gateway


CMD ["/app/video", "&", "/app/user", "&", "/app/relation", "&", "/app/gateway"]