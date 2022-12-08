FROM golang:1.19.1-alpine3.15

ENV PORT 80
EXPOSE 80

RUN mkdir -p /app
WORKDIR /app

COPY code ./

CMD ["go", "run", "index.go"]

