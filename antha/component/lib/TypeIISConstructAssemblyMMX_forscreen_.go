package lib

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
)

// Input parameters for this protocol (data)

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

// Data which is returned from this protocol, and data types

func _TypeIISConstructAssemblyMMX_forscreenRequirements() {}

// Conditions to run on startup
func _TypeIISConstructAssemblyMMX_forscreenSetup(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreenInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _TypeIISConstructAssemblyMMX_forscreenSteps(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreenInput, _output *TypeIISConstructAssemblyMMX_forscreenOutput) {
	samples := make([]*wtype.LHComponent, 0)

	waterSample := mixer.SampleForTotalVolume(_input.Water, _input.ReactionVolume)
	samples = append(samples, waterSample)

	mmxSample := mixer.Sample(_input.MasterMix, _input.MasterMixVolume)
	samples = append(samples, mmxSample)

	for k, part := range _input.Parts {
		fmt.Println("creating dna part num ", k, " comp ", part.CName, " renamed to ", _input.PartNames[k], " vol ", _input.PartVols[k])

		part.Type = wtype.LiquidTypeFromString(_input.LHPolicyName)

		partSample := mixer.Sample(part, _input.PartVols[k])
		partSample.CName = _input.PartNames[k]
		samples = append(samples, partSample)
	}

	// ensure the last step is mixed
	samples[len(samples)-1].Type = wtype.LTDNAMIX
	_output.Reaction = execute.MixTo(_ctx, _input.OutPlate.Type, _input.OutputLocation, _input.OutputPlateNum, samples...)

	// incubate the reaction mixture
	// commented out pending changes to incubate
	execute.Incubate(_ctx, _output.Reaction, _input.ReactionTemp, _input.ReactionTime, false)
	// inactivate
	//Incubate(Reaction, InactivationTemp, InactivationTime, false)
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _TypeIISConstructAssemblyMMX_forscreenAnalysis(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreenInput, _output *TypeIISConstructAssemblyMMX_forscreenOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _TypeIISConstructAssemblyMMX_forscreenValidation(_ctx context.Context, _input *TypeIISConstructAssemblyMMX_forscreenInput, _output *TypeIISConstructAssemblyMMX_forscreenOutput) {
}
func _TypeIISConstructAssemblyMMX_forscreenRun(_ctx context.Context, input *TypeIISConstructAssemblyMMX_forscreenInput) *TypeIISConstructAssemblyMMX_forscreenOutput {
	output := &TypeIISConstructAssemblyMMX_forscreenOutput{}
	_TypeIISConstructAssemblyMMX_forscreenSetup(_ctx, input)
	_TypeIISConstructAssemblyMMX_forscreenSteps(_ctx, input, output)
	_TypeIISConstructAssemblyMMX_forscreenAnalysis(_ctx, input, output)
	_TypeIISConstructAssemblyMMX_forscreenValidation(_ctx, input, output)
	return output
}

func TypeIISConstructAssemblyMMX_forscreenRunSteps(_ctx context.Context, input *TypeIISConstructAssemblyMMX_forscreenInput) *TypeIISConstructAssemblyMMX_forscreenSOutput {
	soutput := &TypeIISConstructAssemblyMMX_forscreenSOutput{}
	output := _TypeIISConstructAssemblyMMX_forscreenRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func TypeIISConstructAssemblyMMX_forscreenNew() interface{} {
	return &TypeIISConstructAssemblyMMX_forscreenElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &TypeIISConstructAssemblyMMX_forscreenInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _TypeIISConstructAssemblyMMX_forscreenRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &TypeIISConstructAssemblyMMX_forscreenInput{},
			Out: &TypeIISConstructAssemblyMMX_forscreenOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type TypeIISConstructAssemblyMMX_forscreenElement struct {
	inject.CheckedRunner
}

type TypeIISConstructAssemblyMMX_forscreenInput struct {
	InactivationTemp   wunit.Temperature
	InactivationTime   wunit.Time
	LHPolicyName       string
	MasterMix          *wtype.LHComponent
	MasterMixVolume    wunit.Volume
	OutPlate           *wtype.LHPlate
	OutputLocation     string
	OutputPlateNum     int
	OutputReactionName string
	PartNames          []string
	PartVols           []wunit.Volume
	Parts              []*wtype.LHComponent
	ReactionTemp       wunit.Temperature
	ReactionTime       wunit.Time
	ReactionVolume     wunit.Volume
	Water              *wtype.LHComponent
}

type TypeIISConstructAssemblyMMX_forscreenOutput struct {
	Reaction *wtype.LHComponent
}

type TypeIISConstructAssemblyMMX_forscreenSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reaction *wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "TypeIISConstructAssemblyMMX_forscreen",
		Constructor: TypeIISConstructAssemblyMMX_forscreenNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/TypeIIsAssembly/TypeIISConstructAssemblyMMX_forscreen/TypeIISConstructAssemblyMMX.an",
			Params: []ParamDesc{
				{Name: "InactivationTemp", Desc: "", Kind: "Parameters"},
				{Name: "InactivationTime", Desc: "", Kind: "Parameters"},
				{Name: "LHPolicyName", Desc: "", Kind: "Parameters"},
				{Name: "MasterMix", Desc: "", Kind: "Inputs"},
				{Name: "MasterMixVolume", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputLocation", Desc: "", Kind: "Parameters"},
				{Name: "OutputPlateNum", Desc: "", Kind: "Parameters"},
				{Name: "OutputReactionName", Desc: "", Kind: "Parameters"},
				{Name: "PartNames", Desc: "", Kind: "Parameters"},
				{Name: "PartVols", Desc: "", Kind: "Parameters"},
				{Name: "Parts", Desc: "", Kind: "Inputs"},
				{Name: "ReactionTemp", Desc: "", Kind: "Parameters"},
				{Name: "ReactionTime", Desc: "", Kind: "Parameters"},
				{Name: "ReactionVolume", Desc: "", Kind: "Parameters"},
				{Name: "Water", Desc: "", Kind: "Inputs"},
				{Name: "Reaction", Desc: "", Kind: "Outputs"},
			},
		},
	})
}
