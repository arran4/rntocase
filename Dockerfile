FROM alpine:3.20
WORKDIR /usr/local/bin
COPY rntocase /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/rntocase"]
