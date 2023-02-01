package marc

import (
	"bufio"
	"encoding/xml"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewMarcFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		path  string
		isXML bool
	}{
		{name: "binary", path: "testdata/test_1a.mrc", isXML: false},
		{name: "XML", path: "testdata/test_10.xml", isXML: true},
		{name: "binary no extension", path: "testdata/bad", isXML: true},
		{name: "misleading extension", path: "testdata/test_bad.xml", isXML: false},
	}

	for _, tt := range tests {
		var want MarcFile
		file := setUpTestFile(tt.path, t)
		got := NewMarcFile(file)
		if tt.isXML {
			want = MarcFile{
				decoder: xml.NewDecoder(file),
				isXML:   true,
			}

			if !(want.isXML && got.isXML) {
				t.Error("expected struct field isXML to be true for both files")
			}
		} else {
			want = MarcFile{
				scanner: bufio.NewScanner(file),
			}

			opt := cmp.Comparer(func(a, b MarcFile) bool {
				return cmp.Equal(a.scanner.Bytes(), b.scanner.Bytes())
			})

			if !cmp.Equal(want, got, opt) {
				t.Errorf("expected %v, got %v", want, got)
			}
		}
	}
}

func TestRecord(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		path  string
		isXML bool
	}{
		{name: "binary", path: "testdata/test_1a.mrc", isXML: false},
		{name: "XML", path: "testdata/test_10.xml", isXML: true},
	}

	for _, tt := range tests {
		want := newRecord(tt.isXML, t)
		file := setUpTestFile(tt.path, t)

		f := NewMarcFile(file)
		f.Scan()
		got, err := f.Record()
		if err != nil {
			t.Fatalf("problem calling Record on MarcFile: %s", err)
		}

		opt := cmp.AllowUnexported(Leader{})

		if !cmp.Equal(want, got, opt) {
			t.Errorf("expected %q, got %q", want, got)
			t.Error(cmp.Diff(want, got, opt))
		}
	}
}

func setUpTestFile(path string, t *testing.T) *os.File {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	return file
}

func newRecord(isXML bool, t *testing.T) Record {
	t.Helper()

	var rawData []byte

	if isXML {
		rawData = []byte("Raw data not supported in XML format\n")
	} else {
		rawData = []byte("01805nam a2200385 i 4500001001200000005001700012006001900029007000700048008004100055040002300096042000800119043001200127074002000139086001700159100005500176245021100231260008500442336002200527337002400549338003300573440003900606500005400645504004100699538015700740650002000897650002000917700002100937776020100958856006501159907003501224998004501259910001201304910002801316945007501344ocm5717594020041206161421.0m        d f      cr cn-041206s1976    dcua    sb   f000 0 eng c  aGPOcGPOdMvIdMvI  apcc  an-us---  a0620-A (online)0 aI 19.4/2:7351 aSwanson, Vernon E.q(Vernon Emmanuel),d1922-1992.10aGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coalh[electronic resource] /cby Vernon E. Swanson and Claude Huffman, Jr.  a[Washington, D.C.] :bU.S. Dept. of the Interior, U.S. Geological Survey,c1976.  atext2rdacontent.  acomputer2rdamedia.  aonline resource2rdacarrier. 0aGeological Survey circular ;v735.  aTitle from title screen (viewed on Dec. 06, 2004)  aIncludes bibliographical references.  aMode of access: Internet from the USGS Web site. Address as of 12/06/04: http://pubs.usgs.gov/circ/c735/index.htm; current access is available via PURL. 0aCoalxAnalysis. 0aCoalxSampling.1 aHuffman, Claude.1 aSwanson, Vernon Emanuel,d1922-tGuidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coalhiv, 11 p.w(OCoLC)2331861.40uhttp://purl.access.gpo.gov/GPO/LPS56007zView online version  a.b37991760b04-08-17c07-26-05  aes001b07-26-05cmdae-fenggdcuh0i1  aMARCIVE  aHathi Trust report None  g0j0lesb  onp$0.00q r s-t255u0v0w0x0y.i138993579z07-26-05")
	}

	return Record{
		Data: rawData,
		Fields: []Field{
			{
				Tag:        "001",
				Value:      "ocm57175940",
				Indicator1: "",
				Indicator2: "",
			},
			{
				Tag:        "005",
				Value:      "20041206161421.0",
				Indicator1: "",
				Indicator2: "",
			},
			{
				Tag:        "006",
				Value:      "m        d f      ",
				Indicator1: "",
				Indicator2: "",
			},
			{
				Tag:        "007",
				Value:      "cr cn-",
				Indicator1: "",
				Indicator2: "",
			},
			{
				Tag:        "008",
				Value:      "041206s1976    dcua    sb   f000 0 eng c",
				Indicator1: "",
				Indicator2: "",
			},
			{
				Tag:        "040",
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
			{
				Tag:        "042",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "pcc",
					},
				},
			},
			{
				Tag:        "043",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "n-us---",
					},
				},
			},
			{
				Tag:        "074",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "0620-A (online)",
					},
				},
			},
			{
				Tag:        "086",
				Indicator1: "0",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "I 19.4/2:735",
					},
				},
			},
			{
				Tag:        "100",
				Indicator1: "1",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Swanson, Vernon E.",
					},
					{
						Code:  "q",
						Value: "(Vernon Emmanuel),",
					},
					{
						Code:  "d",
						Value: "1922-1992.",
					},
				},
			},
			{
				Tag:        "245",
				Indicator1: "1",
				Indicator2: "0",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Guidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coal",
					},
					{
						Code:  "h",
						Value: "[electronic resource] /",
					},
					{
						Code:  "c",
						Value: "by Vernon E. Swanson and Claude Huffman, Jr.",
					},
				},
			},
			{
				Tag:        "260",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "[Washington, D.C.] :",
					},
					{
						Code:  "b",
						Value: "U.S. Dept. of the Interior, U.S. Geological Survey,",
					},
					{
						Code:  "c",
						Value: "1976.",
					},
				},
			},
			{
				Tag:        "336",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "text",
					},
					{
						Code:  "2",
						Value: "rdacontent.",
					},
				},
			},
			{
				Tag:        "337",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "computer",
					},
					{
						Code:  "2",
						Value: "rdamedia.",
					},
				},
			},
			{
				Tag:        "338",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "online resource",
					},
					{
						Code:  "2",
						Value: "rdacarrier.",
					},
				},
			},
			{
				Tag:        "440",
				Indicator1: " ",
				Indicator2: "0",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Geological Survey circular ;",
					},
					{
						Code:  "v",
						Value: "735.",
					},
				},
			},
			{
				Tag:        "500",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Title from title screen (viewed on Dec. 06, 2004)",
					},
				},
			},
			{
				Tag:        "504",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Includes bibliographical references.",
					},
				},
			},
			{
				Tag:        "538",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Mode of access: Internet from the USGS Web site. Address as of 12/06/04: http://pubs.usgs.gov/circ/c735/index.htm; current access is available via PURL.",
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
			{
				Tag:        "700",
				Indicator1: "1",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Huffman, Claude.",
					},
				},
			},
			{
				Tag:        "776",
				Indicator1: "1",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Swanson, Vernon Emanuel,",
					},
					{
						Code:  "d",
						Value: "1922-",
					},
					{
						Code:  "t",
						Value: "Guidelines for sample collecting and analytical methods used in the U.S. Geological Survey for determining chemical composition of coal",
					},
					{
						Code:  "h",
						Value: "iv, 11 p.",
					},
					{
						Code:  "w",
						Value: "(OCoLC)2331861.",
					},
				},
			},
			{
				Tag:        "856",
				Indicator1: "4",
				Indicator2: "0",
				SubFields: []SubField{
					{
						Code:  "u",
						Value: "http://purl.access.gpo.gov/GPO/LPS56007",
					},
					{
						Code:  "z",
						Value: "View online version",
					},
				},
			},
			{
				Tag:        "907",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: ".b37991760",
					},
					{
						Code:  "b",
						Value: "04-08-17",
					},
					{
						Code:  "c",
						Value: "07-26-05",
					},
				},
			},
			{
				Tag:        "998",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "es001",
					},
					{
						Code:  "b",
						Value: "07-26-05",
					},
					{
						Code:  "c",
						Value: "m",
					},
					{
						Code:  "d",
						Value: "a",
					},
					{
						Code:  "e",
						Value: "-",
					},
					{
						Code:  "f",
						Value: "eng",
					},
					{
						Code:  "g",
						Value: "dcu",
					},
					{
						Code:  "h",
						Value: "0",
					},
					{
						Code:  "i",
						Value: "1",
					},
				},
			},
			{
				Tag:        "910",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "MARCIVE",
					},
				},
			},
			{
				Tag:        "910",
				Indicator1: " ",
				Indicator2: " ",
				SubFields: []SubField{
					{
						Code:  "a",
						Value: "Hathi Trust report None",
					},
				},
			},
			{
				Tag:        "945",
				Value:      "",
				Indicator1: " ",
				Indicator2: " ",
				SubFields:  []SubField{{Code: "g", Value: "0"}, {Code: "j", Value: "0"}, {Code: "l", Value: "esb  "}, {Code: "o", Value: "n"}, {Code: "p", Value: "$0.00"}, {Code: "q", Value: " "}, {Code: "r", Value: " "}, {Code: "s", Value: "-"}, {Code: "t", Value: "255"}, {Code: "u", Value: "0"}, {Code: "v", Value: "0"}, {Code: "w", Value: "0"}, {Code: "x", Value: "0"}, {Code: "y", Value: ".i138993579"}, {Code: "z", Value: "07-26-05"}},
			},
		},
		Leader: Leader{
			raw:           []byte("01805nam a2200385 i 4500"),
			dataOffset:    385,
			Status:        byte('n'),
			Type:          byte('a'),
			BibLevel:      byte('m'),
			Control:       byte(' '),
			EncodingLevel: byte(' '),
			Form:          byte('i'),
			Multipart:     byte(' '),
		},
	}
}
