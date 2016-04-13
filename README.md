# go-xerror

Package xerror extends the functionality of Go's built-in error interface, in several ways:

- errors carry a stack trace of the location where they were first created
- it is possible to attach debug objects to errors for later reporting
- errors hold a list of messages, and can be wrapped by prepending new messages to the list
- errors can natively be treated as Go built-in errors and serialized to a string message or JSON representation

#### Latest stable version (v1):

[![Build Status](https://api.travis-ci.org/ibrt/go-xerror.svg?branch=v1)](https://travis-ci.org/ibrt/go-xerror?branch=v1)
[![Coverage Status](https://coveralls.io/repos/github/ibrt/go-xerror/badge.svg?branch=v1)](https://coveralls.io/github/ibrt/go-xerror?branch=v1)
[![GoDoc](https://godoc.org/gopkg.in/ibrt/go-xerror.v1/xerror?status.svg)](https://godoc.org/gopkg.in/ibrt/go-xerror.v1/xerror)

Source code: https://github.com/ibrt/go-xerror/tree/v1
Installation: ```go get gopkg.in/ibrt/go-xerror.v1/xerror```

#### Development version:

[![Build Status](https://api.travis-ci.org/ibrt/go-xerror.svg?branch=master)](https://travis-ci.org/ibrt/go-xerror?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/ibrt/go-xerror/badge.svg?branch=master)](https://coveralls.io/github/ibrt/go-xerror?branch=master)
[![GoDoc](https://godoc.org/github.com/ibrt/go-xerror/xerror?status.svg)](https://godoc.org/github.com/ibrt/go-xerror/xerror)

Source code: https://github.com/ibrt/go-xerror
Installation: ```go get github.com/ibrt/go-xerror/xerror```
