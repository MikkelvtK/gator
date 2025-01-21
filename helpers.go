package main

import (
	"fmt"
	"strings"
)

func numArgError(n int, argTypes ...string) error {
	return fmt.Errorf("expected %d arguments: %v", n, strings.Join(argTypes, ", "))
}
