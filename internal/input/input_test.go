package input

import (
	"strings"
	"testing"
)

func TestFromReader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			"simple ips",
			"8.8.8.8\n1.1.1.1\n",
			[]string{"8.8.8.8", "1.1.1.1"},
		},
		{
			"skip empty lines",
			"8.8.8.8\n\n\n1.1.1.1\n",
			[]string{"8.8.8.8", "1.1.1.1"},
		},
		{
			"skip comments",
			"# this is a comment\n8.8.8.8\n# another comment\n1.1.1.1\n",
			[]string{"8.8.8.8", "1.1.1.1"},
		},
		{
			"trim whitespace",
			"  8.8.8.8  \n\t1.1.1.1\t\n",
			[]string{"8.8.8.8", "1.1.1.1"},
		},
		{
			"trim carriage return",
			"8.8.8.8\r\n1.1.1.1\r\n",
			[]string{"8.8.8.8", "1.1.1.1"},
		},
		{
			"empty input",
			"",
			nil,
		},
		{
			"only comments and blanks",
			"# comment\n\n# another\n",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromReader(strings.NewReader(tt.input))
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d\ngot: %v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
