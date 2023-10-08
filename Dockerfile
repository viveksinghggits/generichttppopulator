FROM alpine

COPY generichttppopulator /usr/local/bin
ENTRYPOINT [ "/usr/local/bin/generichttppopulator" ]
