FROM node:20-bookworm-slim@sha256:9f938a1eeb3f85ca7691e1b4b5e9ab91e1d2efa7afc1d3e495b6f8158b7e2d39 as base

# build step
FROM base as build

RUN mkdir /app
COPY / app/
WORKDIR /app
RUN npm install

# first, build once so that licenses file gets created.
RUN GENERATE_LICENSES=true npm run build
# build again to include the licenses file in the bundle.
RUN npm run build

# final image
FROM base

RUN mkdir -p /app/build
COPY --from=build /app/build /app/package.json /app/entry.js /app/
COPY --from=build /usr/local/LICENSE /app/NODE-LICENSE

ENV NODE_ENV=production PORT=8080

CMD ["node", "/app/entry.js"]
