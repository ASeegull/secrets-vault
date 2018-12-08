FROM alpine

RUN apk add --update --no-cache ca-certificates

WORKDIR /usr/lib/secrets-vault

COPY ./secrets-server /usr/lib/secrets-vault/secrets-server

RUN chmod +x /usr/lib/secrets-vault/secrets-server

ENTRYPOINT [ "/usr/lib/secrets-vault/secrets-server" ]
EXPOSE 5000