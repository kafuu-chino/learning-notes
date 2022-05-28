package trap

import "testing"

func TestTrap(t *testing.T) {
	testData := []struct {
		input  []int
		output int
	}{
		{
			input:  []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1},
			output: 6,
		},
		{
			input:  []int{4, 2, 0, 3, 2, 5},
			output: 9,
		},
		{
			input:  []int{0, 1, 0, 1, 0, 0},
			output: 1,
		},
	}

	testFn := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "trap1",
			fn: func(t *testing.T) {
				for _, v := range testData {
					if trap1(v.input) != v.output {
						t.Fatal("An error occurred")
					}
				}
			},
		},
		{
			name: "trap2",
			fn: func(t *testing.T) {
				for _, v := range testData {
					if trap2(v.input) != v.output {
						t.Fatal("An error occurred")
					}
				}
			},
		},
	}

	for _, v := range testFn {
		t.Run(v.name, v.fn)
	}
}
