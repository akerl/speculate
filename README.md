speculate
=========

[![Build Status](https://img.shields.io/travis/com/akerl/speculate.svg)](https://travis-ci.com/akerl/speculate)
[![GitHub release](https://img.shields.io/github/release/akerl/speculate.svg)](https://github.com/akerl/speculate/releases)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)

Tool for assuming roles in AWS accounts

## Usage

To assume the "admin" role in the local account:

```
speculate env admin
```

Speculate is designed to be a bare-bones implementation of the underlying AWS role management calls, so that other tools can layer on top of it. `speculate env` uses the current environment's creds (set in standard AWS environment variables or ~/.aws/credentials or similar) to assume the named role, and it returns the resulting creds in environment export format.

To use these creds, you probably want to do the following:

```
eval $(speculate env admin)
```

That will source the new creds into your environment.

A `console` subcommand is also provided, for creating a URL for AWS's console:

```
speculate console
```

This uses the current environment's credentials, and can thus be chained together like so:

```
eval $(speculate admin) && open $(speculate console)
```

The env command accepts some options for controlling which account to assume a role on, as well as for MFA and other assumption settings:

```
# speculate env --help
Generate temporary credentials, either by assuming a role or requesting a session token

Usage:
  speculate env [ROLENAME] [flags]

Flags:
  -a, --account string   Account ID to assume role on (defaults to source account)
  -h, --help             help for env
  -l, --lifetime int     Set lifetime of credentials in seconds. For SessionToken, must be between 900 and 129600 (default 3600). For AssumeRole, must be between 900 and 43200 (default 3600)
  -m, --mfa              Use MFA when assuming role
      --mfacode string   Code to use for MFA
      --policy string    Set a IAM policy in JSON for the assumed credentials
  -s, --session string   Set session name for assumed role (defaults to origin user name)
```

## Installation

## License

speculate is released under the MIT License. See the bundled LICENSE file for details.
