FROM scratch
COPY go-smtp-gateway /
ENTRYPOINT ["/go-smtp-gateway"]