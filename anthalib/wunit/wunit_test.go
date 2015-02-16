// wunit/wunit_test.go: Part of the Antha language
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
// 1 Royal College St, London NW1 0NH UK

package wunit

import (
	"testing"
	"fmt"
)

func StoreMeasurement(m Measurement){
	// do something
}

func TestBasic(*testing.T){
	ExampleBasic()
}
func TestTwo(*testing.T){
	ExampleTwo()
}
func TestFour(*testing.T){
	ExampleFour()
}
func TestFive(*testing.T){
	ExampleFive()
}

func ExampleBasic(){
	degreeC:=GenericPrefixedUnit{GenericUnit{"DegreeC", "C", 1.0, "C"}, SIPrefix{"m", 1e-03}}	
	TdegreeC:=Temperature{ConcreteMeasurement{NewWFloat(1.0), &degreeC}}
	fmt.Println(TdegreeC.SIValue())
	// Output:
	// 0.001
}
func ExampleTwo(){
	Joule:=GenericPrefixedUnit{GenericUnit{"Joule", "J", 1.0, "J"}, SIPrefix{"k", 1e3}}	
	NJoule:=Energy{ConcreteMeasurement{NewWFloat(23.4), &Joule}}
	fmt.Println(NJoule.SIValue())
	// Output:
	// 23400
}

func ExampleFour(){
	p:=NewPressure(56.2, "Pa")
	fmt.Println(p.RawValue())

	p.SetValue(34.0)

	fmt.Println(p.RawValue())

	// Output:
	// 56.2
	// 34
}

func ExampleFive(){

	fmt.Println(PrefixMul("m", "m"))

	// Output: 
	// u
}

func ExampleSix(){
	fmt.Println(k)
	fmt.Println(G)
	fmt.Println(p)
	// Output: 
	// 3
	// 9
	// -12

}
