// antha/AnthaStandardLibrary/Packages/enzymes/exporttofile.go: Part of the Antha language
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

package export

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// function to export a standard report of sequence properties to a txt file
func Exporttofile(dir string, seq wtype.BioSequence) string {
	anthapath.CreatedotAnthafolder()

	filename := fmt.Sprintf("%s%c%s_%s.txt", anthapath.Dirpath(), os.PathSeparator, dir, seq.Name())

	//f, _ := os.Create(filepath.Join(anthapath.Dirpath(), "iGem_registry.txt"))

	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}

	// GC content
	GC := GCcontent(seq.Sequence())

	// Find all orfs:
	orfs := DoublestrandedORFS(seq.Sequence())

	fmt.Fprintln(f, ">", dir[2:]+"_"+seq.Name())
	fmt.Fprintln(f, seq.Sequence())

	fmt.Fprintln(f, "Sequence length:", len(seq.Sequence()))
	fmt.Fprintln(f, "Molecular weight:", wutil.RoundInt(MassDNA(seq.Sequence(), false, true)), "g/mol")
	fmt.Fprintln(f, "GC Content:", wutil.RoundInt((GC * 100)), "%")

	fmt.Fprintln(f, (len(orfs.TopstrandORFS) + len(orfs.BottomstrandORFS)), "Potential Open reading frames found:")
	//fmt.Fprintln(f, "Top strand")
	for _, strandorf := range orfs.TopstrandORFS {
		fmt.Fprintln(f, "Topstrand")
		fmt.Fprintln(f, "Position:", strandorf.StartPosition, "..", strandorf.EndPosition)

		fmt.Fprintln(f, " DNA Sequence:", strandorf.DNASeq)

		fmt.Fprintln(f, "Translated Amino Acid Sequence:", strandorf.ProtSeq)
		fmt.Fprintln(f, "Length of Amino acid sequence:", len(strandorf.ProtSeq)-1)
		fmt.Fprintln(f, "molecular weight:", Molecularweight(strandorf), "kDA")
	}
	//fmt.Fprintln(f, "Bottom strand")
	for _, strandorf := range orfs.BottomstrandORFS {
		fmt.Fprintln(f, "Bottom strand")
		fmt.Fprintln(f, "Position:", strandorf.StartPosition, "..", strandorf.EndPosition)

		fmt.Fprintln(f, " DNA Sequence:", strandorf.DNASeq)

		fmt.Fprintln(f, "Translated Amino Acid Sequence:", strandorf.ProtSeq)
		fmt.Fprintln(f, "Length of Amino acid sequence:", len(strandorf.ProtSeq)-1)
		fmt.Fprintln(f, "molecular weight:", Molecularweight(strandorf), "kDA")
	}
	f.Close()

	return filename
}

// function to export a sequence to a txt file
func ExportFasta(dir string, seq wtype.BioSequence) string {
	anthapath.CreatedotAnthafolder()

	filename := fmt.Sprintf("%s%c%s_%s.fasta", anthapath.Dirpath(), os.PathSeparator, dir, seq.Name())
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())

	f.Close()

	return filename
}

// function to export a sequence to a txt file
func ExportFastaDir(dir string, file string, seq wtype.BioSequence) string {
	filename := fmt.Sprintf("%s%c%s_%s.fasta", anthapath.Dirpath(), os.PathSeparator, dir, seq.Name())
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())

	f.Close()

	return filename
}

func ExportReport(dir string, seq wtype.BioSequence) string {
	anthapath.CreatedotAnthafolder()

	filename := fmt.Sprintf("%s%c%s_%s.txt", anthapath.Dirpath(), os.PathSeparator, dir, seq.Name())
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())

	f.Close()

	return filename
}

// function to export multiple sequences in fasta format into a single txt file
// Modify this for the more general case
func Makefastaserial(dir string, seqs []*wtype.DNASequence) string {
	anthapath.CreatedotAnthafolder()
	filename := fmt.Sprintf("%s%c%s.fasta", anthapath.Dirpath(), os.PathSeparator, dir)
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}

	for _, seq := range seqs {
		fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())
	}

	f.Close()
	return filename
}

func Makefastaserial2(dir string, seqs []wtype.DNASequence) string {
	anthapath.CreatedotAnthafolder()
	filename := fmt.Sprintf("%s%c%s.fasta", anthapath.Dirpath(), os.PathSeparator, dir)
	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}

	for _, seq := range seqs {
		fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())
	}

	f.Close()
	return filename
}

func ExportFastaandSummaryforEachSeq(assemblyparameters enzymes.Assemblyparameters) (err error) {

	enzymename := strings.ToUpper(assemblyparameters.Enzymename)

	// should change this to rebase lookup; what happens if this fails?
	//enzyme := TypeIIsEnzymeproperties[enzymename]
	enzyme, err := lookup.TypeIIsLookup(enzymename)
	if err != nil {
		return err
	}
	//assemble (note that sapIenz is found in package enzymes)
	_, plasmidproductsfromXprimaryseq, err := enzymes.JoinXnumberofparts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)

	if err != nil {
		return err
	}

	for _, assemblyproduct := range plasmidproductsfromXprimaryseq {

		fileprefix := anthapath.Dirpath() + "/"
		tojoin := make([]string, 0)
		tojoin = append(tojoin, fileprefix, assemblyparameters.Constructname)
		filename := strings.Join(tojoin, "")
		Exporttofile(filename, &assemblyproduct)
		ExportFasta(filename, &assemblyproduct)
	}
	return err
}

func ExportFastaSerialfromMultipleAssemblies(dirname string, multipleassemblyparameters []enzymes.Assemblyparameters) (filename string, err error) {

	seqs := make([]wtype.DNASequence, 0)

	for _, assemblyparameters := range multipleassemblyparameters {

		enzymename := strings.ToUpper(assemblyparameters.Enzymename)

		// should change this to rebase lookup; what happens if this fails?
		//enzyme := TypeIIsEnzymeproperties[enzymename]
		enzyme, err := lookup.TypeIIsLookup(enzymename)
		if err != nil {
			return filename, err
		}
		//assemble (note that sapIenz is found in package enzymes)
		_, plasmidproductsfromXprimaryseq, err := enzymes.JoinXnumberofparts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)

		if err != nil {
			return filename, err
		}

		for _, assemblyproduct := range plasmidproductsfromXprimaryseq {

			/*	fileprefix := anthapath.Dirpath() + "/"
				tojoin := make([]string, 0)
				tojoin = append(tojoin, fileprefix, assemblyparameters.Constructname)
				filename := strings.Join(tojoin, "")
				Exporttofile(filename, &assemblyproduct)
				ExportFasta(filename, &assemblyproduct)*/

			seqs = append(seqs, assemblyproduct)
		}

	}

	filename = Makefastaserial2(dirname, seqs)

	return filename, err
}

func ExporttoTextFile(filename string, data []string) (err error) {

	//filename := fmt.Sprintf("%s%c%s_%s.csv", anthapath.Dirpath(), os.PathSeparator, dir, name)

	//f, _ := os.Create(filepath.Join(anthapath.Dirpath(), "iGem_registry.txt"))

	f, e := os.Create(filename)
	if e != nil {
		log.Fatal(e)
	}

	for _, str := range data {

		fmt.Fprintln(f, str)

		/*	_, err = f.WriteString(str)

			if err != nil {
				fmt.Println(err.Error())
				return
			}
		*/
	}
	f.Close()

	return
}
