// Evaluator - evaluator.go
// --------------------------------------------------------
// Finally, the very last step: running the program
// --------------------------------------------------------
package evaluator

import (
	"fmt"
	"reflect"
	"strconv"

	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compctl"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/span"
	"bytespace.network/rerect/symbols"
)

type Evaluator struct {
    Functions map[*symbols.FunctionSymbol]*boundnodes.BoundBlockStatementNode
    StackFrames []*StackFrame

    Globals map[symbols.VariableSymbol]interface{}
}

type StackFrame struct {
    InstPtr int
    Labels map[boundnodes.BoundLabel]int
    Locals map[symbols.VariableSymbol]interface{} 

    ReturnValue interface{}
    HasReturned bool
}

// --------------------------------------------------------
// Helpers
// --------------------------------------------------------
func (evl *Evaluator) setVar(vari symbols.VariableSymbol, val interface{}) {
    if vari.Type() == symbols.ST_Global {
        evl.Globals[vari] = val
    } else {
        evl.stackFrame().Locals[vari] = val
    }
}

func (evl *Evaluator) getVar(vari symbols.VariableSymbol) interface{} {
    if vari.Type() == symbols.ST_Global {

        val, ok := evl.Globals[vari]

        if !ok {
            error.Report(error.NewError(error.RNT, span.Internal(), "Someone fucked up the runtime global lookup :)"))
            return nil
        }

        return val
    } else {

        val, ok := evl.stackFrame().Locals[vari]

        if !ok {
            error.Report(error.NewError(error.RNT, span.Internal(), "Someone fucked up the runtime variable lookup :)"))
            return nil
        }

        return val
    }
}

func (evl *Evaluator) stackFrame() *StackFrame {
    return evl.StackFrames[len(evl.StackFrames)-1]
}

// --------------------------------------------------------
// Evaluation
// --------------------------------------------------------
func Evaluate(prg *compctl.CompilationResult) {
    // create a new evaluator
    evl := Evaluator{
        Functions: prg.Functions,
        StackFrames: make([]*StackFrame, 0),
        Globals: make(map[symbols.VariableSymbol]interface{}),
    }

    // create all globals
    for _, glb := range prg.Globals {
        // initialize with default value for each datatype
        evl.Globals[glb] = glb.VarType().Default
    }

    // look for a "main()" function in a "main" package
    var main *symbols.FunctionSymbol = nil
    for sym := range prg.Functions {
        if sym.FuncName == "main" && sym.ParentPackage.Name() == "main" {
            main = sym
            break
        }
    }

    // no entry point found
    if main == nil {
        error.Report(error.NewError(error.RNT, span.Internal(), "Could not find 'main()' function! An entry point is needed for execution."))
        return
    }

    // otherwise -> run main function
    evl.call(main, []interface{}{})
}

func (evl *Evaluator) call(fnc *symbols.FunctionSymbol, args []interface{}) interface{} {
    // create new call stack
    evl.StackFrames = append(evl.StackFrames, &StackFrame{
        Locals: make(map[symbols.VariableSymbol]interface{}),
        InstPtr: 0,
        ReturnValue: nil,
    })

    // register arguments
    for i := range fnc.Parameters {
        evl.setVar(fnc.Parameters[i], args[i])
    }

    // run the function body
    val := evl.run(evl.Functions[fnc])

    // destroy the stack frame
    evl.StackFrames = evl.StackFrames[:len(evl.StackFrames)-1]

    // return the functions return value
    return val
}

func (evl *Evaluator) callVM(fnc *symbols.FunctionSymbol, args []interface{}) interface{} {
   return fnc.Pointer(args) 
}

func (evl *Evaluator) run(body *boundnodes.BoundBlockStatementNode) interface{} {
    // index all labels
    lbls := make(map[boundnodes.BoundLabel]int)

    for i, v := range body.Statements {
        if v.Type() == boundnodes.BT_LabelIStmt {
            lbl := v.(*boundnodes.BoundLabelStatementNode)
            lbls[lbl.Label] = i
        }
    }

    // store the labels somewhere
    evl.stackFrame().Labels = lbls

    // execute the statements
    for evl.stackFrame().InstPtr < len(body.Statements) {

        // evaluate some cool statement
        evl.evalStatement(body.Statements[evl.stackFrame().InstPtr])

        // did we return?
        if evl.stackFrame().HasReturned {
            return evl.stackFrame().ReturnValue
        }

        // next instruction
        evl.stackFrame().InstPtr++
    }

    // someone forgot to return lol
    return evl.stackFrame().ReturnValue
}

func (evl *Evaluator) evalStatement(stmt boundnodes.BoundStatementNode) {
    if stmt.Type() == boundnodes.BT_DeclarationStmt {
        evl.evalDeclarationStatement(stmt.(*boundnodes.BoundDeclarationStatementNode))

    } else if stmt.Type() == boundnodes.BT_ReturnStmt {
        evl.evalReturnStatement(stmt.(*boundnodes.BoundReturnStatementNode))

    } else if stmt.Type() == boundnodes.BT_ExpressionStmt {
        evl.evalExpressionStatement(stmt.(*boundnodes.BoundExpressionStatementNode))

    } else if stmt.Type() == boundnodes.BT_GoToIStmt {
        evl.evalGotoStatement(stmt.(*boundnodes.BoundGotoStatementNode))

    } else if stmt.Type() == boundnodes.BT_GoToIfIStmt {
        evl.evalGotoIfStatement(stmt.(*boundnodes.BoundGotoIfStatementNode))

    } else if stmt.Type() == boundnodes.BT_DeleteIStmt {
        evl.evalDeleteStatement(stmt.(*boundnodes.BoundDeleteStatementNode))

    } else if stmt.Type() == boundnodes.BT_ApproachIStmt {
        evl.evalApproachStatement(stmt.(*boundnodes.BoundApproachStatementNode))

    } else if stmt.Type() == boundnodes.BT_LabelIStmt {
        // literally do nothing

    } else {
        error.Report(error.NewError(error.RNT, stmt.Source().Position(), "Statement evaluation not implemented! You should implement NOW! (%s)", stmt.Type()))
    }
}

func (evl *Evaluator) evalDeclarationStatement(stmt *boundnodes.BoundDeclarationStatementNode) {
    init := stmt.Variable.VarType().Default

    if stmt.HasInitializer {
        init = evl.evalExpression(stmt.Initializer)
    }

    evl.setVar(stmt.Variable, init)
}

func (evl *Evaluator) evalReturnStatement(stmt *boundnodes.BoundReturnStatementNode) {
    if stmt.HasReturnValue {
        evl.stackFrame().ReturnValue = evl.evalExpression(stmt.ReturnValue)
    }

    evl.stackFrame().HasReturned = true
}

func (evl *Evaluator) evalExpressionStatement(stmt *boundnodes.BoundExpressionStatementNode) {
    evl.evalExpression(stmt.Expression)
}

func (evl *Evaluator) evalGotoStatement(stmt *boundnodes.BoundGotoStatementNode) {
    idx, ok := evl.stackFrame().Labels[stmt.Label]

    if !ok {
        error.Report(error.NewError(error.RNT, span.Internal(), "Someone fucked up the runtime label lookup :)"))
        return
    }

    evl.stackFrame().InstPtr = idx
}

func (evl *Evaluator) evalGotoIfStatement(stmt *boundnodes.BoundGotoIfStatementNode) {
    expr := evl.evalExpression(stmt.Condition)

    if expr == true {
        idx, ok := evl.stackFrame().Labels[stmt.Label]

        if !ok {
            error.Report(error.NewError(error.RNT, span.Internal(), "Someone fucked up the runtime label lookup :)"))
            return
        }

        evl.stackFrame().InstPtr = idx
    }
}

func (evl *Evaluator) evalDeleteStatement(stmt *boundnodes.BoundDeleteStatementNode) {
    delete(evl.stackFrame().Locals, stmt.Variable)
}

func (evl *Evaluator) evalApproachStatement(stmt *boundnodes.BoundApproachStatementNode) {
    vari := evl.getVar(stmt.Iterator)
    val := evl.evalExpression(stmt.Target)

    if vari.(int32) < val.(int32) {
        evl.setVar(stmt.Iterator, vari.(int32)+1)
    } 

    if vari.(int32) > val.(int32) {
        evl.setVar(stmt.Iterator, vari.(int32)-1)
    } 
}

// --------------------------------------------------------
// Expressions
// --------------------------------------------------------
func (evl *Evaluator) evalExpression(expr boundnodes.BoundExpressionNode) interface{} {
    if expr.Type() == boundnodes.BT_LiteralExpr {
        return evl.evalLiteralExpression(expr.(*boundnodes.BoundLiteralExpressionNode))

    } else if expr.Type() == boundnodes.BT_AssignmentExpr {
        return evl.evalAssignmentExpression(expr.(*boundnodes.BoundAssignmentExpressionNode))

    } else if expr.Type() == boundnodes.BT_UnaryExpr {
        return evl.evalUnaryExpression(expr.(*boundnodes.BoundUnaryExpressionNode))

    } else if expr.Type() == boundnodes.BT_BinaryExpr {
        return evl.evalBinaryExpression(expr.(*boundnodes.BoundBinaryExpressionNode))

    } else if expr.Type() == boundnodes.BT_CallExpr {
        return evl.evalCallExpression(expr.(*boundnodes.BoundCallExpressionNode))

    } else if expr.Type() == boundnodes.BT_NameExpr {
        return evl.evalNameExpression(expr.(*boundnodes.BoundNameExpressionNode))

    } else if expr.Type() == boundnodes.BT_ConversionExpr {
        return evl.evalConversionExpression(expr.(*boundnodes.BoundConversionExpressionNode))

    } else {
        error.Report(error.NewError(error.RNT, expr.Source().Position(), "Expression evaluation not implemented! You should implement NOW! (%s)", expr.Type()))
        return nil
    }
}

func (evl *Evaluator) evalLiteralExpression(expr *boundnodes.BoundLiteralExpressionNode) interface{} {
    return expr.LiteralValue
}

func (evl *Evaluator) evalAssignmentExpression(expr *boundnodes.BoundAssignmentExpressionNode) interface{} {
    val := evl.evalExpression(expr.Value)
    evl.setVar(expr.Variable, val)
    return val
}

func (evl *Evaluator) evalUnaryExpression(expr *boundnodes.BoundUnaryExpressionNode) interface{} {
    operand := evl.evalExpression(expr.Operand)

    switch expr.Operator.Operation {
    case boundnodes.UO_Identity:
        return operand

    case boundnodes.UO_Negation:
        // Integers
        // --------
        if expr.Operator.Operand.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return -(operand.(int64))

        } else if expr.Operator.Operand.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return -(operand.(int32))

        } else if expr.Operator.Operand.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return -(operand.(int16))

        } else if expr.Operator.Operand.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return -(operand.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Operand.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return -(operand.(float64))

        } else if expr.Operator.Operand.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return -(operand.(float32))

        }

    case boundnodes.UO_LogicalNegation:
        // Booleans
        // --------
        if expr.Operator.Operand.Equal(compunit.GlobalDataTypeRegister["bool"]) {
            return !(operand.(bool))
        }
    }

    error.Report(error.NewError(error.RNT, expr.Source().Position(), "Unary operator not implemented! You should implement NOW!"))
    return nil
}

func (evl *Evaluator) evalBinaryExpression(expr *boundnodes.BoundBinaryExpressionNode) interface{} {
    left  := evl.evalExpression(expr.Left)
    right := evl.evalExpression(expr.Right)

    switch expr.Operator.Operation {
    case boundnodes.BO_Equal:
        return left == right

    case boundnodes.BO_UnEqual:
        return left != right

    case boundnodes.BO_LessThan:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) < (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) < (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) < (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  < (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) < (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) < (right.(float32))
        }
    
    case boundnodes.BO_LessEqual:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) <= (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) <= (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) <= (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  <= (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) <= (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) <= (right.(float32))
        }

    case boundnodes.BO_GreaterThan:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) > (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) > (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) > (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  > (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) > (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) > (right.(float32))
        }

    case boundnodes.BO_GreaterEqual:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) >= (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) >= (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) >= (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  >= (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) >= (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) >= (right.(float32))
        }

    case boundnodes.BO_Addition:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) + (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) + (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) + (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  + (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) + (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) + (right.(float32))
        }

    case boundnodes.BO_Subtraction:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) - (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) - (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) - (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  - (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) - (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) - (right.(float32))
        }

    case boundnodes.BO_Multiplication:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) * (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) * (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) * (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  * (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) * (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) * (right.(float32))
        }

    case boundnodes.BO_Division:
        // Integers
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["long"]) {
            return (left.(int64)) / (right.(int64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["int"]) {
            return (left.(int32)) / (right.(int32))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["word"]) {
            return (left.(int16)) / (right.(int16))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["byte"]) {
            return (left.(int8))  / (right.(int8))
        }

        // Floats
        // ------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["double"]) {
            return (left.(float64)) / (right.(float64))

        } else if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["float"]) {
            return (left.(float32)) / (right.(float32))
        }

    case boundnodes.BO_LogicalAnd:
        // Booleans
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["bool"]) {
            return (left.(bool)) && (right.(bool))
        }

    case boundnodes.BO_LogicalOr:
        // Booleans
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["bool"]) {
            return (left.(bool)) || (right.(bool))
        }

    case boundnodes.BO_Concat:
        // Booleans
        // --------
        if expr.Operator.Left.Equal(compunit.GlobalDataTypeRegister["string"]) {
            return (left.(string)) + (right.(string))
        }
    }

    error.Report(error.NewError(error.RNT, expr.Source().Position(), "Binary operator not implemented! You should implement NOW!"))
    return nil
}

func (evl *Evaluator) evalCallExpression(expr *boundnodes.BoundCallExpressionNode) interface{} {
    // evaluate all args
    args := []interface{}{}

    for _, arg := range expr.Arguments {
        args = append(args, evl.evalExpression(arg))
    } 

    // is this a native call?
    if expr.Function.IsVMFunction {
        // do a native call
        return evl.callVM(expr.Function, args)
    }

    // otherwise: call normally
    return evl.call(expr.Function, args)
}

func (evl *Evaluator) evalNameExpression(expr *boundnodes.BoundNameExpressionNode) interface{} {
    return evl.getVar(expr.Variable)
}

func (evl *Evaluator) evalConversionExpression(expr *boundnodes.BoundConversionExpressionNode) interface{} {
    val := evl.evalExpression(expr.Value)

    // Casting anything to 'any'
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["any"]) {
        return interface{}(val)
    }

    //fmt.Printf("Converting %s(%s) -> %s\n", expr.Value.ExprType().Name(), expr.Value.Type(), expr.TargetType.Name())

    // Casting to long
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["long"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return val;

        case int32:
            return int64(v)

        case int16:
            return int64(v)

        case int8:
            return int64(v)

        // Cross cast
        // ----------
        case float64:
            return int64(v)

        case float32:
            return int64(v)

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 64)
            if err != nil {
                panic(err)
            }

            return int64(vl)
        }
    }

    // Casting to int
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["int"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return int32(v);

        case int32:
            return v

        case int16:
            return int32(v)

        case int8:
            return int32(v)

        // Cross cast
        // ----------
        case float64:
            return int32(v)

        case float32:
            return int32(v)

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 32)
            if err != nil {
                panic(err)
            }

            return int32(vl)
        }
    }
    
    // Casting to word
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["word"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return int16(v);

        case int32:
            return int16(v)

        case int16:
            return v

        case int8:
            return int16(v)

        // Cross cast
        // ----------
        case float64:
            return int16(v)

        case float32:
            return int16(v)

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 16)
            if err != nil {
                panic(err)
            }

            return int16(vl)
        }
    }
    
    // Casting to byte
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["byte"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return int8(v);

        case int32:
            return int8(v)

        case int16:
            return int8(v)

        case int8:
            return v

        // Cross cast
        // ----------
        case float64:
            return int8(v)

        case float32:
            return int8(v)

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 8)
            if err != nil {
                panic(err)
            }

            return int8(vl)
        }
    }
    
    // Casting to double
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["double"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case float64:
            return v

        case float32:
            return float64(v)

        // Cross casts
        // -----------
        case int64:
            return float64(v);

        case int32:
            return float64(v)

        case int16:
            return float64(v)

        case int8:
            return float64(v)

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseFloat(v, 64)
            if err != nil {
                panic(err)
            }

            return float64(vl)
        }
    }

    // Casting to float
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["float"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case float64:
            return float32(v)

        case float32:
            return v

        // Cross casts
        // -----------
        case int64:
            return float32(v);

        case int32:
            return float32(v)

        case int16:
            return float32(v)

        case int8:
            return float32(v)

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseFloat(v, 32)
            if err != nil {
                panic(err)
            }

            return float32(vl)
        }
    }

    // Casting to string
    if expr.TargetType.Equal(compunit.GlobalDataTypeRegister["string"]) {
        switch v := val.(type) {

        // Integers
        // --------
        case int64:
            return fmt.Sprintf("%d", v)

        case int32:
            return fmt.Sprintf("%d", v)

        case int16:
            return fmt.Sprintf("%d", v)

        case int8:
            return fmt.Sprintf("%d", v)

        // Floats
        // ------
        case float64:
            return fmt.Sprintf("%d", v)

        case float32:
            return fmt.Sprintf("%d", v)

        // Booleans
        // --------
        case bool:
            if v {
                return "true"
            } else {
                return "false"
            }

        case ArrayInstance:
            return fmt.Sprintf("[%s]", v.Type.Name())

        // Strings
        // -------
        case string:
            return v
        }
    }

    // Casting to array
    if expr.TargetType.TypeGroup == symbols.ARR {
        switch v := val.(type) {
        case ArrayInstance:
            // only cast when the internal types match
            if v.Type.Equal(expr.TargetType) {
                return v
            }
        }
    }

    error.Report(error.NewError(error.RNT, expr.Source().Position(), "Unable to cast %s to %s!", reflect.TypeOf(val), expr.TargetType.Name()))
    return nil
}
