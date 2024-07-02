FROM node:20-bookworm-slim@sha256:ae54182fb22d5e567440c9a27d2f4f856f630e5159390af1f69ab0b5e58f6d66 as base

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
