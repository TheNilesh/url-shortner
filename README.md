# url-shortner

An app that provides easy to remember links to the pages on the internet

## Installation

    # with docker
    docker run -p 8080:8080thenilesh/url-shortner:latest # or simply make docker-run
    # We need redis to run this hence
    docker run -d -p 6379:6379 redis # make redis-run
    # if network configuration is not messy on your end then redis should be discovered

## Development

    # If go version 1.20+ is installed
    make build

    # test and coverage
    make test cover

    # docker build
    make docker-build

    # Feel free to explore makefile

## Design Overview for Developers

    rest package -> svc package -> store

- The rest package is front controller.
- svc package is use case package. It holds domain logic
- store package is repository package. It holds logic related to persistence