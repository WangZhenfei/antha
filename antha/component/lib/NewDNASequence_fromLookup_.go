package lib

import (
	"fmt"
	//"math"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/igem"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/entrez"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol

// Data which is returned from this protocol

// Physical inputs to this protocol

// Physical outputs from this protocol

func _NewDNASequence_fromLookupRequirements() {

}

// Actions to perform before protocol itself
func _NewDNASequence_fromLookupSetup(_ctx context.Context, _input *NewDNASequence_fromLookupInput) {

}

// Core process of the protocol: steps to be performed for each input
func _NewDNASequence_fromLookupSteps(_ctx context.Context, _input *NewDNASequence_fromLookupInput, _output *NewDNASequence_fromLookupOutput) {

	var err error

	if _input.EntrezID {

		if _input.DNAID {

			_output.DNA, err = entrez.RetrieveSequence(_input.ID, "nucleotide")

		}
	} else if _input.BiobrickID {

		partdetails := igem.LookUp([]string{_input.ID})

		seq := partdetails.Sequence(_input.ID)

		_output.DNA = wtype.MakeLinearDNASequence(_input.ID, seq)

	} //else {Status = fmt.Sprintln("correct conditions not met")}

	orfs := sequences.FindallORFs(_output.DNA.Seq)
	features := sequences.ORFs2Features(orfs)

	_output.DNA = wtype.Annotate(_output.DNA, features)

	_output.Status = fmt.Sprintln(
		text.Print("DNA_Seq: ", _output.DNA),
		text.Print("ORFs: ", _output.DNA.Features),
	)
	_output.Warnings = err
}

// Actions to perform after steps block to analyze data
func _NewDNASequence_fromLookupAnalysis(_ctx context.Context, _input *NewDNASequence_fromLookupInput, _output *NewDNASequence_fromLookupOutput) {

}

func _NewDNASequence_fromLookupValidation(_ctx context.Context, _input *NewDNASequence_fromLookupInput, _output *NewDNASequence_fromLookupOutput) {

}
func _NewDNASequence_fromLookupRun(_ctx context.Context, input *NewDNASequence_fromLookupInput) *NewDNASequence_fromLookupOutput {
	output := &NewDNASequence_fromLookupOutput{}
	_NewDNASequence_fromLookupSetup(_ctx, input)
	_NewDNASequence_fromLookupSteps(_ctx, input, output)
	_NewDNASequence_fromLookupAnalysis(_ctx, input, output)
	_NewDNASequence_fromLookupValidation(_ctx, input, output)
	return output
}

func NewDNASequence_fromLookupRunSteps(_ctx context.Context, input *NewDNASequence_fromLookupInput) *NewDNASequence_fromLookupSOutput {
	soutput := &NewDNASequence_fromLookupSOutput{}
	output := _NewDNASequence_fromLookupRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func NewDNASequence_fromLookupNew() interface{} {
	return &NewDNASequence_fromLookupElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &NewDNASequence_fromLookupInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _NewDNASequence_fromLookupRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &NewDNASequence_fromLookupInput{},
			Out: &NewDNASequence_fromLookupOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type NewDNASequence_fromLookupElement struct {
	inject.CheckedRunner
}

type NewDNASequence_fromLookupInput struct {
	BiobrickID bool
	DNAID      bool
	EntrezID   bool
	ID         string
}

type NewDNASequence_fromLookupOutput struct {
	DNA      wtype.DNASequence
	Status   string
	Warnings error
}

type NewDNASequence_fromLookupSOutput struct {
	Data struct {
		DNA      wtype.DNASequence
		Status   string
		Warnings error
	}
	Outputs struct {
	}
}

func init() {
	addComponent(Component{Name: "NewDNASequence_fromLookup",
		Constructor: NewDNASequence_fromLookupNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/NewDNASequence/NewDNASequence_fromLookup.an",
			Params: []ParamDesc{
				{Name: "BiobrickID", Desc: "", Kind: "Parameters"},
				{Name: "DNAID", Desc: "", Kind: "Parameters"},
				{Name: "EntrezID", Desc: "", Kind: "Parameters"},
				{Name: "ID", Desc: "", Kind: "Parameters"},
				{Name: "DNA", Desc: "", Kind: "Data"},
				{Name: "Status", Desc: "", Kind: "Data"},
				{Name: "Warnings", Desc: "", Kind: "Data"},
			},
		},
	})
}
