// buffercomponents.go

// Package for dealing with manipulation of buffers
package buffers

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Pubchem"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	"fmt"
)

func StockConcentration(nameofmolecule string, massofmoleculeactuallyaddedinG wunit.Mass, diluent string, totalvolumeinL wunit.Volume) (actualconc wunit.Concentration) {

	molecule := pubchem.MakeMolecule(nameofmolecule)

	// in particular, the molecular weight
	molecularweight := molecule.MolecularWeight

	//diluentmolecule := pubchem.MakeMolecule(diluent)

	fmt.Println("SI value of mass:", massofmoleculeactuallyaddedinG.SIValue())

	actualconcfloat := (massofmoleculeactuallyaddedinG.SIValue() * 1000) / (molecularweight * totalvolumeinL.SIValue())

	actualconc = wunit.NewConcentration(actualconcfloat, "M/l")

	return
}

func Dilute(moleculename string, stockconc wunit.Concentration, stockvolume wunit.Volume, diluentname string, diluentvoladded wunit.Volume) (dilutedconc wunit.Concentration) {

	molecule := pubchem.MakeMolecule(moleculename)

	stockMperL := stockconc.MolPerL(molecule.MolecularWeight)

	diluentSI := diluentvoladded.SIValue()

	stockSI := stockvolume.SIValue()

	dilutedconcMperL := stockMperL.SIValue() * stockSI / (stockSI + diluentSI)

	dilutedconc = wunit.NewConcentration(dilutedconcMperL, "M/l")
	fmt.Println(diluentname)
	return
}

/*
func DiluteToTargetConc(moleculename string, stockconc wunit.Concentration, targetconc wunit.Concentration, diluentname string, targetVolume wunit.Volume) (diluentvolume wunit.Volume, stockvolumetouse wunit.Volume) {

	molecule := pubchem.MakeMolecule(moleculename)

	stockMperL := stockconc.MolPerL(molecule.MolecularWeight)

	diluentSI := diluentvoladded.SIValue()

	//stockSI := stockvolume.SIValue()

	dilutedconcMperL := stockMperL.SIValue() * stockSI / (stockSI + diluentSI)

	dilutedconc = wunit.NewConcentration(dilutedconcMperL, "M/l")
	fmt.Println(diluentname)
	return
}
*/
/*
From pubchem...
type Molecule struct {
	Moleculename     string
	MolecularFormula string  `json:"MolecularFormula"`
	MolecularWeight  float64 `json:"MolecularWeight"`
	CID              int     `json:"CID"`
}


type Substance struct {
	Substancename string
	SID           int `json:"SID"`
}

*/
// approximate formula for substance?
// allow empty field?
// have distinct struct?

/*
type Proteincomponent struct {
wtype.LHComponent
pubchem.Molecule
Conc wunit.Concentration
Seq wtype.ProteinSequence
}

type DNAcomponent struct {
wtype.LHComponent
pubchem.Molecule
conc wunit.Concentration
Seq wtype.DNASequence
}

type Substancecomponent struct {
wtype.LHComponent
pubchem.Substance
}

type Moleculecomponent struct {
wtype.LHComponent
Molecule pubchem.Molecule
Conc wunit.Concentration
}

type Organismcomponent struct {
wtype.LHComponent
Organism wtype.Organism

}


type Component struct {
wtype.LHComponent
Conc wunit.Concentration
Molecule *pubchem.Molecule
Sequence *Seq
}
*/
/*
type Buffercomponent struct {
	wtype.LHComponent
	Type       int
	Typestruct interface{}
}

type Buffer struct {
	Components []Buffercomponent
	BufferPH   *PH
}

const (
	Molecule  = iota // e.g. nh4
	Substance        // e.g. yeast extract
	Protein          //e.g.
	DNA              //e.g.
	Organism
)
*/
/*
func (b *Buffercomponent) MolecularWeight() (mw float64, err error) {

	//b.

	if b.Type == 0 {
		mw = b.Typestruct.MolecularWeight
		return mw, nil
	} else if b.Type == 1 {
		mw = b.Typestruct.MolecularWeight()
		err = fmt.Errorf("Only approximate molecular weight possible with Substance component")
		return mw, err
	} else if b.Type == 2 {
		mw = b.Typestruct.MolecularWeight()
		err = fmt.Errorf("Only approximate molecular weight possible with Protein component")
		return mw, err
	} else if b.Type == 3 {
		mw = b.Typestruct.MolecularWeight()
		err = fmt.Errorf("Only approximate molecular weight possible with DNA component")
		return mw, err
	} else if b.Type == 4 {
		mw = b.Typestruct.MolecularWeight()
		err = fmt.Errorf("Only approximate molecular weight possible with organism component")
		return mw, err
	} else {
		err = fmt.Errorf("unkwown component type")
		return err
	}
}
*/
