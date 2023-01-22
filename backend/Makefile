test_only:
	go test -v -count=1 -cover -vet=off $(path)...

test_with_report:
	go test -v -count=1 -cover -coverprofile=report.out -vet=off $(path)...
	go tool cover -html=report.out -o report.html

test_clean:
	go clean -testcache