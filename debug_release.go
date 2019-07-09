// +build !debug

package main

// debug is a no-op for release builds
func debug(_ string, _ ...interface{}) {}
