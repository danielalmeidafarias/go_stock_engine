package postgres

import "testing"

func TestSeedFilePath(t *testing.T) {
	tests := []struct {
		name     string
		useSeed  string
		seedFile string
		want     string
	}{
		{name: "disabled", want: ""},
		{name: "false", useSeed: "false", want: ""},
		{name: "default", useSeed: "true", want: defaultSeedFile},
		{name: "configured", useSeed: "true", seedFile: "/tmp/seed.sql", want: "/tmp/seed.sql"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("USE_SEED", tt.useSeed)
			t.Setenv("SEED_FILE", tt.seedFile)

			if got := seedFilePath(); got != tt.want {
				t.Fatalf("seedFilePath() = %q, want %q", got, tt.want)
			}
		})
	}
}
