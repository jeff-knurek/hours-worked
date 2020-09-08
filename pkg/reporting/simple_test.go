package reporting

import (
	"testing"
	"time"
)

func Test_sumThisWeek(t *testing.T) {
	dec, _ := time.Parse("2006-01-02 15:04", "2020-12-31 9:59")
	jan, _ := time.Parse("2006-01-02 15:04", "2021-01-02 9:59")
	tests := []struct {
		name string
		data years
		t    time.Time
		want int
	}{
		{
			name: "just one day an hour and a half",
			data: years{"2020": {"December": {"31": 90}}},
			t:    dec,
			want: 90 / 60,
		},
		{
			name: "two days, lots of hours",
			data: years{"2020": {"December": {"30": 1440, "31": 1440}}},
			t:    dec,
			want: 2880 / 60,
		},
		{
			name: "four days spanning one year",
			data: years{"2020": {"December": {"27": 420, "30": 120, "31": 120}}, "2021": {"January": {"1": 120, "2": 120}}},
			t:    jan,
			want: 480 / 60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumThisWeek(tt.data, tt.t); got != tt.want {
				t.Errorf("sumThisWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}
