package reporting

import (
	"reflect"
	"testing"
	"time"
)

func Test_incrementTime(t *testing.T) {
	sameHour, _ := time.Parse("2006-01-02 15:04", "2020-12-31 9:59")
	newHour, _ := time.Parse("2006-01-02 15:04", "2020-12-31 10:00")
	newDay, _ := time.Parse("2006-01-02 15:04", "2020-12-30 19:23")
	newYear, _ := time.Parse("2006-01-02 15:04", "2021-01-01 18:59")
	tests := []struct {
		name        string
		y           Years
		t           time.Time
		want        Years
		wantMinutes int
	}{
		{
			name:        "increment existing",
			y:           Years{"2020": {"December": {"31": 4}}},
			t:           sameHour,
			want:        Years{"2020": {"December": {"31": 5}}},
			wantMinutes: 5,
		},
		{
			name:        "new hour",
			y:           Years{"2020": {"December": {"31": 4}}},
			t:           newHour,
			want:        Years{"2020": {"December": {"31": 5}}},
			wantMinutes: 5,
		},
		{
			name:        "new day",
			y:           Years{"2020": {"December": {"31": 180}}},
			t:           newDay,
			want:        Years{"2020": {"December": {"31": 180, "30": 1}}},
			wantMinutes: 1,
		},
		{
			name:        "new year",
			y:           Years{"2020": {"December": {"31": 487}}},
			t:           newYear,
			want:        Years{"2020": {"December": {"31": 487}}, "2021": {"January": {"1": 1}}},
			wantMinutes: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotMin := incrementTime(tt.y, tt.t)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("incrementTime() = %v, want %v", got, tt.want)
			}
			if gotMin != tt.wantMinutes {
				t.Errorf("incrementTime() did not provide the current minutes: %v, want %v", gotMin, tt.wantMinutes)
			}
		})
	}
}
