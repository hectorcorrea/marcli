package marc

import (
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRecordContains(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	containsTests := []struct {
		name             string
		searchValue      string
		searchFieldsList []string
		result           bool
	}{
		{name: "empty searchValue", searchValue: "", searchFieldsList: []string{"650"}, result: true},
		{name: "empty searchFieldsList", searchValue: "Coal", searchFieldsList: []string{}, result: true},
		{name: "non-empty searchFieldsList", searchValue: "Coal", searchFieldsList: []string{"650"}, result: true},
		{name: "case insensitive search non-empty searchFieldsList", searchValue: "coal", searchFieldsList: []string{"650"}, result: true},
		{name: "empty searchFieldsList", searchValue: "Pizza", searchFieldsList: []string{}, result: false},
		{name: "non-empty searchFieldsList", searchValue: "Coal", searchFieldsList: []string{"260"}, result: false},
	}

	for _, tt := range containsTests {
		got := record.Contains(tt.searchValue, tt.searchFieldsList)
		if !got == tt.result {
			t.Errorf("expected record not to contain %q", tt.searchValue)
		}
	}
}

func TestControlNumber(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	want := "ocm57175940"

	got := record.ControlNum()

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestHasFields(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	fieldFilters := FieldFilters{
		Fields: []FieldFilter{
			{Tag: "650", Subfields: "ax"}},
	}

	if !record.HasFields(fieldFilters) {
		t.Errorf("expected to have field(s) %s", fieldFilters)
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	filterTests := []struct {
		name    string
		include FieldFilters
		exclude FieldFilters
		result  []Field
	}{
		{name: "empty include, empty exclude", include: FieldFilters{}, exclude: FieldFilters{}, result: record.Fields},
		{name: "include one tag no subfields, empty exclude", include: FieldFilters{Fields: []FieldFilter{{Tag: "650", Subfields: ""}}}, exclude: FieldFilters{}, result: record.FieldsByTag("650")},
		{
			name:    "include one tag one subfield, empty exclude",
			include: FieldFilters{Fields: []FieldFilter{{Tag: "650", Subfields: "a"}}},
			exclude: FieldFilters{},
			result:  filterOnSubFields(record.FieldsByTag("650"), "a", t),
		},
		{
			name:    "empty include, long exclude",
			include: FieldFilters{},
			exclude: FieldFilters{Fields: []FieldFilter{
				{Tag: "001", Subfields: ""},
				{Tag: "005", Subfields: ""},
				{Tag: "006", Subfields: ""},
				{Tag: "007", Subfields: ""},
				{Tag: "008", Subfields: ""},
				{Tag: "040", Subfields: ""},
				{Tag: "042", Subfields: ""},
				{Tag: "043", Subfields: ""},
				{Tag: "074", Subfields: ""},
				{Tag: "086", Subfields: ""},
				{Tag: "100", Subfields: ""},
				{Tag: "245", Subfields: ""},
				{Tag: "260", Subfields: ""},
				{Tag: "336", Subfields: ""},
				{Tag: "337", Subfields: ""},
				{Tag: "338", Subfields: ""},
				{Tag: "440", Subfields: ""},
				{Tag: "500", Subfields: ""},
				{Tag: "504", Subfields: ""},
				{Tag: "538", Subfields: ""},
				{Tag: "700", Subfields: ""},
				{Tag: "776", Subfields: ""},
				{Tag: "856", Subfields: ""},
				{Tag: "907", Subfields: ""},
				{Tag: "998", Subfields: ""},
				{Tag: "910", Subfields: ""},
				{Tag: "945", Subfields: ""},
			}},
			result: record.FieldsByTag("650"),
		},
	}

	for _, tt := range filterTests {
		t.Run(tt.name, func(t *testing.T) {
			got := record.Filter(tt.include, tt.exclude)

			if !cmp.Equal(tt.result, got) {
				t.Errorf("expected %q\n\ngot %q", tt.result, got)
			}
		})
	}
}

func TestRecordRaw(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	want := []byte("01805nam a2200385 i 4500001001200000005001700012006001900029007000700048008004100055040002300096042000800119043001200127074002000139086001700159100005500176245021100231260008500442336002200527337002400549338003300573440003900606500005400645504004100699538015700740650002000897650002000917700002100937776020100958856006501159907003501224998004501259910001201304910002801316945007501344ocm5717594020041206161421.0m        d f      cr cn-041206s1976    dcua    sb   f000 0 eng c  aGPOcGPOdMvIdMvI  apcc  an-us---  a0620-A (online)0 aI 19.4/2:7351 aSwanson, Vernon E.q(Vernon Emmanuel),d1922-1992.10aGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coalh[electronic resource] /cby Vernon E. Swanson and Claude Huffman, Jr.  a[Washington, D.C.] :bU.S. Dept. of the Interior, U.S. Geological Survey,c1976.  atext2rdacontent.  acomputer2rdamedia.  aonline resource2rdacarrier. 0aGeological Survey circular ;v735.  aTitle from title screen (viewed on Dec. 06, 2004)  aIncludes bibliographical references.  aMode of access: Internet from the USGS Web site. Address as of 12/06/04: http://pubs.usgs.gov/circ/c735/index.htm; current access is available via PURL. 0aCoalxAnalysis. 0aCoalxSampling.1 aHuffman, Claude.1 aSwanson, Vernon Emanuel,d1922-tGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coalhiv, 11 p.w(OCoLC)2331861.40uhttp://purl.access.gpo.gov/GPO/LPS56007zView online version  a.b37991760b04-08-17c07-26-05  aes001b07-26-05cmdae-fenggdcuh0i1  aMARCIVE  aHathi Trust report None  g0j0lesb  onp$0.00q r s-t255u0v0w0x0y.i138993579z07-26-05")

	got := record.Raw()

	if !bytes.Equal(got, want) {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestRecordString(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	want := "Leader: =LDR  01805nam a2200385 i 4500"

	got := record.String()

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestFieldByTag(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	want := []Field{
		{
			Tag:        "650",
			Indicator1: " ",
			Indicator2: "0",
			SubFields: []SubField{
				{
					Code:  "a",
					Value: "Coal",
				},
				{
					Code:  "x",
					Value: "Analysis.",
				},
			},
		},
		{
			Tag:        "650",
			Indicator1: " ",
			Indicator2: "0",
			SubFields: []SubField{
				{
					Code:  "a",
					Value: "Coal",
				},
				{
					Code:  "x",
					Value: "Sampling.",
				},
			},
		},
	}

	got := record.FieldsByTag("650")

	if !cmp.Equal(want, got) {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGetValue(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	getValueTests := []struct {
		name     string
		tag      string
		subfield string
		result   string
	}{
		{name: "control field no subfield", tag: "001", subfield: "", result: "ocm57175940"},
		{name: "control field with (ignored) subfield", tag: "001", subfield: "b", result: "ocm57175940"},
		{name: "data field with subfield", tag: "650", subfield: "x", result: "Analysis."},
		{name: "data field without subfield", tag: "650", subfield: "", result: "=650  \\0$aCoal$xAnalysis."},
	}

	for _, tt := range getValueTests {
		t.Run(tt.name, func(t *testing.T) {
			got := record.GetValue(tt.tag, tt.subfield)

			if got != tt.result {
				t.Errorf("expected %q, got %q", tt.result, got)
			}
		})
	}
}

func TestGetValues(t *testing.T) {
	t.Parallel()

	record := setUpTestRecord("testdata/test_1a.mrc", t)

	getValueTests := []struct {
		name     string
		tag      string
		subfield string
		result   []string
	}{
		{name: "control field no subfield", tag: "001", subfield: "", result: []string{"=001  ocm57175940"}},
		{name: "control field with (ignored) subfield", tag: "001", subfield: "b", result: []string{}},
		{name: "data field with subfield", tag: "650", subfield: "x", result: []string{"Analysis.", "Sampling."}},
		{name: "data field without subfield", tag: "650", subfield: "", result: []string{"=650  \\0$aCoal$xAnalysis.", "=650  \\0$aCoal$xSampling."}},
	}

	for _, tt := range getValueTests {
		t.Run(tt.name, func(t *testing.T) {
			got := record.GetValues(tt.tag, tt.subfield)

			if !cmp.Equal(got, tt.result) {
				t.Errorf("expected %q, got %q", tt.result, got)
			}
		})
	}
}

func setUpTestRecord(path string, t *testing.T) Record {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("problem opening test data file %q: %v", path, err)
	}

	f := NewMarcFile(file)
	f.Scan()
	record, err := f.Record()
	if err != nil {
		t.Fatalf("problem getting record: %v", err)
	}

	return record
}

func filterOnSubFields(fields []Field, subfield string, t *testing.T) []Field {
	t.Helper()

	outFields := []Field{}

	for _, field := range fields {
		subfields := field.GetSubFields(subfield)

		outField := Field{
			Tag:        field.Tag,
			Value:      field.Value,
			Indicator1: field.Indicator1,
			Indicator2: field.Indicator2,
			SubFields:  subfields,
		}

		outFields = append(outFields, outField)
	}

	return outFields
}
