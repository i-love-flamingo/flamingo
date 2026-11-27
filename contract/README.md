# Contract tests

Tests of framework behavior that applications rely on, written from the outside. Most of them start an app.
They stop apps with signals, so they do not run on Windows.

    go test ./contract/...

`harness.Boot` starts an app inside the test process with the modules, areas and config files of the test.
Such a test cannot be parallel: it shares the process environment and the global dingo scopes.

- `bootstrap/`: modules, areas, scopes and commands
- `config/`: config loading and injection
- `routing/`: routes, filters, path prefixes and reverse routing
- `responder/`: results such as render, data, redirect, download and errors
- `auth/`: auth brokers, OpenID Connect and security
- `internal/`: the harness and the test modules

If a change makes one of these tests fail, applications can break too. Fix the change, or update the test and
describe the new behavior in the pull request.
