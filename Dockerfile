# Stage 1: Build the Go binary
FROM golang:1.21.4-bullseye as stage1

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o authgear-sms-gateway-tencent .

# Stage 2: Prepare the actual fs we use to run the program
FROM debian:bullseye-20220125-slim
WORKDIR /app
# /etc/mime.types (mime-support)
# /usr/share/ca-certificates/*/* (ca-certificates)
# /usr/share/zoneinfo/ (tzdata)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libmagic-dev \
    libmagic-mgc \
    ca-certificates \
    mime-support \
    tzdata \
    && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates
COPY --from=stage1 /src/authgear-sms-gateway-tencent /usr/local/bin/
COPY ./docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["docker-entrypoint.sh"]
# update-ca-certificates requires root to run.
#USER nobody
EXPOSE 7000
CMD ["authgear-sms-gateway-tencent"]
