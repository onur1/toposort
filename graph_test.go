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
			desc: "single cyclic case insensitive",
			data: map[string]string{
				"jONas": "Jonas",
			},
			err: toposort.ErrCircular,
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
		{
			desc: "invalid name 1",
			data: map[string]string{
				"Barbara": "Nick1",
			},
			err: toposort.ErrInvalidName,
		},
		{
			desc: "invalid name 2",
			data: map[string]string{
				"a": "Nick",
			},
			err: toposort.ErrInvalidName,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			g, err := toposort.NewGraph(tt.data)
			if err != nil {
				if tt.err == nil {
					t.Fatal(err)
				}
				if !errors.Is(err, tt.err) {
					t.Fatalf("expected error %v != %v", tt.err, err)
				}
			}
			if tt.sorted != nil {
				sortedIDs := g.SortedIDs()
				if !reflect.DeepEqual(sortedIDs, tt.sorted) {
					t.Fatalf("expected sorted value %+v != %+v", tt.sorted, sortedIDs)
				}
			}
		})
	}
}
