FROM golang:1.21 as buildgo

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN CGO_ENABLED=0 go build ./cmd/cedd
RUN CGO_ENABLED=0 go build -o ced-cli ./cmd/ced

FROM node:16-alpine as buildsvelte

RUN mkdir /app
COPY web/ app/
WORKDIR /app
RUN npm install
RUN npm run build

FROM golang:1.21 as overmind
RUN CGO_ENABLED=0 go install github.com/DarthSim/overmind/v2@v2.4.0

FROM node:16-alpine

RUN apk add tmux

RUN mkdir /app
RUN mkdir /app/build
COPY Procfile /app/
COPY --from=buildgo /app/cedd /app/
COPY --from=buildgo /app/ced-cli /app/
COPY --from=buildsvelte /app/build /app/package.json /app/build/
COPY --from=overmind /go/bin/overmind /app/

ENV DB_PATH=/app/ced.db
ENV OVERMIND_PROCFILE=/app/Procfile NODE_ENV=production

CMD ["/app/overmind", "start"]

