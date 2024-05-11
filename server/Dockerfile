FROM golang:1.22@sha256:b1e05e2c918f52c59d39ce7d5844f73b2f4511f7734add8bb98c9ecdd4443365 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/cedd
RUN CGO_ENABLED=0 go build -o ced-cli ./cmd/ced

FROM alpine@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

RUN mkdir /app
COPY --from=buildgo /app/cedd /app/
COPY --from=buildgo /app/ced-cli /app/
COPY --from=buildgo /usr/local/go/LICENSE /app/GO-LICENSE
COPY licenses.txt /app/THIRD-PARTY-LICENSES

ENV DB_PATH=/app/ced.db
ENV HTTP_PORT=8080

CMD ["/app/cedd"]

