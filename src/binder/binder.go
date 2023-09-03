// Binder - binder.go
// --------------------------------------------------------
// The binder is a crucial part of the compilation process
// Its job is to do a semantic analysis on the parsed tree
// --------------------------------------------------------
package binder

import (
	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/span"
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// --------------------------------------------------------
// Function indexing
// --------------------------------------------------------
func IndexFunctions(pck *symbols.PackageSymbol, mem []syntaxnodes.MemberNode) ([]*symbols.FunctionSymbol, []syntaxnodes.StatementNode) {
    syms := []*symbols.FunctionSymbol{}
    bodies := []syntaxnodes.StatementNode{}

    // look through all member nodes
    for _, v := range mem {
        // we're only looking for function nodes
        if v.Type() != syntaxnodes.NT_Function {
            continue
        }

        fncMem := v.(*syntaxnodes.FunctionNode)

        // create parameter symbols
        prms := []*symbols.ParameterSymbol{}
        for i, prm := range fncMem.Parameters {
            prms = append(prms, symbols.NewParameterSymbol(
                prm.ParameterName.Buffer,
                i,
                LookupTypeClause(prm.ParameterType),
            ))
        }

        // register a function symbol for this function
        fnc := symbols.NewFunctionSymbol(
            fncMem.FunctionName.Buffer,
            LookupTypeClause(fncMem.ReturnType),
            prms,
        )

        ok := pck.TryRegisterFunction(fnc) 

        if !ok {
            error.Report(error.NewError(error.BND, fncMem.FunctionName.Position, "Cannot register function '%s'! A function with that name already exists!", fnc.FuncName))
            continue
        }

        syms = append(syms, fnc)
        bodies = append(bodies, fncMem.Body)
    }

    return syms, bodies
}

// --------------------------------------------------------
// Binding
// --------------------------------------------------------
type Binder struct {
    CurrentPackage *symbols.PackageSymbol
    CurrentFunction *symbols.FunctionSymbol
    CurrentScope *Scope
}

func BindFunctions(pck *symbols.PackageSymbol, syms []*symbols.FunctionSymbol, bodies []syntaxnodes.StatementNode) []boundnodes.BoundStatementNode {
    boundBodies := []boundnodes.BoundStatementNode{}

    for i, v := range bodies {
        // create a new binder
        bin := Binder{
            CurrentPackage: pck,
            CurrentFunction: syms[i],
            CurrentScope: NewScope(nil),
        }

        boundBodies = append(boundBodies, bin.bindStatement(v))
    }

    return boundBodies
}

// --------------------------------------------------------
// Statements
// --------------------------------------------------------
func (bin *Binder) bindStatement(stmt syntaxnodes.StatementNode) boundnodes.BoundStatementNode {
    if stmt.Type() == syntaxnodes.NT_DeclarationStmt {
        return bin.bindDeclarationStmt(stmt.(*syntaxnodes.DeclarationStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_ReturnStmt {
        return bin.bindReturnStmt(stmt.(*syntaxnodes.ReturnStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_WhileStmt {
        return bin.bindWhileStmt(stmt.(*syntaxnodes.WhileStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_FromToStmt {
        return bin.bindFromToStmt(stmt.(*syntaxnodes.FromToStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_ForStmt {
        return bin.bindForStmt(stmt.(*syntaxnodes.ForStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_LoopStmt {
        return bin.bindLoopStmt(stmt.(*syntaxnodes.LoopStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_BlockStmt {
        return bin.bindBlockStmt(stmt.(*syntaxnodes.BlockStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_ExpressionStmt {
        return bin.bindExpressionStmt(stmt.(*syntaxnodes.ExpressionStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_IfStmt {
        return bin.bindIfStmt(stmt.(*syntaxnodes.IfStatementNode))

    } else {

        error.Report(error.NewError(error.BND, stmt.Position(), "Unknown statement type '%s'!", stmt.Type()))
        return boundnodes.NewBoundExpressionStatementNode(stmt, boundnodes.NewBoundErrorExpressionNode(stmt))
    }
}

func (bin *Binder) bindDeclarationStmt(stmt *syntaxnodes.DeclarationStatementNode) *boundnodes.BoundDeclarationStatementNode {
    var initializer boundnodes.BoundExpressionNode
    var typ *symbols.TypeSymbol

    // do we have an explicit type or initializer?
    if !stmt.HasExplicitType && !stmt.HasInitializer {
        typ = compunit.GlobalDataTypeRegister["error"]
        error.Report(error.NewError(error.BND, stmt.Position(), "Variable declaration either needs explicit type declaration or initializer!"))
    }

    // if theres an explicit type -> resolve it
    if stmt.HasExplicitType {
        typ = LookupTypeClause(stmt.VarType)
    }

    // if we have an initializer -> bind it
    if stmt.HasInitializer {
        initializer = bin.bindExpression(stmt.Initializer)

        // if theres an explicit type -> make sure they match
        if stmt.HasExplicitType {
            initializer = bin.bindConversion(initializer, typ)

        // if not -> set the variable type
        } else {
            typ = initializer.ExprType()
        }
    }

    // create a variable symbol
    vari := symbols.NewLocalSymbol(stmt.VarName.Buffer, typ)

    // register this variable
    bin.CurrentScope.RegisterVariable(vari)

    // create bound node
    return boundnodes.NewBoundDeclarationStatementNode(stmt, vari, initializer, stmt.HasInitializer)
}

func (bin *Binder) bindReturnStmt(stmt *syntaxnodes.ReturnStatementNode) boundnodes.BoundStatementNode {
    var retValue boundnodes.BoundExpressionNode

    // bind the return value if it exists
    if stmt.HasExpression {
        retValue = bin.bindExpression(stmt.Expression)
    }

    // make sure the return value kind matches the function type
    if retValue == nil && !bin.CurrentFunction.ReturnType.Equal(compunit.GlobalDataTypeRegister["void"]) {
        error.Report(error.NewError(error.BND, stmt.Position(), "A function of type 'void' is not allowed to return a value!"))
        return boundnodes.NewBoundExpressionStatementNode(stmt, boundnodes.NewBoundErrorExpressionNode(stmt))
    }

    if retValue != nil {
        if !retValue.ExprType().Equal(bin.CurrentFunction.ReturnType) {
            error.Report(error.NewError(error.BND, stmt.Position(), "A function of type '%s' is not allowed to return a value of type '%s'!", bin.CurrentFunction.ReturnType.Name(), retValue.ExprType().Name()))
            return boundnodes.NewBoundExpressionStatementNode(stmt, boundnodes.NewBoundErrorExpressionNode(stmt))
        }
    }

    // create new bound node
    return boundnodes.NewBoundReturnStatementNode(stmt, retValue, stmt.HasExpression)
}


// --------------------------------------------------------
// Expressions
// --------------------------------------------------------
func (bin *Binder) bindExpression(expr syntaxnodes.ExpressionNode) boundnodes.BoundExpressionNode {
    return nil
}

// --------------------------------------------------------
// Utils
// --------------------------------------------------------
func (bin *Binder) bindConversion(expr boundnodes.BoundExpressionNode, typ *symbols.TypeSymbol) boundnodes.BoundExpressionNode {
    return expr
}

// --------------------------------------------------------
// Helper functions
// --------------------------------------------------------
func LookupType(name string, pos span.Span) *symbols.TypeSymbol {
    typ, ok := compunit.GlobalDataTypeRegister[name]

    if !ok {
        error.Report(error.NewError(error.BND, pos, "Unknown data type '%s'!", name))
        return compunit.GlobalDataTypeRegister["error"]
    }

    return typ
}

func LookupTypeClause(typ *syntaxnodes.TypeClauseNode) *symbols.TypeSymbol {
   
    // if the type clause does not exists -> void return type
    if typ == nil {
        return compunit.GlobalDataTypeRegister["void"]
    }

    // otherwise -> look up the type
    return LookupType(typ.TypeName.Buffer, typ.Position())
}
