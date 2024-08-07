FROM golang:1.22@sha256:2bd56f00ff47baf33e64eae7996b65846c7cb5e0a46e0a882ef179fd89654afa as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/cedd
RUN CGO_ENABLED=0 go build -o ced-cli ./cmd/ced

FROM alpine@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

RUN mkdir /app
COPY --from=buildgo /app/cedd /app/
COPY --from=buildgo /app/ced-cli /app/
COPY --from=buildgo /usr/local/go/LICENSE /app/GO-LICENSE
COPY licenses.txt /app/THIRD-PARTY-LICENSES

ENV DB_PATH=/app/ced.db
ENV HTTP_PORT=8080

CMD ["/app/cedd"]

