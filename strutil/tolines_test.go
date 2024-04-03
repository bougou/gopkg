package strutil

import "testing"

func Test_ToLines(t *testing.T) {
	tests := []struct {
		name      string
		multiline string
		out       []string
	}{
		{
			name:      "#1",
			multiline: "a\nb\nc\n",
			out:       []string{"a", "b", "c"},
		},
		{
			name:      "#2",
			multiline: "a\nb\nc\n\n",
			out:       []string{"a", "b", "c", ""},
		},
		{
			name: "#3",
			multiline: `

a

b
c



`,
			out: []string{"", "", "a", "", "b", "c", "", "", ""},
		},
	}

	for _, tt := range tests {
		got := ToLines(tt.multiline)
		want := tt.out
		if len(got) != len(want) {
			t.Errorf("case %s: lines number not matched, got: %v, expected: %v", tt.name, got, want)
			return
		}

		for i, line := range got {
			if line != want[i] {
				t.Errorf("case %s: item not matched, got: %v, expected: %v", tt.name, got, want)
			}
		}
	}
}
