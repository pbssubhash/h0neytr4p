FROM golang:1.21-alpine AS builder
WORKDIR /app
#
COPY . .
#
RUN apk -U add git openssl
RUN go mod download
RUN go mod tidy
RUN go build -o /app/main
#
RUN openssl req \
          -nodes \
          -x509 \
          -sha512 \
          -newkey rsa:4096 \
          -keyout "app.key" \
          -out "app.crt" \
          -days 3650 \
          -subj '/C=AU/ST=Some-State/O=Internet Widgits Pty Ltd'
#
FROM scratch
#
COPY --from=builder /app/main /opt/h0neytr4p/h0neytr4p
COPY --from=builder /app/traps /opt/h0neytr4p/traps
COPY --from=builder /app/app.key /opt/h0neytr4p/
COPY --from=builder /app/app.crt /opt/h0neytr4p/
#
WORKDIR /opt/h0neytr4p
CMD ["-cert=app.crt", "-key=app.key", "-log=log/log.json", "-payload=/opt/h0neytr4p/payloads/", "-traps=traps/"]
ENTRYPOINT ["./h0neytr4p"]
