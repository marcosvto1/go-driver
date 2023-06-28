GODRIVER

# Show coverage test and show output file
> go test -coverprofile=coverage.out ./internal/... -v

# Visualize navigator coverage tests output in html
go tool cover -html=coverage.out
