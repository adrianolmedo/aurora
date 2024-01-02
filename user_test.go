package aurora

import "testing"

func TestCheckEmptyFields(t *testing.T) {
	tt := []struct {
		name        string
		user        User
		errExpected bool
	}{
		{
			name:        "empty-struct",
			user:        User{},
			errExpected: true,
		},
		{
			name:        "empty-fields",
			user:        User{Name: ""},
			errExpected: true,
		},
		{
			name: "filled-fields",
			user: User{
				Name: "Adrián",
			},
			errExpected: false,
		},
	}

	for _, tc := range tt {
		err := tc.user.Validate()
		errReceived := err != nil

		if errReceived != tc.errExpected {
			t.Fatalf("%s: validateUser: unexpected error status: %v", tc.name, err)
		}
	}
}
