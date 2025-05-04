# Contributing to Verdex
We appreciate your interest in contributing to verdex!  
This document provides some basic guidelines for contributors.

## Getting Started
- Always base your work from the `develop` branch, which is the development branch with the latest code.
- Before creating a Pull Request (PR), make sure there is a corresponding issue for your contribution. If there isn't one already, please create one.
- Include the problem description in the issue.
- Run test cases on all product versions before submitting the Pull Request (PR):
```bash
# all versions:
go run . -test -product keycloak

# single version:
go run . -test -product keycloak -test-version 26.0.5

# versions range:
go run . -test -product keycloak -test-version "26.*"
go run . -test -product keycloak -test-version ">= 26.0.1 < 26.0.5"

# with real target:
go run . -target https://target.fr -product keycloak
```

## Code Style
Please adhere to the existing coding style for consistency.

## Questions
If you have any questions or need further guidance, please feel free to ask in the issue or PR, or reach out to the maintainers.  
Thank you for your contribution!

## Documentation
Use [Mintlify local development guide](https://mintlify.com/docs/development) to improve documentation in `/docs`.

## Releases
- To create a new CLI release:
    - Create a GitHub tag named `v<cli-version>`
    - Run on this tag the `Build and release binary` GitHub action
    - Once completed, run on the same tag the `Build and release Docker` GitHub action
- To create a new Templates release, create a GitHub tag and release named `templates-<templates-version>`

Note that CLI version and Templates version are different.
