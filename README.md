# ced

A small self-hosted RSVP service.

ced features a lightweight interface to allow guests to easily RSVP to your event.

Features:
- Customizable event title and URL
- CSV import of guest list
- Fuzzy searching
- Multiple themes
- Functions in browsers without Javascript enabled
- Easy-to-manage SQLite database

## Questions?

Please see [GitHub discussions](https://github.com/bradenrayhorn/ced/discussions).

## Quickstart

See the [Quickstart guide](QUICKSTART.md) for an example of setting up ced
with Docker Compose.

## Installation

ced is released as two Docker images.

- `ghcr.io/bradenrayhorn/ced-server`
- `ghcr.io/bradenrayhorn/ced-ui`

**Note:** ced-server currently supports running a single instance per SQLite database.
Running multiple instances of ced-server has not been validated.

**Important:** Only ced-ui should be exposed to public traffic. ced-server is designed
to trust certain headers, particularly, `ced-connecting-ip`, that are set by ced-ui.
Allowing ced-server to be accessed directly will leave ced-server vulnerable to simple
IP spoofing.

### Docker Compose

A Docker Compose file can be found [here](docker/docker-compose.yml).

### Kubernetes

A helm chart is available for use at:

- repository: `oci://ghcr.io/bradenrayhorn`
- chart: `ced-helm`

See [here](https://github.com/bradenrayhorn/ced/pkgs/container/ced-helm) for published versions.

The chart optionally supports setting configuring [litestream](https://litestream.io/).
This enables backup and automated recovery of ced sqlite data.

Check [values.yaml](https://github.com/bradenrayhorn/ced/blob/main/helm/ced/values.yaml)
for example configuration and more details.

## Configuration

Configuration is done using environment variables.

**ced-server**

| Variable | Description | Required | Default | Values |
| - | - | - | - | - |
| `HTTP_PORT` | Port to run on. | Yes | | Example: `8080` |
| `DB_PATH` | File path for SQLite file. | Yes | | Example: `ced.db` |
| `PRETTY_LOG` | If false, logs in JSON format. | No | `false` | `true` `false` |

**ced-ui**

| Variable | Description | Required | Default | Values |
| - | - | - | - | - |
| `PUBLIC_EVENT_TITLE` | Title of the event. | Yes | | Example: `My Big Event` |
| `PUBLIC_EVENT_URL` | URL to page with event details. | Yes | | Example: `https://myevent.com/details` |
| `UNPROXIED_SERVER_URL` | Base URL of ced-server that does not go through a proxy. See [Proxies](#proxies). | Yes | | Example: `http://ced-server.cluster.local` |
| `ORIGIN` | Allowed origin for CORS. | No | | Example: `https://ced.example.com` |
| `PUBLIC_EVENT_THEME` | Theme. | No | `hamlindigo` | `hamlindigo` `cardstock` |
| `ADDRESS_HEADER` | See [IP Logging](#ip-logging). | No | | Example: `CF-Connecting-IP` |

### IP Logging

For certain log messages the IP address of the requestor is included.
Properly configuring ced ensures the IP address is accurate.

With a default configuration, the IP address is the network address that sent the request.

If ced is deployed behind a proxy, this results in misleading IP address information,
as the IP address will always be the IP address of the proxy.

To fix this, set the `ADDRESS_HEADER` environment variable to whichever header contains the
IP address of the original client. It is important to make sure this option is set
**ONLY** when ced is behind a trusted proxy **AND** the proxy is setting the header you specify.

For example, if ced is deployed behind Cloudflare, the correct header should be
`CF-Connecting-IP`.

## Usage

After installing ced, the next step is to create groups.
It is recommended to use the CSV import process.

Groups have three attributes that must be set.

- `name`: Display name of the group.
- `maxAttendees`: Maximum amount of guests in the group.
- `searchHints`: Optional. A comma separated list of names used to aid in searching.

### Searching

The search functionality works by matching on the name and each of the search hints in the group.
There is a small amount of fuziness applied to correct for typos.
The search is case-insensitve and ignores all whitespace.

For example, assume the following group is setup:

```csv
Name=Bob Lob and family
Max Attendees=3
Search Hints="Bob Lob, Trevor Lob, Josie Lob, Rob Lob, Robert Lob"
```

Searching for "Bob Lob", "Trevor Lob", "Josie Lob", "Rob Lob", or "Robert Lob" will all
match the group "Bob Lob and family".

Setting the search hints properly will ensure your guests can easily find their RSVP.

**Search hints recommendations:**

- Add all individuals in the group
- Add all variations of names
  - For example, someone named "Pam" might go by both "Pam" and "Pamela"

### CSV Import

Groups can be imported from a CSV using the below command.
The CSV data is read from stdin.

Example usage:

```
cat mydata.csv | /app/ced-cli group import
```

The CSV must be in the following format, with **NO** header line.

```csv
{name}, {max_attendees}, {search_hints}
```

Example CSV:

```csv
"Bob Lob",1,""
"Jerome and Elaine Johnson",2,"Jerome Johnson, Elaine Johnson"
```

### Create Group Command

Groups can also be created individually using the following command.
See usage by passing `--help` flag.

```
/app/ced-cli group create
```

## Contributing

See [contributing guide](./CONTRIBUTING.md).

## Development

### Structure

ced is split into two components: the server and the web interface.

The `/server` directory contains the server backend, written in Go.

The `/ui` directory contains the web frontend, written in Svelte.

Additionally, the `/helm` directory contains the source code for the helm chart.

### Local Setup

The Makefile contains two targets that can be used to start the server and ui locally.

- `make run/server`
- `make run/ui`

