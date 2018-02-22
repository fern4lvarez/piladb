Changelog
=========

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Test against new versions of Go: 1.9.x and 1.10

## [0.1.4] - 2017-06-20

### Added
- Build `pilad` with go1.8.3
- Introduce `make race` to identify data races in tests
- pila: Add `Stack.UUID()` function to get thread-safe Stack ID

### Fixed
- pila: Fix data race conditions on Database and Stack ID

## [0.1.3] - 2017-04-15

### Added
- Build `pilad` with go1.8.1

### Changed
- pkg/stack: Use RWMutex for concurrent reads

### Fixed
- pkg/stack: Fix data race conditions on Size and Peek
- pila: Fix data race conditions on concurrent updates of Stack dates

## [0.1.2] - 2017-03-05

### Added
- pilad: Add Go Version to Status
- pilad: Add `/_ping` endpoint
- godoc: Extend packages documentation
- vendor: Update dependencies
- dev: Add pila.sh utilities to Docker image

### Changed
- pilad: Show links of interest in `/` endpoint

### Fixed
- Fix decoding bug when pushing a malformed payload

## [0.1.1] - 2017-02-20

### Added
- Build `pilad` with go1.8
- pila: Allow use of external Stack base implementations. See https://github.com/fern4lvarez/piladb/pull/47
- Add support to codecov.io

### Removed
- config: Remove unused Default() method

## [0.1.0] - 2016-12-20

### Added
- First release!

[Unreleased]: https://github.com/fern4lvarez/piladb/compare/v0.1.4...HEAD
[0.1.4]: https://github.com/fern4lvarez/piladb/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/fern4lvarez/piladb/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/fern4lvarez/piladb/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/fern4lvarez/piladb/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/fern4lvarez/piladb/compare/dda6b656cbd635dab8e9fc6c254a46f01e4e43ca...v0.1.0
