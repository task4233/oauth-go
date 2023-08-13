package common_test

import (
	"testing"

	"github.com/task4233/oauth/common"
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

func TestConstructWithQueries(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		uri             string
		queryParameters map[string]string
		wantURL         string
		wantErr         bool
	}{
		"ok:no query": {
			uri:             "https://task4233.dev",
			queryParameters: map[string]string{},
			wantURL:         "https://task4233.dev",
			wantErr:         false,
		},
		"ok:one query": {
			uri:             "https://task4233.dev",
			queryParameters: map[string]string{"a": "b"},
			wantURL:         "https://task4233.dev?a=b",
			wantErr:         false,
		},
		"ok:two queries": {
			uri:             "https://task4233.dev",
			queryParameters: map[string]string{"a": "b", "c": "d"},
			wantURL:         "https://task4233.dev?a=b&c=d",
			wantErr:         false,
		},
		"ok:one query with existing query": {
			uri:             "https://task4233.dev?a=b",
			queryParameters: map[string]string{"c": "d"},
			wantURL:         "https://task4233.dev?a=b&c=d",
			wantErr:         false,
		},
		"ok:two queries with existing query": {
			uri:             "https://task4233.dev?a=b",
			queryParameters: map[string]string{"c": "d", "e": "f"},
			wantURL:         "https://task4233.dev?a=b&c=d&e=f",
			wantErr:         false,
		},
		"ok:one query with existing query with same key": {
			uri:             "https://task4233.dev?a=b",
			queryParameters: map[string]string{"a": "c"},
			wantURL:         "https://task4233.dev?a=c",
			wantErr:         false,
		},
		"ok:two queries with existing query with same key": {
			uri:             "https://task4233.dev?a=b",
			queryParameters: map[string]string{"a": "c", "d": "e"},
			wantURL:         "https://task4233.dev?a=c&d=e",
			wantErr:         false,
		},
		"ng:invalid uri": {
			uri:             "https://task4233.dev:invalid",
			queryParameters: map[string]string{},
			wantURL:         "",
			wantErr:         true,
		},
		"ng:invalid query": {
			uri:             "https://task4233.dev",
			queryParameters: map[string]string{"a": "b", "c": "d&"},
			wantURL:         "",
			wantErr:         true,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := common.ConstructURLWithQueries(tt.uri, tt.queryParameters)
			if tt.wantErr != (err != nil) {
				t.Fatalf("unexpected error, want: %v, got: %v", tt.wantErr, err)
			}
			if err != nil {
				return
			}
			if tt.wantURL != got {
				t.Fatalf("unexpected wantURL, want: %v, but got: %v", tt.wantURL, got)
			}
		})
	}
}
