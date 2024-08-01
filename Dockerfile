FROM alpine:latest
RUN apk add --no-cache curl
WORKDIR /app
COPY ip2region.db /app/ip2region.db
COPY ip2region /app/ip2region
ENTRYPOINT ["ip2region"]
