package tracking

import (
	"reflect"
	"testing"
)

func Test_parseSessionID(t *testing.T) {
	tests := []struct {
		name          string
		sessionOutput string
		want          []string
		wantErr       bool
	}{
		{
			name: "match regex",
			sessionOutput: `
Sessions=c6
`,
			want:    []string{"c6"},
			wantErr: false,
		},
		{
			name: "multiple sessions",
			sessionOutput: `
			Sessions=c6 c10
			`,
			want:    []string{"c6", "c10"},
			wantErr: false,
		},
		{
			name: "invalid input",
			sessionOutput: `
		   SESSION        UID USER             SEAT             TTY
		       c18       1001 jane		       seat0
		       c10       1002 szymon	       seat0
		        c6       1000 sally            seat0

		3 sessions listed.
		`,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSessionID(tt.sessionOutput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSessionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSessionID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseActive(t *testing.T) {
	tests := []struct {
		name    string
		out     string
		want    bool
		wantErr bool
	}{
		{
			name: "match active",
			out: `
State=active
`,
			want:    true,
			wantErr: false,
		},
		{
			name: "match in-active",
			out: `
State=online
`,
			want:    false,
			wantErr: false,
		},
		{
			name: "catch error",
			out: `
Failed to get session path
`,
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseActive(tt.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseActive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
