// oligos_test.go
package oligos

import (
	//"fmt"
	//"strconv"
	//	"strings"

	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	//"github.com/antha-lang/antha/internal/github.com/tealeg/xlsx"
)

// simple reverse complement check to test testing methodology initially
type testpair struct {
	//BasicMeltingTemp test
	sequence    wtype.DNASequence
	meltingtemp wunit.Temperature
	// FWDoligoTest
	mintemp             wunit.Temperature
	maxtemp             wunit.Temperature
	maxGCcontent        float64
	minlength           int
	maxlength           int
	seqstoavoid         []string
	overlapthreshold    int
	outputoligoseq      string
	calculatedGCcontent float64

	// tests of OverlapCheck
	primer1        string
	primer2        string
	overlapPercent float64
	overlapNumber  int
	overlapSeq     string
}

var meltingtemptests = []testpair{

	{sequence: wtype.MakeSingleStrandedDNASequence("whatever", "AAAAAAAAAAAAAAAAAAA"),
		meltingtemp: wunit.NewTemperature(29.511, "C")},
}

var oligotests = []testpair{

	{sequence: wtype.MakeSingleStrandedDNASequence("sfGFP", "ATGAGCAAAGGAGAAGAACTTTTCACTGGAGTTGTCCCAATTCTTGTTGAATTAGATGGTGATGTTAATGGGCACAAATTTTCTGTCCGTGGAGAGGGTGAAGGTGATGCTACAAACGGAAAACTCACCCTTAAATTTATTTGCACTACTGGAAAACTACCTGTTCCATGGCCAACACTTGTCACTACTCTGACCTATGGTGTTCAATGCTTTTCCCGTTATCCGGATCACATGAAACGGCATGACTTTTTCAAGAGTGCCATGCCCGAAGGTTATGTACAGGAACGCACTATATCTTTCAAAGATGACGGGACCTACAAGACGCGTGCTGAAGTCAAGTTTGAAGGTGATACCCTTGTTAATCGTATCGAGTTAAAAGGTATTGATTTTAAAGAAGATGGAAACATTCTCGGACACAAACTCGAGTACAACTTTAACTCACACAATGTATACATCACGGCAGACAAACAAAAGAATGGAATCAAAGCTAACTTCAAAATTCGCCACAACGTTGAAGATGGTTCCGTTCAACTAGCAGACCATTATCAACAAAATACTCCAATTGGCGATGGCCCTGTCCTTTTACCAGACAACCATTACCTGTCGACACAATCTGTCCTTTCGAAAGATCCCAACGAAAAGCGTGACCACATGGTCCTTCTTGAGTTTGTAACTGCTGCTGGGATTACACATGGCATGGATGAGCTCTACAAATAA"),
		meltingtemp:         wunit.NewTemperature(52.764, "C"),
		mintemp:             wunit.NewTemperature(50, "C"),
		maxtemp:             wunit.NewTemperature(85, "C"),
		maxGCcontent:        0.6,
		minlength:           25,
		maxlength:           45,
		seqstoavoid:         []string{""},
		overlapthreshold:    45,
		outputoligoseq:      "ATGAGCAAAGGAGAAGAACTTTTCA",
		calculatedGCcontent: 0.36},
}

var overlaptests = []testpair{

	{primer1: "AATACCAGTAGGGTAGAAGAGCACG",
		primer2:        "AATACCAGTAGGGTAGAAGAGCAC",
		overlapPercent: 1.0,
		overlapNumber:  24,
		overlapSeq:     "AATACCAGTAGGGTAGAAGAGCAC"},
	{primer1: "AAAAAAAAAAAAAAAA",
		primer2:        "TTTTTTTTTTTTT",
		overlapPercent: 0.0,
		overlapNumber:  0,
		overlapSeq:     ""},
	{primer1: "AATACCAGTAGGGTAGAAGAGCACG",
		primer2:        "CCAAAAAATAGAAGAGCAC",
		overlapPercent: 0.5789473684210527,
		overlapNumber:  11,
		overlapSeq:     "TAGAAGAGCAC"},
	{primer2: "TAGAAGAGCACGCCCCCCCC",
		primer1:        "CCAAAAAATAGAAGAGCAC",
		overlapPercent: 0.5789473684210527,
		overlapNumber:  11,
		overlapSeq:     "TAGAAGAGCAC"},
	{primer2: "",
		primer1:        "CCAAAAAATAGAAGAGCAC",
		overlapPercent: 0.0,
		overlapNumber:  0,
		overlapSeq:     ""},
}

func TestBasicMeltingTemp(t *testing.T) {
	for _, oligo := range meltingtemptests {
		result := BasicMeltingTemp(oligo.sequence)
		if result.ToString() != oligo.meltingtemp.ToString() {
			t.Error(
				"For", oligo.sequence, "/n",
				"expected", oligo.meltingtemp.ToString(), "\n",
				"got", result.ToString(), "\n",
			)
		}
	}

}

func TestFWDOligoSeq(t *testing.T) {
	for _, oligo := range oligotests {
		oligoseq, err := FWDOligoSeq(oligo.sequence, oligo.maxGCcontent, oligo.minlength, oligo.maxlength, oligo.mintemp, oligo.maxtemp, oligo.seqstoavoid, oligo.overlapthreshold)
		if oligoseq.Sequence() != oligo.outputoligoseq {
			t.Error(
				"For", oligo.sequence, "/n",
				"expected", oligo.outputoligoseq, "\n",
				"got", oligoseq.Sequence(), "\n",
			)
		}
		if oligoseq.GCContent != oligo.calculatedGCcontent {
			t.Error(
				"For", oligo.sequence, "/n",
				"expected", oligo.calculatedGCcontent, "\n",
				"got", oligoseq.GCContent, "\n",
			)
		}
		if err != nil {
			t.Error(
				"errors:", err.Error(), "\n",
			)
		}

	}
}
func TestOverlapCheck(t *testing.T) {
	for _, test := range overlaptests {
		percent, number, seq := OverlapCheck(test.primer1, test.primer2)

		if percent != test.overlapPercent {
			t.Error(
				"For", test.primer1, " and ", test.primer2, "/n",
				"expected", test.overlapPercent, "\n",
				"got", percent, "\n",
			)
		}

		if number != test.overlapNumber {
			t.Error(
				"For", test.primer1, " and ", test.primer2, "/n",
				"expected", test.overlapNumber, "\n",
				"got", number, "\n",
			)
		}

		if seq != test.overlapSeq {
			t.Error(
				"For", test.primer1, " and ", test.primer2, "/n",
				"expected", test.overlapSeq, "\n",
				"got", seq, "\n",
			)
		}
	}
}
