package toposort_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/onur1/toposort"
)

func TestGraph(t *testing.T) {
	testCases := []struct {
		desc   string
		sorted []string
		data   map[string]string
		err    error
	}{
		{
			desc:   "example",
			sorted: []string{"Jonas", "Sophie", "Nick", "Barbara"},
			data: map[string]string{
				"Barbara": "Nick",
				"Nick":    "Sophie",
				"Sophie":  "Jonas",
			},
		},
		{
			desc:   "single row",
			sorted: []string{"Jonas", "Sophie"},
			data: map[string]string{
				"Sophie": "Jonas",
			},
		},
		{
			desc: "cyclic",
			data: map[string]string{
				"Barbara": "Nick",
				"Nick":    "Sophie",
				"Sophie":  "Jonas",
				"Jonas":   "Barbara",
			},
			err: toposort.ErrCircular,
		},
		{
			desc: "multiple cyclic",
			data: map[string]string{
				"Barbara": "Nick",
				"Nick":    "Sophie",
				"Sophie":  "Jonas",
				"Jonas":   "Barbara",
				"Daniel":  "Ruby",
				"Jason":   "Daniel",
				"Ruby":    "Jason",
			},
			err: toposort.ErrCircular,
		},
		{
			desc:   "single case sensitive",
			sorted: []string{"Jonas", "jONas"},
			data: map[string]string{
				"jONas": "Jonas",
			},
		},
		{
			desc: "multiple roots",
			data: map[string]string{
				"Barbara": "Nick",
				"Nick":    "Sophie",
				"Sophie":  "Jonas",
				"Ruby":    "Daniel",
			},
			err: toposort.ErrMultipleRoots,
		},
	}
	for _, tt := range testCases {
		tt := tt

		t.Run(tt.desc, func(t *testing.T) {
			sorted, err := toposort.Sort(tt.data)
			if tt.err == nil {
				if err != nil {
					t.Fatal(err)
				}
				if reflect.DeepEqual(sorted, tt.sorted) {
					return
				}
				t.Fatalf("expected sorted value %+v != %+v", tt.sorted, sorted)
			}
			if errors.Is(err, tt.err) {
				return
			}
			t.Fatalf("expected error %v != %v", tt.err, err)
		})
	}
}
