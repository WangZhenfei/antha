package lib

import (
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

func _Synthesis_validationRequirements() {
}

func _Synthesis_validationSetup(_ctx context.Context, _input *Synthesis_validationInput) {
}

func _Synthesis_validationSteps(_ctx context.Context, _input *Synthesis_validationInput, _output *Synthesis_validationOutput) {

	// check if sequence meets requirements for synthesis
	_output.ValiadationStatus, _output.Validated = sequences.ValidateSynthesis(_input.PartsWithOverhangs, _input.Vector.Name(), "GenScript")

	if _output.Validated {
		_output.ValidatedPartsWithOverhangs = _input.PartsWithOverhangs
	}

}

func _Synthesis_validationAnalysis(_ctx context.Context, _input *Synthesis_validationInput, _output *Synthesis_validationOutput) {

}

func _Synthesis_validationValidation(_ctx context.Context, _input *Synthesis_validationInput, _output *Synthesis_validationOutput) {

}
func _Synthesis_validationRun(_ctx context.Context, input *Synthesis_validationInput) *Synthesis_validationOutput {
	output := &Synthesis_validationOutput{}
	_Synthesis_validationSetup(_ctx, input)
	_Synthesis_validationSteps(_ctx, input, output)
	_Synthesis_validationAnalysis(_ctx, input, output)
	_Synthesis_validationValidation(_ctx, input, output)
	return output
}

func Synthesis_validationRunSteps(_ctx context.Context, input *Synthesis_validationInput) *Synthesis_validationSOutput {
	soutput := &Synthesis_validationSOutput{}
	output := _Synthesis_validationRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func Synthesis_validationNew() interface{} {
	return &Synthesis_validationElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &Synthesis_validationInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _Synthesis_validationRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &Synthesis_validationInput{},
			Out: &Synthesis_validationOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type Synthesis_validationElement struct {
	inject.CheckedRunner
}

type Synthesis_validationInput struct {
	PartsWithOverhangs []wtype.DNASequence
	SynthesisProvider  string
	Vector             wtype.DNASequence
}

type Synthesis_validationOutput struct {
	ValiadationStatus           string
	Validated                   bool
	ValidatedPartsWithOverhangs []wtype.DNASequence
}

type Synthesis_validationSOutput struct {
	Data struct {
		ValiadationStatus           string
		Validated                   bool
		ValidatedPartsWithOverhangs []wtype.DNASequence
	}
	Outputs struct {
	}
}

func init() {
	addComponent(Component{Name: "Synthesis_validation",
		Constructor: Synthesis_validationNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Data/DNA/GeneDesign/SynthesisValidation.an",
			Params: []ParamDesc{
				{Name: "PartsWithOverhangs", Desc: "", Kind: "Parameters"},
				{Name: "SynthesisProvider", Desc: "", Kind: "Parameters"},
				{Name: "Vector", Desc: "", Kind: "Parameters"},
				{Name: "ValiadationStatus", Desc: "", Kind: "Data"},
				{Name: "Validated", Desc: "", Kind: "Data"},
				{Name: "ValidatedPartsWithOverhangs", Desc: "", Kind: "Data"},
			},
		},
	})
}
