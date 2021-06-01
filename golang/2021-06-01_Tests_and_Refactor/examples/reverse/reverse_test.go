package reverse

import (
	"testing"
)

func TestWordReverse(t *testing.T) {
	cases := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "a name for the success test case",
			str:      "Hello World",
			expected: "dlroW olleH",
		},
		{
			name:     "empty str",
			str:      "",
			expected: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := WordReverse(c.str)

			if result != c.expected {
				t.Errorf("got %s, expected %s", result, c.expected)
			}
		})
	}
}

func BenchmarkWordReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WordReverse("Hello World")
	}
}
