package mocks

//go:generate mockery --name=.* --recursive=false --case=underscore --dir ./.. --output . --with-expecter

//go:generate mockery --name=.* --recursive=false --case=underscore --dir ./../../../pkg/logger --output . --with-expecter
