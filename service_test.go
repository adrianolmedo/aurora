package main

import (
	"strings"
	"testing"
)

func TestSignUp(t *testing.T) {
	tt := []struct {
		name           string
		input          *User
		errExpected    bool
		wantErrContain string
	}{
		{
			name: "successful",
			input: &User{
				Name: "Adrián",
			},
			errExpected:    false,
			wantErrContain: "",
		},
		{
			name: "empty-field",
			input: &User{
				Name: "",
			},
			errExpected:    true,
			wantErrContain: "the name can't be empty",
		},
	}

	for _, tc := range tt {
		err := signUp(tc.input)
		if (err != nil) != tc.errExpected {
			t.Fatalf("%s: unexpected error value %v", tc.name, err)
		}

		if err != nil && !strings.Contains(err.Error(), tc.wantErrContain) {
			t.Fatalf("want error string %q to contain %q", err.Error(), tc.wantErrContain)
		}
	}
}
