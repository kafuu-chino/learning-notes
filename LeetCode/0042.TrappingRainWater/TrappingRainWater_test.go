package TrappingRainWater

import "testing"

func Test_trap1(t *testing.T) {
	type args struct {
		height []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Example 1",
			args: args{[]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}},
			want: 6,
		},
		{
			name: "Example 2",
			args: args{[]int{4, 2, 0, 3, 2, 5}},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trap1(tt.args.height); got != tt.want {
				t.Errorf("trap1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trap2(t *testing.T) {
	type args struct {
		height []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Example 1",
			args: args{[]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}},
			want: 6,
		},
		{
			name: "Example 2",
			args: args{[]int{4, 2, 0, 3, 2, 5}},
			want: 9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trap2(tt.args.height); got != tt.want {
				t.Errorf("trap2() = %v, want %v", got, tt.want)
			}
		})
	}
}
