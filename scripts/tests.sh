# This script runs all Go tests in the project with ensuring cache is cleared.
## Pre-requisite: Go should be installed and Docker should be running.
go clean -testcache
go test -v ./...
