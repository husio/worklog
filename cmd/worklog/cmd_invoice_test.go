package main

import "testing"

func TestPrettyFormatNumberDE(t *testing.T) {
	cases := map[string]struct {
		input int
		want  string
	}{
		"zero": {
			input: 0,
			want:  "0",
		},
		"no separation": {
			input: 12,
			want:  "12",
		},
		"another no separation": {
			input: 123,
			want:  "123",
		},
		"a single separation": {
			input: 1234,
			want:  "1.234",
		},
		"two separations": {
			input: 1234567,
			want:  "1.234.567",
		},
		"almost three separations": {
			input: 123456789,
			want:  "123.456.789",
		},
		"three separations": {
			input: 1234567890,
			want:  "1.234.567.890",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := prettyFormatNumberDE(tc.input)
			if got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}
