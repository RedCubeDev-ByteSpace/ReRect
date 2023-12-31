package boundnodes

import (
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/symbols"
)

// List of conversion types
// ------------------------
type ConversionType string;
const (
    CT_Identity ConversionType = "Identity conversion"
    CT_Implicit ConversionType = "Implicit conversion"
    CT_Explicit ConversionType = "Explicit conversion"
    CT_None     ConversionType = "No conversion"
)

func ClassifyConversion(from *symbols.TypeSymbol, to *symbols.TypeSymbol) ConversionType {
    // identity casts be gaming
    if from.Equal(to) {
        return CT_Identity
    }

    // Anything can be cast to any
    if to.Equal(compunit.GlobalDataTypeRegister["any"]) &&
       !from.Equal(compunit.GlobalDataTypeRegister["void"]) {
        return CT_Implicit
    }

    // any can be cast to anything
    if !to.Equal(compunit.GlobalDataTypeRegister["void"]) &&
        from.Equal(compunit.GlobalDataTypeRegister["any"]){
        return CT_Implicit
    }

    // up and down casts
    if (from.TypeGroup == symbols.INT   && to.TypeGroup == symbols.INT) ||
       (from.TypeGroup == symbols.FLOAT && to.TypeGroup == symbols.FLOAT) {
        
        // allow implicit upcasts
        if to.TypeSize > from.TypeSize {
            return CT_Implicit
        }

        // down casts need to be explicit
        if to.TypeSize < from.TypeSize {
            return CT_Explicit
        }
    }

    // allow explicit int -> float
    if from.TypeGroup == symbols.INT && to.TypeGroup == symbols.FLOAT {
        return CT_Explicit
    }

    // allow explicit float -> int
    if from.TypeGroup == symbols.FLOAT && to.TypeGroup == symbols.INT {
        return CT_Explicit
    }

    // allow anything explicitly to string
    if !from.Equal(compunit.GlobalDataTypeRegister["void"]) &&
        to.Equal(compunit.GlobalDataTypeRegister["string"]) {
        return CT_Explicit
    }

    // allow anything explicitly from string
    if  from.Equal(compunit.GlobalDataTypeRegister["string"]) &&
       !to.Equal(compunit.GlobalDataTypeRegister["void"]) {
        return CT_Explicit
    }

    // allow implicit container -> trait if the container implements the trait
    if from.TypeGroup == symbols.CONT &&
       to.TypeGroup   == symbols.TRT {

        // check if the container implements the trait
        cnt := from.Container
        for _, t := range cnt.Traits {
            if t.TraitType.Equal(to) {
                return CT_Implicit
            }
        }

        // if it doesnt -> death
        return CT_None
    } 

    // allow explicit trait -> container if the container implements the trait
    if from.TypeGroup == symbols.TRT &&
       to.TypeGroup   == symbols.CONT {

        // check if the container implements the trait
        cnt := to.Container
        for _, t := range cnt.Traits {
            if t.TraitType.Equal(from) {
                return CT_Explicit
            }
        }

        // if it doesnt -> death
        return CT_None
    } 

    // otherwise -> dont convert
    return CT_None
}
