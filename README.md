# url-shortner

An app that provides easy to remember links to the pages on the internet

## Installation

    docker run -d -p 6379:6379 redis
    # if network configuration is not messy on your end then redis should be discovered

    docker run -p 8080:8080 --network=host thenilesh/url-shortner:latest

    # Install REST Client extension in vs code
    # and use request samples from testdata/tests.http

## Development

    # If go version 1.20+ is installed
    make build

    # test and coverage
    make test cover

    # docker build
    make docker-build

    # Check Makefile for other options

## Design Overview for Developers

    rest package -> svc package -> store

- The rest package is front controller.
- svc package is use case package. It holds domain logic
- store package is repository package. It holds logic related to persistence.