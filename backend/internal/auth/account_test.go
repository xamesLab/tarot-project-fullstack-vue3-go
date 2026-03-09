package auth

import "testing"

func TestValidateCredentials(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		login        string
		password     string
		wantLogin    string
		wantPassword string
		wantErr      bool
	}{
		{
			name:         "valid and trimmed",
			login:        "  alice  ",
			password:     "  strongpass  ",
			wantLogin:    "alice",
			wantPassword: "strongpass",
		},
		{
			name:     "login too short",
			login:    "ab",
			password: "strongpass",
			wantErr:  true,
		},
		{
			name:     "password too short",
			login:    "alice",
			password: "short",
			wantErr:  true,
		},
		{
			name:     "empty login",
			login:    "   ",
			password: "strongpass",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotLogin, gotPassword, err := validateCredentials(tt.login, tt.password)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotLogin != tt.wantLogin {
				t.Fatalf("unexpected login: got %q want %q", gotLogin, tt.wantLogin)
			}
			if gotPassword != tt.wantPassword {
				t.Fatalf("unexpected password: got %q want %q", gotPassword, tt.wantPassword)
			}
		})
	}
}
