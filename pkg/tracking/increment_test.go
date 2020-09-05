package tracking

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
