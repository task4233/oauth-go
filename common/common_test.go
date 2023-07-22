package common_test

import (
	"testing"

	"github.com/task4233/oauth-go/common"
)

func TestAreTwoUnorderedSlicesSame(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		s    []string
		t    []string
		want bool
	}{
		"same slice": {
			s:    []string{"a", "b", "c"},
			t:    []string{"a", "b", "c"},
			want: true,
		},
		"same slice but different order": {
			s:    []string{"a", "b", "c"},
			t:    []string{"c", "b", "a"},
			want: true,
		},
		"same slice but different length": {
			s:    []string{"a", "b", "c"},
			t:    []string{"a", "b"},
			want: false,
		},
		"different slice": {
			s:    []string{"a", "b", "c"},
			t:    []string{"a", "b", "d"},
			want: false,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := common.AreTwoUnorderedSlicesSame(tt.s, tt.t)
			if got != tt.want {
				t.Fatalf("want: %v, but got: %v", tt.want, got)
			}
		})
	}
}
