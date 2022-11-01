FROM alpine:latest

RUN mkdir /app

COPY echoRest /app

CMD [ "/app/echoRest" ]