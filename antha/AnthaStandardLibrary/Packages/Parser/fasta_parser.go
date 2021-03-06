// antha/AnthaStandardLibrary/Packages/Parser/fasta_parser.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package parser

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type Fasta struct {
	Id   string
	Desc string
	Seq  string
}

// This will retrieve seq from FASTA file
func RetrieveSeqFromFASTA(id string, filename string) (seq wtype.DNASequence, err error) {

	var nofeatures []wtype.Feature

	allparts, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fastaFh := bytes.NewReader(allparts)

	// then retrieve the particular record
	for record := range FastaParse(fastaFh) {
		if strings.Contains(record.Id, id) {
			seq = wtype.DNASequence{record.Id, record.Seq, true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", nofeatures}
			return
		}
	}

	seq = wtype.DNASequence{"", "", true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", nofeatures} // blank seq
	if seq.Seq == "" {
		err = errors.New("Record not found in file")
		return
	}
	return
}

// This will retrieve seq from FASTA file
func FASTAtoLinearDNASeqs(filename string) (seqs []wtype.DNASequence, err error) {

	var nofeatures []wtype.Feature

	seqs = make([]wtype.DNASequence, 0)

	var seq wtype.DNASequence

	allparts, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fastaFh := bytes.NewReader(allparts)

	// then retrieve the particular record
	for record := range FastaParse(fastaFh) {
		//if strings.Contains(record.Id, id) {
		seq = wtype.DNASequence{record.Id, record.Seq, false, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", nofeatures}

		seqs = append(seqs, seq)

	}
	return

}

// This will retrieve seq from FASTA file
func FASTAtoPlasmidDNASeqs(filename string) (seqs []wtype.DNASequence, err error) {

	var nofeatures []wtype.Feature

	seqs = make([]wtype.DNASequence, 0)

	var seq wtype.DNASequence

	allparts, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fastaFh := bytes.NewReader(allparts)

	// then retrieve the particular record
	for record := range FastaParse(fastaFh) {
		//if strings.Contains(record.Id, id) {
		seq = wtype.DNASequence{record.Id, record.Seq, true, false, wtype.Overhang{0, 0, 0, "", false}, wtype.Overhang{0, 0, 0, "", false}, "", nofeatures}

		seqs = append(seqs, seq)

	}
	return

}

func FastatoDNASequences(inputfilename string) (seqs []wtype.DNASequence, err error) {
	fastaFh, err := os.Open(inputfilename)
	if err != nil {
		return
	}
	defer fastaFh.Close()

	seqs = make([]wtype.DNASequence, 0)

	var seq wtype.DNASequence

	/*records := make([][]string, 0)
	seq := make([]string, 0)
	seq = []string{"#Name", "Sequence", "Plasmid?", "Seq Type", "Class"}
	records = append(records, seq)*/
	for record := range FastaParse(fastaFh) {
		plasmidstatus := ""
		//seqtype := "DNA"
		//class := "not specified"
		if strings.Contains(strings.ToUpper(record.Desc), "PLASMID") || strings.Contains(strings.ToUpper(record.Desc), "CIRCULAR") || strings.Contains(strings.ToUpper(record.Desc), "VECTOR") {
			plasmidstatus = "PLASMID"
		}
		/*	if strings.Contains(record.Desc, "Amino acid") || strings.Contains(record.Id, "aa") {
				seqtype = "AA"
			}

			if strings.Contains(record.Desc, "Class:") {
				uptoclass := strings.Index(record.Desc, "Class:")
				prefix := uptoclass + len("class:")
				class = record.Desc[prefix:]
			}*/
		seq, err = wtype.MakeDNASequence(record.Id, record.Seq, []string{plasmidstatus})
		if err != nil {
			return seqs, err
		}
		seqs = append(seqs, seq)
	}

	return
}

func Build_fasta(header string, seq bytes.Buffer) (Record Fasta) {
	fields := strings.SplitN(header, " ", 2)

	var record Fasta

	if len(fields) > 1 {
		record.Id = fields[0]
		record.Desc = "`" + fields[1] + "`"
	} else {
		record.Id = fields[0]
		record.Desc = ""
	}

	record.Seq = seq.String()

	Record = record

	return Record
}

// new new version

func FastaParse(fastaFh io.Reader) chan Fasta {

	outputChannel := make(chan Fasta)

	scanner := bufio.NewScanner(fastaFh)
	// scanner.Split(bufio.ScanLines)
	header := ""
	var seq bytes.Buffer

	go func() {
		// Loop over the letters in inputString
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			line = strings.Replace(line, "\t", "", -1)
			sublines := strings.Split(line, "\r")
			for _, line := range sublines {

				if len(line) == 0 {
					continue
				}

				// line := scanner.Text()

				if line[0] == '>' {
					// If we stored a previous identifier, get the DNA string and map to the
					// identifier and clear the string
					if header != "" {
						// outputChannel <- build_fasta(header, seq.String())
						outputChannel <- Build_fasta(header, seq)
						header = ""
						seq.Reset()
					}

					// Standard FASTA identifiers look like: ">id desc"
					header = line[1:]
					//fmt.Println(" header: ", header)
					//header = strings.Join(line[1:], " ")

				} else {
					// Append here since multi-line DNA strings are possible
					seq.WriteString(line)
				}
			}

		}

		outputChannel <- Build_fasta(header, seq)

		// Close the output channel, so anything that loops over it
		// will know that it is finished.
		close(outputChannel)
	}()

	return outputChannel
}

/*
// new version
func FastaParse(fastaFh io.Reader) chan Fasta {

	outputChannel := make(chan Fasta)

	scanner := bufio.NewScanner(fastaFh)
	// scanner.Split(bufio.ScanLines)
	header := ""
	var seq bytes.Buffer

	go func() {
		// Loop over the letters in inputString
		for scanner.Scan() {
			fullline := strings.TrimSpace(scanner.Text())
			sublines := strings.Fields(fullline)
			for _, line := range sublines {

				if len(line) == 0 {
					continue
				}

				// line := scanner.Text()

				if line[0] == '>' {
					// If we stored a previous identifier, get the DNA string and map to the
					// identifier and clear the string
					if header != "" {
						// outputChannel <- build_fasta(header, seq.String())
						outputChannel <- Build_fasta(header, seq)
						header = ""
						seq.Reset()
					}

					// Standard FASTA identifiers look like: ">id desc"
					header = line[1:]
				} else {
					// Append here since multi-line DNA strings are possible
					seq.WriteString(line)
				}
			}

		}

		outputChannel <- Build_fasta(header, seq)

		// Close the output channel, so anything that loops over it
		// will know that it is finished.
		close(outputChannel)
	}()

	return outputChannel
}
*/
/*
func FastaParse(fastaFh io.Reader) chan Fasta {

	outputChannel := make(chan Fasta)

	scanner := bufio.NewScanner(fastaFh)
	//scanner.Split(bufio.ScanLines)
	header := ""
	var seq bytes.Buffer

	go func() {
		// Loop over the letters in inputString
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			line = strings.Replace(line, "^M", "", -1)
			if len(line) == 0 {
				continue
			}

			// line := scanner.Text()

			if line[0] == '>' {
				// If we stored a previous identifier, get the DNA string and map to the
				// identifier and clear the string
				if header != "" {
					// outputChannel <- build_fasta(header, seq.String())
					outputChannel <- Build_fasta(header, seq)
					header = ""
					seq.Reset()
				}

				// Standard FASTA identifiers look like: ">id desc"
				header = line[1:]
			} else {
				// Append here since multi-line DNA strings are possible
				seq.WriteString(line)
			}

		}

		outputChannel <- Build_fasta(header, seq)

		// Close the output channel, so anything that loops over it
		// will know that it is finished.
		close(outputChannel)
	}()

	return outputChannel
}
*/
func Fastatocsv(inputfilename string, outputfileprefix string) (csvfile *os.File, err error) {
	fastaFh, err := os.Open(inputfilename)
	if err != nil {
		return
	}
	defer fastaFh.Close()

	//csvfile, err := os.Create(outputfilename)
	csvfile, err = ioutil.TempFile(os.TempDir(), "csv")
	if err != nil {
		return
	}

	records := make([][]string, 0)
	seq := make([]string, 0)
	seq = []string{"#Name", "Sequence", "Plasmid?", "Seq Type", "Class"}
	records = append(records, seq)
	for record := range FastaParse(fastaFh) {
		plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"
		if strings.Contains(record.Desc, "Plasmid") || strings.Contains(record.Id, "Circular") || strings.Contains(record.Id, "Vector") {
			plasmidstatus = "TRUE"
		}
		if strings.Contains(record.Desc, "Amino acid") || strings.Contains(record.Id, "aa") {
			seqtype = "AA"
		}

		if strings.Contains(record.Desc, "Class:") {
			uptoclass := strings.Index(record.Desc, "Class:")
			prefix := uptoclass + len("class:")
			class = record.Desc[prefix:]
		}
		seq = []string{record.Id, record.Seq, plasmidstatus, seqtype, class}
		records = append(records, seq)
	}

	writer := csv.NewWriter(csvfile)
	for _, record := range records {
		err = writer.Write(record)
		if err != nil {
			return
		}
	}

	writer.Flush()
	return
}
