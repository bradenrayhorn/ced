FROM golang:1.23@sha256:8c10f21bec412f08f73aa7b97ca5ac5f28a39d8a88030ad8a339fd0a781d72b4 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/cedd
RUN CGO_ENABLED=0 go build -o ced-cli ./cmd/ced

FROM alpine@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099

RUN mkdir /app
COPY --from=buildgo /app/cedd /app/
COPY --from=buildgo /app/ced-cli /app/
COPY --from=buildgo /usr/local/go/LICENSE /app/GO-LICENSE
COPY licenses.txt /app/THIRD-PARTY-LICENSES

ENV DB_PATH=/app/ced.db
ENV HTTP_PORT=8080

CMD ["/app/cedd"]

