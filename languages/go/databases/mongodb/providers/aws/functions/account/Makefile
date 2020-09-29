.PHONY: build clean deploy

build:
	export GO111MODULE=on
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/email_validation email_validation/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/forgot_password forgot_password/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/is_logged_in is_logged_in/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/login login/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/logout logout/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/register register/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/resend_email_validation resend_email_validation/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
