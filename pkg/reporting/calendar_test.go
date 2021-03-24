package reporting

import (
	"strings"
	"testing"
	"time"
)



func TestTextCalendar(t *testing.T) {
	dec, _ := time.Parse("2006-01-02 15:04", "2020-12-31 9:59")
	dec, _ = time.Parse("2006-01-02 15:04", "2020-12-31 9:59")
	jan, _ := time.Parse("2006-01-02 15:04", "2021-01-02 9:59")
	tests := []struct {
		name string
		data Years
		t    time.Time
		want string
	}{
		{
			name: "one month",
			data: Years{"2020": {"December": {"27": 420, "30": 120, "31": 120}}},
			t:    dec,
			want: strings.Replace(`December

		| Sun | Mon | Tue | Wed | Thu | Fri | Sat |
		|----:|----:|----:|----:|----:|----:|----:|
		|     |     |   0 |   0 |   0 |   0 |   0 |
		|   0 |   0 |   0 |   0 |   0 |   0 |   0 |
		|   0 |   0 |   0 |   0 |   0 |   0 |   0 |
		|   0 |   0 |   0 |   0 |   0 |   0 |   0 |
		| 7.0 |   0 |   0 | 2.0 | 2.0 |     |     |`, "\t", "", -1),
		},
		{
			name: "no data this month",
			data: Years{"2020": {"October": {"27": 420, "30": 120, "31": 120}}},
			t:    dec,
			want: "--no data tracked for this month--",
		},
		{
			name: "two months",
			data: Years{"2020": {"December": {"27": 420, "30": 120}}, "2021": {"January": {"1": 120, "2": 120}}},
			t:    jan,
			want: strings.Replace(`January

		| Sun | Mon | Tue | Wed | Thu | Fri | Sat |
		|----:|----:|----:|----:|----:|----:|----:|
		|     |     |     |     |     | 2.0 | 2.0 |
		|   . |   . |   . |   . |   . |   . |   . |
		|   . |   . |   . |   . |   . |   . |   . |
		|   . |   . |   . |   . |   . |   . |   . |
		|   . |   . |   . |   . |   . |   . |   . |
		|   . |     |     |     |     |     |     |`, "\t", "", -1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := textCalendar(tt.data, tt.t)
			if got != tt.want {
				t.Errorf("textCalendar() \nGot:\n%v,\nwant:\n%v\n", got, tt.want)
			}
		})
	}
}
