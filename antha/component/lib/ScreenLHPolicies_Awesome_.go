package lib

import (
	"fmt"
	antha "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"path/filepath"
	"strconv"
)

// Input parameters for this protocol (data)

// Data which is returned from this protocol, and data types

//[]string //map[string]string

//NeatSamplewells []string

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _ScreenLHPolicies_AwesomeRequirements() {
}

// Conditions to run on startup
func _ScreenLHPolicies_AwesomeSetup(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _ScreenLHPolicies_AwesomeSteps(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput, _output *ScreenLHPolicies_AwesomeOutput) {

	// validate presence of doe design file in anthapath
	if antha.Anthafileexists(_input.LHDOEFile) == false {
		fmt.Println("This DOE file ", _input.LHDOEFile, " was not found in anthapath ~.antha. Please move it there, change file name and type in antha-lang/antha/microarch/driver/makelhpolicy.go and recompile antha to use this liquidhandling doe design")
		fmt.Println("currently set to ", liquidhandling.DOEliquidhandlingFile, " type ", liquidhandling.DXORJMP)
	} else {
		fmt.Println("found lhpolicy doe file", _input.LHDOEFile)
	}

	// declare some global variables for use later
	var rotate = false
	var autorotate = true
	var wellpositionarray = make([]string, 0)
	var perconditionuntowelllocationmap = make([]string, 0)
	var alphabet = wutil.MakeAlphabetArray()
	_output.Runtowelllocationmap = make(map[string]string)
	_output.Blankwells = make([]string, 0)
	counter := 0
	var platenum = 1
	// work out plate layout based on picture or just in order

	if _input.Printasimage {
		chosencolourpalette := image.AvailablePalettes["Palette1"]
		positiontocolourmap, _, _ := image.ImagetoPlatelayout(_input.Imagefilename, _input.OutPlate, &chosencolourpalette, rotate, autorotate)

		//Runtowelllocationmap = make([]string,0)

		for location, colour := range positiontocolourmap {
			R, G, B, A := colour.RGBA()

			if uint8(R) == 242 && uint8(G) == 243 && uint8(B) == 242 && uint8(A) == 255 {
				continue
			} else {
				wellpositionarray = append(wellpositionarray, location)
			}
		}

	} else {

		for j := 0; j < _input.OutPlate.WlsX; j++ {
			for i := 0; i < _input.OutPlate.WlsY; i++ { //countingfrom1iswhatmakesushuman := j + 1
				//k = k + 1
				wellposition := string(alphabet[i]) + strconv.Itoa(j+1)
				//fmt.Println(wellposition, k)
				wellpositionarray = append(wellpositionarray, wellposition)
			}

		}
	}
	reactions := make([]*wtype.LHComponent, 0)

	//policies, names := liquidhandling.PolicyMaker(liquidhandling.Allpairs, "DOE_run",false)

	//intfactors := []string{"Pre_MIX","POST_MIX"}
	_, names, runs, err := liquidhandling.PolicyMakerfromDesign(_input.DXORJMP, _input.LHDOEFile, "DOE_run")
	if err != nil {
		panic(err)
	}

	//newruns := make([]doe.Run,len(runs))

	for l := 0; l < len(_input.TestSolVolumes); l++ {
		for k := 0; k < len(_input.TestSols); k++ {
			for j := 0; j < _input.NumberofReplicates; j++ {
				for i := 0; i < len(runs); i++ {

					if counter == ((_input.OutPlate.WlsX * _input.OutPlate.WlsY) + _input.NumberofBlanks) {
						fmt.Println("plate full, counter = ", counter)
						platenum++
						counter = 0
					}

					//eachreaction := make([]*wtype.LHComponent, 0)

					// keep default policy for diluent
					//Diluent.Type = names[i]
					//fmt.Println(Diluent.Type)

					// diluent first
					bufferSample := mixer.SampleForTotalVolume(_input.Diluent, _input.TotalVolume)
					//eachreaction = append(eachreaction,bufferSample)

					solution := execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, bufferSample)

					// now test sample

					// change liquid class
					_input.TestSols[k].Type = wtype.LiquidTypeFromString(names[i])

					//sample
					testSample := mixer.Sample(_input.TestSols[k], _input.TestSolVolumes[l])

					//eachreaction = append(eachreaction,testSample)

					// pipette out
					solution = execute.MixTo(_ctx, _input.OutPlate.Type, wellpositionarray[counter], platenum, testSample)

					perconditionuntowelllocationmap = append(perconditionuntowelllocationmap, wtype.LiquidTypeName(_input.TestSols[k].Type)+":"+wellpositionarray[counter])

					// get annotation info
					doerun := wtype.LiquidTypeName(_input.TestSols[k].Type)

					volume := _input.TestSolVolumes[l].ToString() //strconv.Itoa(wutil.RoundInt(number))+"ul"

					solutionname := _input.TestSols[k].CName

					description := volume + "_" + solutionname + "_replicate" + strconv.Itoa(j+1) + "_platenum" + strconv.Itoa(platenum)
					//setpoints := volume+"_"+solutionname+"_replicate"+strconv.Itoa(j+1)+"_platenum"+strconv.Itoa(platenum)

					// add run to well position lookup table
					_output.Runtowelllocationmap[doerun+"_"+description] = wellpositionarray[counter]
					reactions = append(reactions, solution)
					counter = counter + 1

					// add additional info for each run
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "Location_"+description, wellpositionarray[counter])

					// add run order:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "runorder_"+description, counter)

					// add setpoint printout to double check correct match up:
					runs[i] = doe.AddAdditionalHeaderandValue(runs[i], "Additional", "doerun"+description, doerun)
					//runs[i].AddAdditionalValue("Replicate", strconv.Itoa(j+1))
					//runs[i].AddAdditionalValue("Solution name", TestSols[k].CName)
					//runs[i].AddAdditionalValue("Volume", strconv.Itoa(wutil.RoundInt(TestSolVolumes[l].RawValue()))+"ul)

				}

				// export DOE design file per set of conditions
				outputsandwich := strconv.Itoa(wutil.RoundInt(_input.TestSolVolumes[l].RawValue())) + _input.TestSols[k].CName + strconv.Itoa(j+1)

				outputfilename := filepath.Join(antha.Dirpath(), "DOE2"+"_"+outputsandwich+".xlsx")

				_output.Errors = append(_output.Errors, doe.AddWelllocations(_input.DXORJMP, filepath.Join(antha.Dirpath(), _input.LHDOEFile), 0, perconditionuntowelllocationmap, "DOE_run", outputfilename, []string{"Volume", "Solution", "Replicate"}, []interface{}{_input.TestSolVolumes[l].ToString(), _input.TestSols[k].CName, string(j)}))

				// other things to add to check for covariance
				// order in which wells were pippetted
				// plate ID
				// row
				// column
				// ambient temp

				// empty
				perconditionuntowelllocationmap = make([]string, 0)
			}
		}
	}

	// export overall DOE design file showing all well locations for all conditions
	_ = doe.JMPXLSXFilefromRuns(runs, _input.OutputFilename)

	// add blanks after

	for n := 0; n < platenum; n++ {
		for m := 0; m < _input.NumberofBlanks; m++ {
			//eachreaction := make([]*wtype.LHComponent, 0)

			// use defualt policy for blank

			bufferSample := mixer.Sample(_input.Diluent, _input.TotalVolume)
			//eachreaction = append(eachreaction,bufferSample)

			// add blanks to last column of plate
			well := alphabet[_input.OutPlate.WlsY-1-m] + strconv.Itoa(_input.OutPlate.WlsX)
			fmt.Println("blankwell", well)
			reaction := execute.MixTo(_ctx, _input.OutPlate.Type, well, n+1, bufferSample)
			//fmt.Println("where am I?",wellpositionarray[counter])
			_output.Runtowelllocationmap["Blank"+strconv.Itoa(m+1)+" platenum"+strconv.Itoa(n+1)] = well

			_output.Blankwells = append(_output.Blankwells, well)

			reactions = append(reactions, reaction)
			counter = counter + 1

		}

	}

	_output.Reactions = reactions
	_output.Runcount = len(_output.Reactions)
	_output.Pixelcount = len(wellpositionarray)
	_output.Runs = runs

}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _ScreenLHPolicies_AwesomeAnalysis(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput, _output *ScreenLHPolicies_AwesomeOutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _ScreenLHPolicies_AwesomeValidation(_ctx context.Context, _input *ScreenLHPolicies_AwesomeInput, _output *ScreenLHPolicies_AwesomeOutput) {
}
func _ScreenLHPolicies_AwesomeRun(_ctx context.Context, input *ScreenLHPolicies_AwesomeInput) *ScreenLHPolicies_AwesomeOutput {
	output := &ScreenLHPolicies_AwesomeOutput{}
	_ScreenLHPolicies_AwesomeSetup(_ctx, input)
	_ScreenLHPolicies_AwesomeSteps(_ctx, input, output)
	_ScreenLHPolicies_AwesomeAnalysis(_ctx, input, output)
	_ScreenLHPolicies_AwesomeValidation(_ctx, input, output)
	return output
}

func ScreenLHPolicies_AwesomeRunSteps(_ctx context.Context, input *ScreenLHPolicies_AwesomeInput) *ScreenLHPolicies_AwesomeSOutput {
	soutput := &ScreenLHPolicies_AwesomeSOutput{}
	output := _ScreenLHPolicies_AwesomeRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func ScreenLHPolicies_AwesomeNew() interface{} {
	return &ScreenLHPolicies_AwesomeElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &ScreenLHPolicies_AwesomeInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _ScreenLHPolicies_AwesomeRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &ScreenLHPolicies_AwesomeInput{},
			Out: &ScreenLHPolicies_AwesomeOutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type ScreenLHPolicies_AwesomeElement struct {
	inject.CheckedRunner
}

type ScreenLHPolicies_AwesomeInput struct {
	DXORJMP            string
	Diluent            *wtype.LHComponent
	Imagefilename      string
	LHDOEFile          string
	NumberofBlanks     int
	NumberofReplicates int
	OutPlate           *wtype.LHPlate
	OutputFilename     string
	Printasimage       bool
	TestSolVolumes     []wunit.Volume
	TestSols           []*wtype.LHComponent
	TotalVolume        wunit.Volume
}

type ScreenLHPolicies_AwesomeOutput struct {
	Blankwells           []string
	Errors               []error
	Pixelcount           int
	Reactions            []*wtype.LHComponent
	Runcount             int
	Runs                 []doe.Run
	Runtowelllocationmap map[string]string
}

type ScreenLHPolicies_AwesomeSOutput struct {
	Data struct {
		Blankwells           []string
		Errors               []error
		Pixelcount           int
		Runcount             int
		Runs                 []doe.Run
		Runtowelllocationmap map[string]string
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	addComponent(Component{Name: "ScreenLHPolicies_Awesome",
		Constructor: ScreenLHPolicies_AwesomeNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/FindbestLHPolicy/ScreenLHPolicies_Awesome.an",
			Params: []ParamDesc{
				{Name: "DXORJMP", Desc: "", Kind: "Parameters"},
				{Name: "Diluent", Desc: "", Kind: "Inputs"},
				{Name: "Imagefilename", Desc: "", Kind: "Parameters"},
				{Name: "LHDOEFile", Desc: "", Kind: "Parameters"},
				{Name: "NumberofBlanks", Desc: "", Kind: "Parameters"},
				{Name: "NumberofReplicates", Desc: "", Kind: "Parameters"},
				{Name: "OutPlate", Desc: "", Kind: "Inputs"},
				{Name: "OutputFilename", Desc: "", Kind: "Parameters"},
				{Name: "Printasimage", Desc: "", Kind: "Parameters"},
				{Name: "TestSolVolumes", Desc: "", Kind: "Parameters"},
				{Name: "TestSols", Desc: "", Kind: "Inputs"},
				{Name: "TotalVolume", Desc: "", Kind: "Parameters"},
				{Name: "Blankwells", Desc: "", Kind: "Data"},
				{Name: "Errors", Desc: "", Kind: "Data"},
				{Name: "Pixelcount", Desc: "", Kind: "Data"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
				{Name: "Runcount", Desc: "", Kind: "Data"},
				{Name: "Runs", Desc: "", Kind: "Data"},
				{Name: "Runtowelllocationmap", Desc: "[]string //map[string]string\n", Kind: "Data"},
			},
		},
	})
}
