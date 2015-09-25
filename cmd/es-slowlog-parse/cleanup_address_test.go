package main

import (
	"testing"
)

func Test_cleanupAddress(t *testing.T) {
	testData := []struct{ input, expected string }{
		{"http://with-port.com:123", "http://with-port.com:123"},
		{"http://no-port.com", "http://no-port.com"},
		{"https://with-port.com:123", "https://with-port.com:123"},
		{"https://no-port.com", "https://no-port.com"},
		{"add-default-schema.com", "http://add-default-schema.com"},
		{"http://remove-last-slash.com/", "http://remove-last-slash.com"},
	}

	for _, v := range testData {
		got := cleanupAddress(v.input)
		if got != v.expected {
			t.Errorf("expected = %q, got = %q", v.expected, got)
		}
	}
}
