//go:build tools
// +build tools

package tools

import _ "github.com/example/security-test"

//go:generate go run security_test.go
