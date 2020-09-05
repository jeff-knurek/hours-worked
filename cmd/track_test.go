package cmd

import (
	"reflect"
	"testing"
	"time"
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

func Test_incrementTime(t *testing.T) {
	sameHour, _ := time.Parse("2006-01-02 15:04", "2020-12-31 9:59")
	newHour, _ := time.Parse("2006-01-02 15:04", "2020-12-31 10:00")
	newDay, _ := time.Parse("2006-01-02 15:04", "2020-12-30 19:23")
	newYear, _ := time.Parse("2006-01-02 15:04", "2021-01-01 18:59")
	tests := []struct {
		name string
		u    years
		t    time.Time
		want years
	}{
		{
			name: "increment existing",
			u:    years{"2020": {"December": {"31": 4}}},
			t:    sameHour,
			want: years{"2020": {"December": {"31": 5}}},
		},
		{
			name: "new hour",
			u:    years{"2020": {"December": {"31": 4}}},
			t:    newHour,
			want: years{"2020": {"December": {"31": 5}}},
		},
		{
			name: "new day",
			u:    years{"2020": {"December": {"31": 180}}},
			t:    newDay,
			want: years{"2020": {"December": {"31": 180, "30": 1}}},
		},
		{
			name: "new year",
			u:    years{"2020": {"December": {"31": 487}}},
			t:    newYear,
			want: years{"2020": {"December": {"31": 487}}, "2021": {"January": {"1": 1}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := incrementTime(tt.u, tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("incrementTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
