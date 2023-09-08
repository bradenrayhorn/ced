## ced

A small RSVP service.

### Structure

ced is split into two components: the server and the web interface.

The `/server` directory contains the server backend, written in Go.

The `/ui` directory contains the web frontend, written in Svelte.

### Usage

Two docker images are produced by ced that correspond with the components listed above.

- `ghcr.io/bradenrayhorn/ced-server`
- `ghcr.io/bradenrayhorn/ced-ui`

#### Tags

Docker images are tagged in the following manner.

- `latest`: Latest release version of ced.
- `{major}.{minor}.{patch}`: Targets a specific release version of ced.
- `next`: Latest build of ced from the `main` branch. This is less tested than the release build.
- `next-SHA`: Targets a specific commit from the `main` branch.

