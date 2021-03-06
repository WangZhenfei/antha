// liquidhandling/setupagent.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/logger"
)

// v2.0 should be another linear program - basically just want to optimize
// positioning in the face of constraints

// default setup agent
func BasicSetupAgent(request *LHRequest, params *liquidhandling.LHProperties) (*LHRequest, error) {
	// this is quite tricky and requires extensive interaction with the liquid handling
	// parameters

	// the principal question is how to define constraints on the system

	// I think this needs to remain tbd for now
	// instead we can rely on the preference system I already use

	plate_lookup := make(map[string]string, 5)
	//tip_lookup := make([]*wtype.LHTipbox, 0, 5)

	//	tip_preferences := params.Tip_preferences
	input_preferences := params.Input_preferences
	output_preferences := params.Output_preferences

	// how do we set the below?
	// we don't know how many tips we need until we generate
	// instructions; ditto input or output plates until we've done layout

	// input plates
	input_plates := request.Input_plates

	// output plates
	output_plates := request.Output_plates

	// tips
	tips := request.Tips

	// just need to set the tip types
	// these should be distinct... we should check really
	// ...eventually
	if len(tips) != 0 {
		for _, tb := range tips {
			if tb == nil {
				continue
			}
			params.Tips = append(params.Tips, tb.Tips[0][0])
		}
	}

	setup := make(map[string]interface{})
	// make sure anything in setup is in synch

	for pos, id := range params.PosLookup {
		if id != "" {
			p := params.PlateLookup[id]
			setup[pos] = p
		}

	}

	// this logic may not transfer well but I expect that outputs are more constrained
	// than inputs for the simple reason that most output takes place to single wells
	// while input (sometimes) takes place from reservoirs

	// outputs

	for _, p := range output_plates {
		position := get_first_available_preference(output_preferences, setup)
		if position == "" {
			//RaiseError("No positions left for output")
			err := wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "No positions left for output")
			return request, err
		}
		setup[position] = p
		plate_lookup[p.ID] = position
		params.AddPlate(position, p)
		logger.Info(fmt.Sprintf("Output plate of type %s in position %s", p.Type, position))
	}

	for _, p := range input_plates {
		position := get_first_available_preference(input_preferences, setup)
		if position == "" {
			//RaiseError("No positions left for input")
			err := wtype.LHError(wtype.LH_ERR_NO_DECK_SPACE, "No positions left for input")
			return request, err
		}
		//fmt.Println("PLAATE: ", position)
		setup[position] = p
		plate_lookup[p.ID] = position
		params.AddPlate(position, p)
		logger.Info(fmt.Sprintf("Input plate of type %s in position %s", p.Type, position))
	}

	//request.Setup = setup
	request.Plate_lookup = plate_lookup
	return request, nil
}

func get_first_available_preference(prefs []string, setup map[string]interface{}) string {
	for _, pref := range prefs {
		_, ok := setup[pref]
		if !ok {
			return pref
		}

	}
	return ""
}
