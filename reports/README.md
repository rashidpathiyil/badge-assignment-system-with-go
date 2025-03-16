# Reports Directory

This directory contains generated reports from the badge assignment system.

## Available Reports

- `coverage.out`: Test coverage report generated with `go test -coverprofile=coverage.out`

To view coverage as HTML:

```bash
go tool cover -html=reports/coverage.out
```

**Note:** These files should not be committed to version control.
