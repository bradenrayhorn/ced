FROM golang:1.21 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/cedd
RUN CGO_ENABLED=0 go build -o ced-cli ./cmd/ced

FROM alpine

RUN mkdir /app
COPY --from=buildgo /app/cedd /app/
COPY --from=buildgo /app/ced-cli /app/

ENV DB_PATH=/app/ced.db
ENV HTTP_PORT=8080

CMD ["/app/cedd"]

