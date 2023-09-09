// Evaluator - evaluator.go
// --------------------------------------------------------
// Finally, the very last step: running the program
// --------------------------------------------------------
package evaluator

import (
	"reflect"

	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compctl"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	evalobjects "bytespace.network/rerect/eval_objects"
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

    This interface{}

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
        evl.Globals[glb] = evl.getDefault(glb.VarType())
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

// Functions
// ---------
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
   return fnc.FunctionPointer(args) 
}

// Methods
// -------
func (evl *Evaluator) callMethod(fnc *symbols.FunctionSymbol, instance interface{}, args []interface{}) interface{} {
    // create new call stack
    evl.StackFrames = append(evl.StackFrames, &StackFrame{
        Locals: make(map[symbols.VariableSymbol]interface{}),
        InstPtr: 0,
        ReturnValue: nil,
        This: instance,
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

func (evl *Evaluator) callMethodVM(fnc *symbols.FunctionSymbol, instance interface{}, args []interface{}) interface{} {
   return fnc.MethodPointer(instance, args) 
}

// Run any sort of function body (function or method)
// --------------------------------------------------
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
    init := evl.getDefault(stmt.Variable.VarType())

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

    } else if expr.Type() == boundnodes.BT_AccessCallExpr {
        return evl.evalAccessCallExpression(expr.(*boundnodes.BoundAccessCallExpressionNode))

    } else if expr.Type() == boundnodes.BT_NameExpr {
        return evl.evalNameExpression(expr.(*boundnodes.BoundNameExpressionNode))

    } else if expr.Type() == boundnodes.BT_ConversionExpr {
        return evl.evalConversionExpression(expr.(*boundnodes.BoundConversionExpressionNode))

    } else if expr.Type() == boundnodes.BT_MakeArrayExpr {
        return evl.evalMakeArrayExpression(expr.(*boundnodes.BoundMakeArrayExpressionNode))
        
    } else if expr.Type() == boundnodes.BT_ArrayIndexExpr {
        return evl.evalArrayIndexExpression(expr.(*boundnodes.BoundArrayIndexExpressionNode))

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

    // classic variable assignment
    if expr.Expression.Type() == boundnodes.BT_NameExpr {
        evl.setVar(expr.Expression.(*boundnodes.BoundNameExpressionNode).Variable, val)

    // array index assignment
    } else if expr.Expression.Type() == boundnodes.BT_ArrayIndexExpr {
        exp := expr.Expression.(*boundnodes.BoundArrayIndexExpressionNode)

        // get the source array
        src := evl.evalExpression(exp.SourceArray).(*evalobjects.ArrayInstance)
        idx := evl.evalExpression(exp.Index).(int32)

        // assign the value
        src.Elements[idx] = val
    }

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

func (evl *Evaluator) evalAccessCallExpression(expr *boundnodes.BoundAccessCallExpressionNode) interface{} {
    // evaluate the call source
    src := evl.evalExpression(expr.Expression)

    // evaluate all args
    args := []interface{}{}

    for _, arg := range expr.Arguments {
        args = append(args, evl.evalExpression(arg))
    } 

    // is this a native call?
    if expr.Function.IsVMFunction {
        // do a native call
        return evl.callMethodVM(expr.Function, src, args)
    }

    // otherwise: call normally
    return evl.callMethod(expr.Function, src, args)
}

func (evl *Evaluator) evalNameExpression(expr *boundnodes.BoundNameExpressionNode) interface{} {
    return evl.getVar(expr.Variable)
}

func (evl *Evaluator) evalConversionExpression(expr *boundnodes.BoundConversionExpressionNode) interface{} {
    val := evl.evalExpression(expr.Value)

    res, ok := evalobjects.EvalConversion(val, expr.TargetType)
    if ok {
        return res
    }

    error.Report(error.NewError(error.RNT, expr.Source().Position(), "Unable to cast %s to %s!", reflect.TypeOf(val), expr.TargetType.Name()))
    return nil
}


func (evl *Evaluator) evalMakeArrayExpression(expr *boundnodes.BoundMakeArrayExpressionNode) interface{} {
    // This is a length defined array
    if !expr.HasInitializer {
        // figure out the length of the new array
        length := evl.evalExpression(expr.Length).(int32)

        // create an array object
        arr := &evalobjects.ArrayInstance {
            Type: expr.ArrType,
            Elements: make([]interface{}, 0),
        }

        // fill the array with default values
        for i := 0; int32(i) < length; i++ {
            arr.Elements = append(arr.Elements, evl.getDefault(expr.ArrType.SubTypes[0]))
        }

        return arr

    // This is an element defined array
    } else {
        // create new empty array
        arr := &evalobjects.ArrayInstance {
            Type: expr.ArrType,
            Elements: make([]interface{}, 0),
        }

        // insert all defined elements
        for _, v := range expr.Initializer {
            arr.Elements = append(arr.Elements, evl.evalExpression(v))
        }

        return arr
    }
}

func (evl *Evaluator) evalArrayIndexExpression(expr *boundnodes.BoundArrayIndexExpressionNode) interface{} {
    // evaluate the source
    src := evl.evalExpression(expr.SourceArray).(*evalobjects.ArrayInstance)

    // evaluate the index
    idx := evl.evalExpression(expr.Index).(int32)

    // make sure the index isnt out of bounds
    if idx < 0 || idx >= int32(len(src.Elements)) {     
        error.Report(error.NewError(error.RNT, expr.Source().Position(), "Index out of bounds! (index: %d, length of array: %d)", idx, len(src.Elements)))
        return evl.getDefault(src.Type.SubTypes[0])
    }
    
    // if its not -> return the value at the index
    return src.Elements[idx]
}

// --------------------------------------------------------
// Helpers
// --------------------------------------------------------
func (evl *Evaluator) getDefault(typ *symbols.TypeSymbol) interface{} {
    // Arrays need some special care because theyre reference types
    if typ.TypeGroup == symbols.ARR {
        return &evalobjects.ArrayInstance{
            Type: typ,
            Elements: make([]interface{}, 0),
        }
    }

    // otherwise: return the predefined default
    return typ.Default
}
