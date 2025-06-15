FROM alpine:3.20
WORKDIR /usr/local/bin
COPY rntodarwin /usr/local/bin/
COPY rntocamel /usr/local/bin/
COPY rnreverse /usr/local/bin/
COPY rntodelimited /usr/local/bin/
COPY rntokebab /usr/local/bin/
COPY rntolowerleading /usr/local/bin/
COPY rntosnake /usr/local/bin/
COPY rntrim /usr/local/bin/
COPY rntoupperleading /usr/local/bin/
COPY rntotitle /usr/local/bin/
COPY rntopascal /usr/local/bin/
COPY rntoupper /usr/local/bin/
COPY rntolower /usr/local/bin/
COPY rnacronym /usr/local/bin/
ENTRYPOINT ["/bin/sh"]
