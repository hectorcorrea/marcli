package marc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewFieldFilters(t *testing.T) {
	t.Parallel()

	newFieldFiltersTests := []struct {
		name      string
		fieldsStr string
		result    FieldFilters
	}{
		{name: "empty fieldsStr", fieldsStr: "", result: FieldFilters{Fields: nil}},
		{name: "one field no subfields", fieldsStr: "700", result: FieldFilters{Fields: []FieldFilter{{Tag: "700"}}}},
		{name: "one field with subfields", fieldsStr: "700h", result: FieldFilters{Fields: []FieldFilter{{Tag: "700", Subfields: "h"}}}},
		{name: "one field with subfields", fieldsStr: "700h,245ahc", result: FieldFilters{Fields: []FieldFilter{{Tag: "700", Subfields: "h"}, {Tag: "245", Subfields: "ahc"}}}},
	}

	for _, tt := range newFieldFiltersTests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFieldFilters(tt.fieldsStr)

			if !cmp.Equal(got, tt.result) {
				t.Errorf("expected %q, got %q", got, tt.result)
				t.Errorf(cmp.Diff(got, tt.result))
			}
		})
	}
}

func TestNewFieldFilter(t *testing.T) {
	t.Parallel()

	t.Run("short field string", func(t *testing.T) {
		want, _ := FieldFilter{}, ErrInvalidFieldString

		got, err := NewFieldFilter("70")

		if want != got {
			t.Errorf("expected %q, got %q", want, got)
		}

		if !errors.Is(err, ErrInvalidFieldString) {
			t.Errorf("expected %q, got %q", ErrInvalidFieldString, err)
		}
	})

	newFieldFilterTests := []struct {
		name     string
		fieldStr string
		filter   FieldFilter
		err      error
	}{
		{name: "field string without subfield", fieldStr: "700", filter: FieldFilter{Tag: "700"}, err: nil},
		{name: "field string with one subfield", fieldStr: "700h", filter: FieldFilter{Tag: "700", Subfields: "h"}, err: nil},
		{name: "field string with multiple subfields", fieldStr: "245ahc", filter: FieldFilter{Tag: "245", Subfields: "ahc"}, err: nil},
	}

	for _, tt := range newFieldFilterTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFieldFilter(tt.fieldStr)

			if err != nil {
				t.Errorf("unexpected error: %q", err)
			}

			if !cmp.Equal(got, tt.filter) {
				t.Errorf("expected %q, got %q", got, tt.filter)
			}
		})
	}
}

func TestIncludeField(t *testing.T) {
	t.Parallel()

	includeFieldTests := []struct {
		name         string
		field        string
		fieldFilters FieldFilters
		result       bool
	}{
		{
			name:         "field included",
			field:        "245",
			fieldFilters: FieldFilters{Fields: []FieldFilter{{Tag: "700", Subfields: "h"}, {Tag: "245", Subfields: "ahc"}}},
			result:       true,
		},
		{
			name:         "field not included",
			field:        "LDR",
			fieldFilters: FieldFilters{Fields: []FieldFilter{{Tag: "700", Subfields: "h"}, {Tag: "245", Subfields: "ahc"}}},
			result:       false,
		},
	}

	for _, tt := range includeFieldTests {
		got := tt.fieldFilters.IncludeField(tt.field)

		if !got == tt.result {
			t.Errorf("expected IncludeField() called on %q to return true", tt.fieldFilters)
		}
	}
}

func TestIncludeLeader(t *testing.T) {
	t.Parallel()

	fieldFilters := FieldFilters{Fields: []FieldFilter{{Tag: "LDR", Subfields: ""}, {Tag: "245", Subfields: "ahc"}}}

	if !fieldFilters.IncludeLeader() {
		t.Errorf("expected IncludeLeader() called on %q to return true", fieldFilters)
	}
}

func TestFieldFiltersString(t *testing.T) {
	t.Parallel()

	fieldFilters := FieldFilters{Fields: []FieldFilter{{Tag: "LDR", Subfields: ""}, {Tag: "245", Subfields: "ahc"}}}

	want := "Filters {\r\n\tTag: LDR\r\n\tTag: 245 subfields: ahc\r\n}\r\n"
	got := fieldFilters.String()

	if want != got {
		t.Errorf("expected %q, got %q", want, got)
	}
}
