package marc

import (
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const dirs = "001001200000005001700012006001900029007000700048008004100055040002300096042000800119043001200127074002000139086001700159100005500176245021100231260008500442336002200527337002400549338003300573440003900606500005400645504004100699538015700740650002000897650002000917700002100937776020100958856006501159907003501224998004501259910001201304910002801316945007501344"
const data = "ocm5717594020041206161421.0m        d f      cr cn-041206s1976    dcua    sb   f000 0 eng c  aGPOcGPOdMvIdMvI  apcc  an-us---  a0620-A (online)0 aI 19.4/2:7351 aSwanson, Vernon E.q(Vernon Emmanuel),d1922-1992.10aGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coalh[electronic resource] /cby Vernon E. Swanson and Claude Huffman, Jr.  a[Washington, D.C.] :bU.S. Dept. of the Interior, U.S. Geological Survey,c1976.  atext2rdacontent.  acomputer2rdamedia.  aonline resource2rdacarrier. 0aGeological Survey circular ;v735.  aTitle from title screen (viewed on Dec. 06, 2004)  aIncludes bibliographical references.  aMode of access: Internet from the USGS Web site. Address as of 12/06/04: http://pubs.usgs.gov/circ/c735/index.htm; current access is available via PURL. 0aCoalxAnalysis. 0aCoalxSampling.1 aHuffman, Claude.1 aSwanson, Vernon Emanuel,d1922-tGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coalhiv, 11 p.w(OCoLC)2331861.40uhttp://purl.access.gpo.gov/GPO/LPS56007zView online version  a.b37991760b04-08-17c07-26-05  aes001b07-26-05cmdae-fenggdcuh0i1  aMARCIVE  aHathi Trust report None  g0j0lesb  onp$0.00q r s-t255u0v0w0x0y.i138993579z07-26-05"

func TestMakeField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		offset int
		want   Field
	}{
		{
			name:   "control field",
			offset: 0,
			want: Field{
				Tag:        "001",
				Value:      "ocm57175940",
				Indicator1: "",
				Indicator2: "",
			},
		},
		{
			name: "data field",
			// 60 is the index of the first data field in the example data, 040
			offset: 60,
			want: Field{
				Tag:        "040",
				Value:      "",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "GPO",
					},
					{
						Code:  "c",
						Value: "GPO",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag, fdata := setUpDirsAndData([]byte(dirs), []byte(data), tt.offset, t)
			got, _ := MakeField(tag, fdata)
			compareFields(tt.want, got, t)
		})
	}
}

func TestIsControlField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  Field
		result bool
	}{
		{
			name: "returns true for a control field",
			input: Field{
				Tag:        "001",
				Value:      "ocm57175940",
				Indicator1: "",
				Indicator2: "",
			},
			result: true,
		},
		{
			name: "returns false for a data field",
			input: Field{
				Tag:        "040",
				Value:      "",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "GPO",
					},
					{
						Code:  "c",
						Value: "GPO",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
				},
			},
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.IsControlField() != tt.result {
				t.Errorf("expected IsControlField() to return %v for %v", tt.result, tt.input)
			}
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  Field
		arg    string
		regEx  string
		result bool
	}{
		{
			name: "control field contains string",
			input: Field{
				Tag:        "001",
				Value:      "ocm57175940",
				Indicator1: "",
				Indicator2: "",
			},
			arg:    "759",
			result: true,
		},
		{
			name: "control field does not contain string",
			input: Field{
				Tag:        "001",
				Value:      "ocm57175940",
				Indicator1: "",
				Indicator2: "",
			},
			arg:    "abc",
			result: false,
		},
		{
			name: "data field contains string",
			input: Field{
				Tag:        "040",
				Value:      "",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "GPO",
					},
					{
						Code:  "c",
						Value: "GPO",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
				},
			},
			arg:    "MvI",
			result: true,
		},
		{
			name: "data field does not contain string",
			input: Field{
				Tag:        "040",
				Value:      "",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "GPO",
					},
					{
						Code:  "c",
						Value: "GPO",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
				},
			},
			arg:    "abc",
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Add test for regEx searches
			if tt.input.Contains(tt.arg, tt.regEx) != tt.result {
				t.Errorf("expected Contains() call on %v to return %v for %v", tt.input, tt.result, tt.arg)
			}
		})
	}
}

func TestGetSubFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		want  []SubField
		field Field
		input string
	}{
		{
			name: "calling GetSubFields on a data field returns slice of matching SubFields",
			want: []SubField{
				{
					Code:  "d",
					Value: "MvI",
				},
				{
					Code:  "d",
					Value: "MvI",
				},
			},
			field: Field{
				Tag:        "040",
				Value:      "",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "GPO",
					},
					{
						Code:  "c",
						Value: "GPO",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
					{
						Code:  "d",
						Value: "MvI",
					},
				},
			},
			input: "d",
		},
		{
			name: "calling GetSubFields on a data field returns slice of matching SubFields",
			want: []SubField{},
			field: Field{
				Tag:        "001",
				Value:      "ocm57175940",
				Indicator1: "",
				Indicator2: "",
			},
			input: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.field.GetSubFields(tt.input)

			if !cmp.Equal(tt.want, got) {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func setUpDirsAndData(dirs []byte, data []byte, offset int, t *testing.T) (string, []byte) {
	t.Helper()

	dirs = dirs[offset:]
	tag := string(dirs[:tagEnd])
	length, _ := strconv.Atoi(string(dirs[lengthOfFieldStart:lengthOfFieldEnd]))
	begin, _ := strconv.Atoi(string(dirs[startCharPosStart:startCharPosEnd]))
	fdata := data[begin : begin+length-1]

	return tag, fdata
}

func (a Field) Equal(b Field) bool {
	less := func(a, b string) bool { return a < b }
	subfieldsEqual := cmp.Equal(a.SubFields, b.SubFields, cmpopts.SortSlices(less))
	return a.Tag == b.Tag && a.Value == b.Value && a.Indicator1 == b.Indicator1 &&
		a.Indicator2 == b.Indicator2 &&
		subfieldsEqual
}

func compareFields(want, got Field, t *testing.T) {
	t.Helper()

	if !cmp.Equal(want, got) {
		t.Errorf("expected %q, got %q", want, got)
		t.Error(cmp.Diff(want, got))
	}
}
