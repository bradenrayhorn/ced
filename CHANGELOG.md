# Changelog

## v0.8.0

- Simplify UI of RSVP page

Full changelog: [v0.7.0...v0.8.0](https://github.com/bradenrayhorn/ced/compare/v0.7.0...v0.8.0)

## v0.7.0

- Add export command to CLI
- Update dependencies

Full changelog: [v0.6.1...v0.7.0](https://github.com/bradenrayhorn/ced/compare/v0.6.1...v0.7.0)

## v0.6.1

- Fix ced-ui image failing to start on arm64 architectures

Full changelog: [v0.6.0...v0.6.1](https://github.com/bradenrayhorn/ced/compare/v0.6.0...v0.6.1)

## v0.6.0

**Breaking Changes**

- Environment variable configuration has been simplified
  - `TRUSTED_CLIENT_IP_HEADER` has been removed from `ced-server` and `ced-ui`. Replaced with `ADDRESS_HEADER` on `ced-ui` only.
  - `ORIGIN` has been removed from `ced-server`. This is no longer needed.
  - `PUBLIC_BASE_API_URL` and `UNPROXIED_BASE_API_URL` have been removed from `ced-ui`. Replaced with `UNPROXIED_SERVER_URL`. A public API URL is no longer required.
  - See [README](https://github.com/bradenrayhorn/ced/blob/v0.6.0/README.md) for details on configuration
- `ced-helm@0.7.1` has also been released. It is compatible with ced >= 0.6.0. It is incompatible with earlier releases.
  - Previous `ced-helm` versions are incompatible with ced >= 0.6.0


**Changes**

- Add about page
- Simplify proxy configuration
- Change license to AGPLv3
- Update dependencies

Full changelog: [v0.5.7...v0.6.0](https://github.com/bradenrayhorn/ced/compare/v0.5.7...v0.6.0)

## v0.5.7

- Update dependencies

Full changelog: [v0.5.6...v0.5.7](https://github.com/bradenrayhorn/ced/compare/v0.5.6...v0.5.7)

## v0.5.6

- Update dependencies

Full changelog: [v0.5.5...v0.5.6](https://github.com/bradenrayhorn/ced/compare/v0.5.5...v0.5.6)

## v0.5.5

- Update dependencies

Full changelog: [v0.5.4...v0.5.5](https://github.com/bradenrayhorn/ced/compare/v0.5.4...v0.5.5)

## v0.5.4

- Update cardstock theme
- Update no results error message
- Update dependencies

Full changelog: [v0.5.3...v0.5.4](https://github.com/bradenrayhorn/ced/compare/v0.5.3...v0.5.4)

## v0.5.3

- Update fonts of cardstock theme
- Update dependencies

Full changelog: [v0.5.2...v0.5.3](https://github.com/bradenrayhorn/ced/compare/v0.5.2...v0.5.3)

## v0.5.2

- Improve iOS Safari experience

Full changelog: [v0.5.1...v0.5.2](https://github.com/bradenrayhorn/ced/compare/v0.5.1...v0.5.2)

## v0.5.1

- Update font loading
- Update dependencies

Full changelog: [v0.5.0...v0.5.1](https://github.com/bradenrayhorn/ced/compare/v0.5.0...v0.5.1)

## v0.5.0

- Preinstall sqlite on server image
- Update dependencies

Full changelog: [v0.4.0...v0.5.0](https://github.com/bradenrayhorn/ced/compare/v0.4.0...v0.5.0)

## v0.4.0

- Improve error handling
- Add additional logs
- Update fonts of cardstock theme

Full changelog: [v0.3.0...v0.4.0](https://github.com/bradenrayhorn/ced/compare/v0.3.0...v0.4.0)

## v0.3.0

- Improved search capabilities, using levenshtein distance
- CSV import command
- Support showing multiple search results
- Update fonts of cardstock theme
- Update dependencies

Full changelog: [v0.2.1...v0.3.0](https://github.com/bradenrayhorn/ced/compare/v0.2.0...v0.3.0)

## v0.2.1

- Update cardstock theme
- Fix some proxy logic

Full changelog: [v0.2.0...v0.2.1](https://github.com/bradenrayhorn/ced/compare/v0.2.1...v0.2.1)

## v0.2.0

- Add favicon
- Update dependencies
- Improve loading state
- Log message when groups updated
- Update proxy logic

Full changelog: [v0.1.0...v0.2.0](https://github.com/bradenrayhorn/ced/compare/v0.1.0...v0.2.0)

## v0.1.0

- Initial release
