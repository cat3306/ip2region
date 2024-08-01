FROM xyz-fun.tencentcloudcr.com/booster/alpine:1.0
COPY ip2region.db /app/ip2region.db
COPY ip2region /usr/bin/ip2region
WORKDIR /app
ENTRYPOINT ["ip2region"]