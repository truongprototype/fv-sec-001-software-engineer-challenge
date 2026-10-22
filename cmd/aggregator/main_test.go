package main

import (
	"strings"
	"testing"
)

func TestRunValidatesArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "missing input", args: []string{"--output", "results"}, want: "--input is required"},
		{name: "missing output", args: []string{"--input", "input.csv"}, want: "--output is required"},
		{name: "positional argument", args: []string{"--input", "input.csv", "--output", "results", "extra"}, want: "unexpected positional arguments"},
		{name: "unknown flag", args: []string{"--unknown"}, want: "flag provided but not defined"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := run(test.args)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("run() error = %v, want containing %q", err, test.want)
			}
		})
	}
}
