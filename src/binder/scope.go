package binder

import (
	"bytespace.network/rerect/symbols"
)

// Scope struct
// ------------
type Scope struct {
    Parent *Scope

    Variables []symbols.VariableSymbol
}

// Constructor
func NewScope(parent *Scope) *Scope {
    return &Scope{
        Parent: parent,
        Variables: make([]symbols.VariableSymbol, 0),
    }
}

func (scp *Scope) RegisterVariable(vari symbols.VariableSymbol) bool {
    // make sure this variable name isnt already in use
    for _, v := range scp.Variables {
        if v.Name() == vari.Name() {
            return false
        }
    }

    // if it isnt -> register this variable
    scp.Variables = append(scp.Variables, vari)
    return true
}

func (scp *Scope) LookupVariable(name string) symbols.VariableSymbol {
    // Look for this variable locally
    for _, v := range scp.Variables {
        if v.Name() == name {
            return v
        }
    }

    // if it wasnt found -> do we have a parent scope?
    if scp.Parent != nil {
        // do a lookup on the parent
        return scp.Parent.LookupVariable(name)

    // otherwise: no fucking clue
    } else {
        return nil
    }
}
