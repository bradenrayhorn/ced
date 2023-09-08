# ced

A small RSVP service.

## Structure

ced is split into two components: the server and the web interface.

The `/server` directory contains the server backend, written in Go.

The `/ui` directory contains the web frontend, written in Svelte.

## Installation

### Docker

Two docker images are produced by ced that correspond with the components listed above.

- `ghcr.io/bradenrayhorn/ced-server`
- `ghcr.io/bradenrayhorn/ced-ui`

Theoretically, ced-server supports running multiple instances of the app connected to the same volume.
Each instance would access the same sqlite database.
**This is not well tested.**

#### Tags

Docker images are tagged in the following manner.

- `latest`: Latest release version of ced.
- `{major}.{minor}.{patch}`: Targets a specific release version of ced.
- `next`: Latest build of ced from the `main` branch. This is less tested than the release build.
- `next-{sha}`: Targets a specific commit from the `main` branch.

### Kubernetes

A helm chart is available for installation. The chart is available at:

- repository: `oci://ghcr.io/bradenrayhorn`
- chart: `ced-helm`

The chart is versioned independently of ced itself.

The chart is currently somewhat inflexible; there are not many configurations supported.

