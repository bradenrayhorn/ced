FROM node:20-bookworm-slim@sha256:df85129996d6b7a4ec702ebf2142cfa683f28b1d33446faec12168d122d3410d as base

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
