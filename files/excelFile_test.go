package files

import "testing"

func TestGetPositionString(t *testing.T) {
	type args struct {
		row    int
		column int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{1, 27}, "2AB"},
		{"test2", args{2, 2}, "3C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPositionString(tt.args.row, tt.args.column); got != tt.want {
				t.Errorf("GetPositionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
