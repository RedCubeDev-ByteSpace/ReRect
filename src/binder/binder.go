// Binder - binder.go
// --------------------------------------------------------
// The binder is a crucial part of the compilation process
// Its job is to do a semantic analysis on the parsed tree
// --------------------------------------------------------
package binder

import (
	"fmt"
	"strconv"

	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/lexer"
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

    BreakLabels []boundnodes.BoundLabel
    ContinueLabels []boundnodes.BoundLabel
    LabelCount int
}

func (bin *Binder) EnterNewScope() {
    // create new scope
    scp := NewScope(bin.CurrentScope)

    // use this new scope
    bin.CurrentScope = scp
}

func (bin *Binder) LeaveScope() {
    bin.CurrentScope = bin.CurrentScope.Parent
}

func (bin *Binder) PushLabels(brk boundnodes.BoundLabel, cnt boundnodes.BoundLabel) {
    bin.BreakLabels    = append(bin.BreakLabels, brk)
    bin.ContinueLabels = append(bin.ContinueLabels, cnt)
}

func (bin *Binder) PopLabels() {
	bin.BreakLabels    = bin.BreakLabels[:len(bin.BreakLabels)-1]
	bin.ContinueLabels = bin.ContinueLabels[:len(bin.ContinueLabels)-1]
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

        // register the function parameters as variables
        for _, v := range syms[i].Parameters {
            bin.CurrentScope.RegisterVariable(v)
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

    } else if stmt.Type() == syntaxnodes.NT_BreakStmt {
        return bin.bindBreakStmt(stmt.(*syntaxnodes.BreakStatementNode))

    } else if stmt.Type() == syntaxnodes.NT_ContinueStmt {
        return bin.bindContinueStmt(stmt.(*syntaxnodes.ContinueStatementNode))

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
            initializer = bin.bindConversion(initializer, typ, true)

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

func (bin *Binder) bindLoopBody(stmt syntaxnodes.StatementNode) (boundnodes.BoundStatementNode, boundnodes.BoundLabel, boundnodes.BoundLabel) {
   
    // generate loop labels
    bin.LabelCount++
    brk := boundnodes.BoundLabel(fmt.Sprintf("break%d", bin.LabelCount))
    cnt := boundnodes.BoundLabel(fmt.Sprintf("continue%d", bin.LabelCount))

    // push loop labels
    bin.PushLabels(brk, cnt)

    // bind the body
    body := bin.bindStatement(stmt)

    // pop the labels
    bin.PopLabels()

    return body, brk, cnt
}

func (bin *Binder) bindWhileStmt(stmt *syntaxnodes.WhileStatementNode) boundnodes.BoundStatementNode {
    // bind the while condition
    cond := bin.bindExpression(stmt.Expression)

    // make sure the expression is a boolean
    cond = bin.bindConversion(cond, compunit.GlobalDataTypeRegister["bool"], false)

    // bind the loop body
    bin.EnterNewScope()
    body, brk, cnt := bin.bindLoopBody(stmt.Body)
    bin.LeaveScope()

    // create new node
    return boundnodes.NewBoundWhileStatementNode(stmt, cond, body, brk, cnt)
}

func (bin *Binder) bindFromToStmt(stmt *syntaxnodes.FromToStatementNode) boundnodes.BoundStatementNode {
    // create the iterator
    vari := symbols.NewLocalSymbol(stmt.Iterator.Buffer, compunit.GlobalDataTypeRegister["int"])

    bin.EnterNewScope()
    bin.CurrentScope.RegisterVariable(vari) // will always work because the scope is empty

    // bind the lower bound
    lb := bin.bindExpression(stmt.LowerBound)

    // bind the upper bound
    ub := bin.bindExpression(stmt.UpperBound)

    // bind the loop body
    body, brk, cnt := bin.bindLoopBody(stmt.Body)

    bin.LeaveScope()

    // create new node
    return boundnodes.NewBoundFromToStatementNode(stmt, vari, lb, ub, body, brk, cnt)
}

func (bin *Binder) bindForStmt(stmt *syntaxnodes.ForStatementNode) boundnodes.BoundStatementNode {
    // register a new scope
    bin.EnterNewScope()

    // bind the initializer
    init := bin.bindStatement(stmt.Declaration)

    // bind the condition
    cond := bin.bindExpression(stmt.Condition)
    cond = bin.bindConversion(cond, compunit.GlobalDataTypeRegister["bool"], false)

    // bind the action
    action := bin.bindStatement(stmt.Action)

    // bind the body
    body, brk, cnt := bin.bindLoopBody(stmt.Body)

    // leave our new scope
    bin.LeaveScope()

    // create new node
    return boundnodes.NewBoundForStatementNode(stmt, init, cond, action, body, brk, cnt)
}

func (bin *Binder) bindLoopStmt(stmt *syntaxnodes.LoopStatementNode) boundnodes.BoundStatementNode {
    // register a new scope
    bin.EnterNewScope()
    
    // bind the amount of loops requested
    amount := bin.bindExpression(stmt.Expression)

    // bind the loop body
    body, brk, cnt := bin.bindLoopBody(stmt.Body)

    // leave our new scope
    bin.LeaveScope()

    // create a new node
    return boundnodes.NewBoundLoopStatementNode(stmt, amount, body, brk, cnt)
}

func (bin *Binder) bindBreakStmt(stmt *syntaxnodes.BreakStatementNode) boundnodes.BoundStatementNode {
    // are there actually any loops around rn?
    if len(bin.BreakLabels) == 0 {
        error.Report(error.NewError(error.BND, stmt.Position(), "Unable to use break statement outside of a loop!"))
        return boundnodes.NewBoundExpressionStatementNode(stmt, boundnodes.NewBoundErrorExpressionNode(stmt))
    }

    // if there are -> create a goto to the closest break label
    return boundnodes.NewBoundGotoStatementNode(stmt, bin.BreakLabels[len(bin.BreakLabels)-1])
}

func (bin *Binder) bindContinueStmt(stmt *syntaxnodes.ContinueStatementNode) boundnodes.BoundStatementNode {
    // are there actually any loops around rn?
    if len(bin.ContinueLabels) == 0 {
        error.Report(error.NewError(error.BND, stmt.Position(), "Unable to use continue statement outside of a loop!"))
        return boundnodes.NewBoundExpressionStatementNode(stmt, boundnodes.NewBoundErrorExpressionNode(stmt))
    }

    // if there are -> create a goto to the closest break label
    return boundnodes.NewBoundGotoStatementNode(stmt, bin.ContinueLabels[len(bin.ContinueLabels)-1])
}

func (bin *Binder) bindBlockStmt(stmt *syntaxnodes.BlockStatementNode) boundnodes.BoundStatementNode {
    // register a new scope
    bin.EnterNewScope()

    // bind all our statements
    stmts := []boundnodes.BoundStatementNode{}
    for _, v := range stmt.Statements {
        stmts = append(stmts, bin.bindStatement(v))
    }

    // leave our new scope
    bin.LeaveScope()

    // create a new node
    return boundnodes.NewBoundBlockStatementNode(stmt, stmts)
}

func (bin *Binder) bindExpressionStmt(stmt *syntaxnodes.ExpressionStatementNode) boundnodes.BoundStatementNode {
    // bind the expression in question
    expr := bin.bindExpression(stmt.Expression)

    // is this expression allowed to be a statement?
    if expr.Type() != boundnodes.BT_CallExpr && 
       expr.Type() != boundnodes.BT_AssignmentExpr &&
       expr.Type() != boundnodes.BT_ErrorExpr {

        error.Report(error.NewError(error.BND, stmt.Expression.Position(), "Expression of type '%s' is not allowed to be used as a statement!", expr.ExprType().Name()))
    }

    // create a new node
    return boundnodes.NewBoundExpressionStatementNode(stmt, expr)
}

func (bin *Binder) bindIfStmt(stmt *syntaxnodes.IfStatementNode) boundnodes.BoundStatementNode {
    // bind the condition
    cond := bin.bindExpression(stmt.Expression)

    // bind if block
    bin.EnterNewScope()
    body := bin.bindStatement(stmt.Body)
    bin.LeaveScope()

    // bind else block if it exists
    var elseBody boundnodes.BoundStatementNode

    if stmt.HasElseClause {
        bin.EnterNewScope()
        elseBody = bin.bindStatement(stmt.Else)
        bin.LeaveScope()
    }

    // create new node
    return boundnodes.NewBoundIfStatementNode(stmt, cond, body, elseBody, stmt.HasElseClause)
}

// --------------------------------------------------------
// Expressions
// --------------------------------------------------------
func (bin *Binder) bindExpression(expr syntaxnodes.ExpressionNode) boundnodes.BoundExpressionNode {

    if expr.Type() == syntaxnodes.NT_LiteralExpr {
        return bin.bindLiteralExpression(expr.(*syntaxnodes.LiteralExpressionNode))

    } else if expr.Type() == syntaxnodes.NT_ParenthesizedExpr {
        return bin.bindParenthesizedExpression(expr.(*syntaxnodes.ParenthesizedExpressionNode))

    } else if expr.Type() == syntaxnodes.NT_AssignmentExpr {
        return bin.bindAssignmentExpression(expr.(*syntaxnodes.AssignmentExpressionNode))

    } else if expr.Type() == syntaxnodes.NT_UnaryExpr {
        return bin.bindUnaryExpression(expr.(*syntaxnodes.UnaryExpressionNode))

    } else if expr.Type() == syntaxnodes.NT_BinaryExpr {
        return bin.bindBinaryExpression(expr.(*syntaxnodes.BinaryExpressionNode))
    
    } else if expr.Type() == syntaxnodes.NT_CallExpr {
        return bin.bindCallExpression(expr.(*syntaxnodes.CallExpressionNode))

    } else if expr.Type() == syntaxnodes.NT_NameExpr {
        return bin.bindNameExpression(expr.(*syntaxnodes.NameExpressionNode))

    } else {
        error.Report(error.NewError(error.BND, expr.Position(), "Unknown expression type '%s'!", expr.Type()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }
}

func (bin *Binder) bindLiteralExpression(expr *syntaxnodes.LiteralExpressionNode) boundnodes.BoundExpressionNode {

    // literal value
    var value interface{}

    // literal type
    var typ *symbols.TypeSymbol

    // evaluate the literal expression
    if expr.Literal.Type == lexer.TT_String {
        value = expr.Literal.Buffer
        typ = compunit.GlobalDataTypeRegister["string"]

    } else if expr.Literal.Type == lexer.TT_KW_True {
        value = true
        typ = compunit.GlobalDataTypeRegister["bool"]

    } else if expr.Literal.Type == lexer.TT_KW_False {
        value = false
        typ = compunit.GlobalDataTypeRegister["bool"]

    } else if expr.Literal.Type == lexer.TT_Integer {
        val, err := strconv.Atoi(expr.Literal.Buffer) 
        
        if err != nil {
            error.Report(error.NewError(error.BND, expr.Position(), "Could not convert '%s' to an integer!", expr.Literal.Buffer))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        value = val
        typ = compunit.GlobalDataTypeRegister["int"]

    } else if expr.Literal.Type == lexer.TT_Float {
        val, err := strconv.ParseFloat(expr.Literal.Buffer, 32)
        
        if err != nil {
            error.Report(error.NewError(error.BND, expr.Position(), "Could not convert '%s' to a float!", expr.Literal.Buffer))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        value = val
        typ = compunit.GlobalDataTypeRegister["float"]

    } else {
        error.Report(error.NewError(error.BND, expr.Position(), "Expected literal value, got: '%s' (%s)!", expr.Literal.Buffer, expr.Literal.Type))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // create a new node
    return boundnodes.NewBoundLiteralExpressionNode(expr, typ, value)
}

func (bin *Binder) bindParenthesizedExpression(expr *syntaxnodes.ParenthesizedExpressionNode) boundnodes.BoundExpressionNode {
    // bind the inner expression
    exp := bin.bindExpression(expr.Expression)

    // done lol
    return exp
}

func (bin *Binder) bindAssignmentExpression(expr *syntaxnodes.AssignmentExpressionNode) boundnodes.BoundExpressionNode {
    // look up variable
    vari := bin.CurrentScope.LookupVariable(expr.VarName.Buffer)

    // did we find one?
    if vari == nil {
        error.Report(error.NewError(error.BND, expr.Position(), "Could not find variable called '%s'!", expr.VarName.Buffer))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // bind assignment value
    val := bin.bindExpression(expr.Expression)

    // make sure the data types match
    val = bin.bindConversion(val, vari.VarType(), false)

    // cool
    return boundnodes.NewBoundAssignmentExpressionNode(expr, vari, val)
}

func (bin *Binder) bindUnaryExpression(expr *syntaxnodes.UnaryExpressionNode) boundnodes.BoundExpressionNode {
    // bind the operand
    operand := bin.bindExpression(expr.Operand)

    // bind a unary operator
    op := boundnodes.GetUnaryOperator(expr.Operator.Type, operand.ExprType())

    // did we find a fitting operator?
    if op == nil {
        error.Report(error.NewError(error.BND, expr.Position(), "Operator '%s' is not defined for data type '%s'!", expr.Operator.Type, operand.ExprType().Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    return boundnodes.NewBoundUnaryExpressionNode(expr, op, operand)
}

func (bin *Binder) bindBinaryExpression(expr *syntaxnodes.BinaryExpressionNode) boundnodes.BoundExpressionNode {
    // bind the left and right sides
    left  := bin.bindExpression(expr.Left)
    right := bin.bindExpression(expr.Right)

    // bind a binary operator
    op := boundnodes.GetBinaryOperator(expr.Operator.Type, left.ExprType(), right.ExprType())

    // did we find a fitting operator?
    if op == nil {
        error.Report(error.NewError(error.BND, expr.Position(), "Operator '%s' is not defined for data types '%s' and '%s'!", expr.Operator.Type, left.ExprType().Name(), right.ExprType().Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // do we need to up cast the left side?
    if left.ExprType().TypeSize < right.ExprType().TypeSize {
        left = bin.bindConversion(left, right.ExprType(), false)
    }

    // do we need to up cast the right side?
    if left.ExprType().TypeSize > right.ExprType().TypeSize {
        right = bin.bindConversion(right, left.ExprType(), false)
    }

    return boundnodes.NewBoundBinaryExpressionNode(expr, op, left, right)
}

func (bin *Binder) bindCallExpression(expr *syntaxnodes.CallExpressionNode) boundnodes.BoundExpressionNode {

    // is this actually a cast?
    if len(expr.Parameters) == 1 {
        // are we calling a type name?
        typ := LookupType(expr.Identifier.Buffer, expr.Identifier.Position, true)
        
        // if so -> bind a conversion
        if typ != nil {
            exp := bin.bindExpression(expr.Parameters[0])
            return bin.bindConversion(exp, typ, true)
        }
    }

    // otherwise -> bind a call
    // ------------------------

    // lookup the function
    fnc := bin.LookupFunction(expr.Identifier.Buffer)

    if fnc == nil {
        error.Report(error.NewError(error.BND, expr.Identifier.Position, "Could not find function '%s'!", expr.Identifier.Buffer))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }
    
    // was the right amount of arguments given?
    if len(fnc.Parameters) != len(expr.Parameters) {
        error.Report(error.NewError(error.BND, expr.Position(), "Function '%s' expects %d arguments, got: %d!", fnc.FuncName, len(fnc.Parameters), len(expr.Parameters)))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // bind all args
    args := []boundnodes.BoundExpressionNode{}
    for _, v := range expr.Parameters {
        args = append(args, bin.bindExpression(v))
    }

    // make sure the datatypes match up
    for i := range fnc.Parameters {
        if !fnc.Parameters[i].VarType().Equal(args[i].ExprType()) {
            error.Report(error.NewError(error.BND, expr.Parameters[i].Position(), "Function '%s' expects an argument of type '%s' at index %d, got: '%s'!", fnc.FuncName, fnc.Parameters[i].VarType().Name(), i, args[i].ExprType().Name()))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }
    }

    // ok cool
    return boundnodes.NewBoundCallExpressionNode(expr, fnc, args)
}

func (bin *Binder) bindNameExpression(expr *syntaxnodes.NameExpressionNode) boundnodes.BoundExpressionNode {
    // look up variable
    vari := bin.CurrentScope.LookupVariable(expr.Identifier.Buffer)

    // did we find one?
    if vari == nil {
        error.Report(error.NewError(error.BND, expr.Position(), "Could not find variable called '%s'!", expr.Identifier.Buffer))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // ok cool
    return boundnodes.NewBoundNameExpressionNode(expr, vari)
}

// --------------------------------------------------------
// Utils
// --------------------------------------------------------
func (bin *Binder) bindConversion(expr boundnodes.BoundExpressionNode, typ *symbols.TypeSymbol, explicit bool) boundnodes.BoundExpressionNode {
    // lookup this converion 
    con := boundnodes.ClassifyConversion(expr.ExprType(), typ)

    // no conversion exists
    if con == boundnodes.CT_None {
        error.Report(error.NewError(error.BND, expr.Source().Position(), "Unable to convert type '%s' into '%s'!", expr.ExprType().Name(), typ.Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr.Source())
    }

    // explicit conversion exists, but explicit isnt allowed
    if con == boundnodes.CT_Explicit && !explicit {
        error.Report(error.NewError(error.BND, expr.Source().Position(), "Unable to implicitly convert type '%s' into '%s'! An explicit conversion exists (are you missing a cast?)", expr.ExprType().Name(), typ.Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr.Source())
    }

    // otherwise -> we cool
    return boundnodes.NewBoundConversionExpressionNode(expr.Source(), expr, typ)
}

// --------------------------------------------------------
// Helper functions
// --------------------------------------------------------
func LookupType(name string, pos span.Span, canfail bool) *symbols.TypeSymbol {
    typ, ok := compunit.GlobalDataTypeRegister[name]

    if !ok {
        if canfail {
            return nil
        }

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
    return LookupType(typ.TypeName.Buffer, typ.Position(), false)
}

func (bin *Binder) LookupFunction(name string) *symbols.FunctionSymbol {
    // look in local package first 
    fnc := LookupFunctionInPackage(bin.CurrentPackage, name)

    if fnc != nil {
        return fnc
    }

    // if we didnt find anything -> start looking through included packages
    for _, pname := range bin.CurrentPackage.IncludedPackages {
        pck := compunit.GetPackage(pname)

        fnc := LookupFunctionInPackage(pck, name)
        if fnc != nil {
            return fnc
        }
    }

    // we got nothin man
    return nil
}

func LookupFunctionInPackage(pck *symbols.PackageSymbol, name string) *symbols.FunctionSymbol {
    for _, v := range pck.Functions {
        if v.FuncName == name {
            return v
        }
    }

    return nil
}
