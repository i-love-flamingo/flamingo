package domain

import "testing"

func TestValidateDate(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"valid date 1",
			args{"1981-12-19"},
			true,
		},
		{
			"valid date 2",
			args{"1981-02-30"},
			true,
		},
		{
			"valid date 3",
			args{"1981-11-30"},
			true,
		},
		{
			"valid date 4",
			args{"1981-01-30"},
			true,
		},
		{
			"invalid date format",
			args{"18-11-1981"},
			false,
		},
		{
			"invalid month but correct format",
			args{"1981-19-11"},
			false,
		},
		{
			"invalid day but correct format",
			args{"1981-01-41"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDate(tt.args.date); got != tt.want {
				t.Errorf("ValidateDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAge(t *testing.T) {
	type args struct {
		date string
		age  int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"too young",
			args{"2000-19-11", 1000},
			false,
		},
		{
			"old enough",
			args{"1981-01-30", 4},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateAge(tt.args.date, tt.args.age); got != tt.want {
				t.Errorf("ValidateAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
