# Overlord

[![Build Status](https://travis-ci.org/meritlabs/overlord.svg?branch=master)](https://travis-ci.org/meritlabs/overlord)

The monitoring application, composed from server `Overlord` and agent `Overseer`.
It performs 3 types of checks:
- node availability through ping requests to the agent;
- node blockchain status;
- node version and protocol version status.

The application reports any inconsittencies to Slack.

## Prerequisites

* Go
* Go modules - [description](https://blog.golang.org/using-go-modules)
* Make 

# Getting started

With `Go`, `Make` and `Go modules` configured, clone the repository and change direcory to the neewe folder.
Then execute `make bootstrap` to install all dependencies required

To compile project, use `make` command.

To run the application:
- build project
- copy `overseer` agent binaries to the machines running nodes and start them
- copy `overlord.yaml.example` to `overlord.yaml`
- edit `overlord.yaml` adding IP addresses of nodes and Slack credentials
- run the `overlord` server.

## Contributing

Please, check out our [Contribution guide](./CONTRIBUTING.md) and [Code of Conduct](./CODE_OF_CONDUCT.md).

## License

**Code released under [the MIT license](./LICENSE).**

Copyright (c) 2017-2020 The Merit Foundation.