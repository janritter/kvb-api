FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root

FROM scratch

ADD dist/kvb-api /
COPY --from=0 /etc/ssl/certs/ca-certificates.crt ./etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/kvb-api"]
EXPOSE 8080/tcp
