package marc

type XmlRecord struct {
	Leader        string `xml:"leader"`
	ControlFields []struct {
		Value string `xml:",chardata"`
		Tag   string `xml:"tag,attr"`
	} `xml:"controlfield"`
	DataFields []struct {
		Tag       string `xml:"tag,attr"`
		Ind1      string `xml:"ind1,attr"`
		Ind2      string `xml:"ind2,attr"`
		SubFields []struct {
			Value string `xml:",chardata"`
			Code  string `xml:"code,attr"`
		} `xml:"subfield"`
	} `xml:"datafield"`
}
