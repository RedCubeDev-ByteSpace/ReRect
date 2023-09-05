// Lowerer - lowerer.go
// --------------------------------------------------------
// The lowerer translates complex language concepts into
// nothing more than a series of labels and gotos
// --------------------------------------------------------
package lowerer

import (
	"fmt"

	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	packageprocessor "bytespace.network/rerect/package_processor"
	"bytespace.network/rerect/symbols"
)

// Lower function - outside wrapper
func Lower(file *packageprocessor.CompilationFile) {

    for _, sym := range file.Functions {
        stmt := file.FunctionBodies[sym]

        // if stmt is not a block statement -> create one
        if stmt.Type() != boundnodes.BT_BlockStmt {
            stmt = boundnodes.NewBoundBlockStatementNode(stmt.Source(), []boundnodes.BoundStatementNode{stmt})
        }

        // rewrite the body (simplify statements)
        stmt = rewriteStatement(stmt)

        // flatten the body into one long list of statements (instead of nested blocks)
        stmt = flatten(stmt.(*boundnodes.BoundBlockStatementNode))

        file.FunctionBodies[sym] = stmt
    }
}

func flatten(stmt *boundnodes.BoundBlockStatementNode) *boundnodes.BoundBlockStatementNode {
    stmts := []boundnodes.BoundStatementNode{}
    stack := []boundnodes.BoundStatementNode{}

    // push to any given stack
    pushTo := func(stck *[]boundnodes.BoundStatementNode, stmt boundnodes.BoundStatementNode) {
		*stck = append(*stck, stmt)
	}

    // copy one stack to another
	transferTo := func(stck *[]boundnodes.BoundStatementNode, stmt []boundnodes.BoundStatementNode) {
		*stck = append(*stck, stmt...)
	}

    // pop from any stack
	popFrom := func(stck *[]boundnodes.BoundStatementNode) boundnodes.BoundStatementNode {
		element := (*stck)[len(*stck)-1]
		*stck = (*stck)[:len(*stck)-1]
		return element
	}

    // push the current block statement onto the stack
	pushTo(&stack, stmt)

    for len(stack) > 0 {
        current := popFrom(&stack)

        // if this is a block statement -> it needs to be flattened
        if current.Type() == boundnodes.BT_BlockStmt {
            // create a local stack for this block statement
            // this is so we can inject nodes before these if needed
            localStack := []boundnodes.BoundStatementNode{}

            // collect variables which have been created in this scope
            vars := []symbols.VariableSymbol{}

            // push all elements onto the stack in reverse order (bc stacks be like that sometimes)
            currentBlock := current.(*boundnodes.BoundBlockStatementNode)
            for i := len(currentBlock.Statements) - 1; i >= 0; i-- {
                st := currentBlock.Statements[i]

                if st.Type() == boundnodes.BT_DeclarationStmt {
                    vars = append(vars, st.(*boundnodes.BoundDeclarationStatementNode).Variable)
                }

                pushTo(&localStack, st)
            }


            // delete all variables created here
            for _, v := range vars {
                pushTo(&stack, boundnodes.NewBoundDeleteStatementNode(current.Source(), v))
            }

            // transfer all elements from our local stack to the main one
            transferTo(&stack, localStack)


        // otherwise, simply copy this statement
        } else {
            stmts = append(stmts, current)
        }
    }

    return boundnodes.NewBoundBlockStatementNode(stmt.Source(), stmts)
}

// --------------------------------------------------------
// Helpers
// --------------------------------------------------------
var labelCounter int = 0
func generateLabel() boundnodes.BoundLabel {
    labelCounter++
    return boundnodes.BoundLabel(fmt.Sprintf("label%d", labelCounter))
}

// --------------------------------------------------------
// Statements
// --------------------------------------------------------
func rewriteStatement(stmt boundnodes.BoundStatementNode) boundnodes.BoundStatementNode {
    if stmt.Type() == boundnodes.BT_DeclarationStmt {
        return rewriteDeclarationStatement(stmt.(*boundnodes.BoundDeclarationStatementNode))

    } else if stmt.Type() == boundnodes.BT_ReturnStmt {
        return rewriteReturnStatement(stmt.(*boundnodes.BoundReturnStatementNode))

    } else if stmt.Type() == boundnodes.BT_WhileStmt {
        return rewriteWhileStatement(stmt.(*boundnodes.BoundWhileStatementNode))

    } else if stmt.Type() == boundnodes.BT_FromToStmt {
        return rewriteFromToStatement(stmt.(*boundnodes.BoundFromToStatementNode))

    } else if stmt.Type() == boundnodes.BT_ForStmt {
        return rewriteForStatement(stmt.(*boundnodes.BoundForStatementNode))

    } else if stmt.Type() == boundnodes.BT_LoopStmt {
        return rewriteLoopStatement(stmt.(*boundnodes.BoundLoopStatementNode))

    } else if stmt.Type() == boundnodes.BT_BlockStmt {
        return rewriteBlockStatement(stmt.(*boundnodes.BoundBlockStatementNode))

    } else if stmt.Type() == boundnodes.BT_ExpressionStmt {
        return rewriteExpressionStatement(stmt.(*boundnodes.BoundExpressionStatementNode))

    } else if stmt.Type() == boundnodes.BT_IfStmt {
        return rewriteIfStatement(stmt.(*boundnodes.BoundIfStatementNode))

    } else if stmt.Type() == boundnodes.BT_LabelIStmt {
        return stmt

    } else if stmt.Type() == boundnodes.BT_GoToIStmt {
        return stmt

    } else if stmt.Type() == boundnodes.BT_GoToIfIStmt {
        return stmt

    } else if stmt.Type() == boundnodes.BT_DeleteIStmt {
        return stmt

    } else if stmt.Type() == boundnodes.BT_ApproachIStmt {
        return stmt

    } else {
        error.Report(error.NewError(error.LWR, stmt.Source().Position(), "Unable to rewrite statement '%s', no rewriter implemented! You should implement NOW!", stmt.Type()))
        return stmt
    }
}

func rewriteDeclarationStatement(stmt *boundnodes.BoundDeclarationStatementNode) *boundnodes.BoundDeclarationStatementNode {
    if !stmt.HasInitializer {
        return stmt // nothing to rewrite
    }

    init := rewriteExpression(stmt.Initializer)
    return boundnodes.NewBoundDeclarationStatementNode(stmt.Source(), stmt.Variable, init, true)
}

func rewriteReturnStatement(stmt *boundnodes.BoundReturnStatementNode) *boundnodes.BoundReturnStatementNode {
    if !stmt.HasReturnValue {
        return stmt // nothing to rewrite
    }

    val := rewriteExpression(stmt.ReturnValue)
    return boundnodes.NewBoundReturnStatementNode(stmt.Source(), val, true)
}

func rewriteWhileStatement(stmt *boundnodes.BoundWhileStatementNode) boundnodes.BoundStatementNode {
    // while (<cond>) { <body> }
    // -------------------------
    // goto .continue
    // .body:
    //  <body>
    // .continue:
    // gotoif <condition> .body
    // .break:
    stmts := []boundnodes.BoundStatementNode{}
    bodyLabel := generateLabel()

    stmts = append(stmts, boundnodes.NewBoundGotoStatementNode(stmt.Source(), stmt.ContinueLabel()))
    stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.Source(), bodyLabel))
    stmts = append(stmts, rewriteStatement(stmt.Body))
    stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.Source(), stmt.ContinueLabel()))
    stmts = append(stmts, boundnodes.NewBoundGotoIfStatementNode(stmt.Source(), bodyLabel, rewriteExpression(stmt.Condtion)))
    stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.Source(), stmt.BreakLabel()))

    return boundnodes.NewBoundBlockStatementNode(stmt.Source(), stmts)
}

func rewriteFromToStatement(stmt *boundnodes.BoundFromToStatementNode) boundnodes.BoundStatementNode {
    // from <var> <- <lb> to <ub> { <body> }
    // -------------------------------------
    // var <ub> <- <ub>
    // for (var <var> <- <lb>; <lb> != <ub>; approach <var> <ub>) {
    //      <body>
    // }
    // delete <ub>

    stmts := []boundnodes.BoundStatementNode{}
    inttyp := compunit.GlobalDataTypeRegister["int"]

    // rewrite lower bound value
    lowerBound := rewriteExpression(stmt.LowerBound)

    // create lb variable declaration, deletion and access
    iteratorDeclaration := boundnodes.NewBoundDeclarationStatementNode(stmt.Source(), stmt.Iterator, lowerBound, true)
    iteratorExpression := boundnodes.NewBoundNameExpressionNode(stmt.Source(), stmt.Iterator)

    // rewrite upper bound value
    upperBound := rewriteExpression(stmt.UpperBound)

    // create a new variable symbol for this
    upperBoundVar := symbols.NewLocalSymbol("__upperBound", inttyp)

    // create ub variable declaration, deletion and access
    upperBoundDeclaration := boundnodes.NewBoundDeclarationStatementNode(stmt.Source(), upperBoundVar, upperBound, true)
    upperBoundDeletion := boundnodes.NewBoundDeleteStatementNode(stmt.Source(), upperBoundVar)
    upperBoundExpression := boundnodes.NewBoundNameExpressionNode(stmt.Source(), upperBoundVar)

    // create the while statement condition
    // lb != ub
    condition := boundnodes.NewBoundBinaryExpressionNode(
        stmt.Source(),
        boundnodes.NewBoundBinaryOperator(boundnodes.BO_UnEqual, inttyp, inttyp, compunit.GlobalDataTypeRegister["bool"]),
        iteratorExpression,
        upperBoundExpression,
    )

    // rewrite the original loop body
    body := rewriteStatement(stmt.Body)

    // create approch statement
    approach := boundnodes.NewBoundApproachStatementNode(stmt.Source(), stmt.Iterator, upperBoundExpression)

    // create internal while statement
    forstmt := boundnodes.NewBoundForStatementNode(stmt.Source(), iteratorDeclaration, condition, approach, body, stmt.BreakLbl, stmt.ContinueLbl)

    // assemble it all
    stmts = append(stmts, upperBoundDeclaration)
    stmts = append(stmts, rewriteStatement(forstmt))
    stmts = append(stmts, upperBoundDeletion)

    return boundnodes.NewBoundBlockStatementNode(stmt.Source(), stmts)
}

func rewriteForStatement(stmt *boundnodes.BoundForStatementNode) boundnodes.BoundStatementNode {
    // for (<declaration>; <condition>; <action>) { <body> }
    // -----------------------------------------------------
    // <declaration>
    // while (<condition>) {
    //  <body>
    //  <action>
    // }
    // delete <iterator>
    stmts := []boundnodes.BoundStatementNode{}

    // rewrite the original loop declaration
    decl := rewriteStatement(stmt.Initializer)
    cond := rewriteExpression(stmt.Condition)
    act  := rewriteStatement(stmt.Action)

    // rewrite the original loop body
    body := rewriteStatement(stmt.Body)

    // create internal while statement
    whilestmt := boundnodes.NewBoundWhileStatementNode(stmt.Source(), cond, boundnodes.NewBoundBlockStatementNode(
        stmt.Source(),
        []boundnodes.BoundStatementNode {
            body,
            boundnodes.NewBoundLabelStatementNode(stmt.Source(), stmt.ContinueLabel()),
            act,
        }),
        stmt.BreakLbl,
        generateLabel(),
    )


    // assemble it all
    stmts = append(stmts, decl)
    stmts = append(stmts, rewriteStatement(whilestmt))

    // delete variable if decl was used for declaration
    if decl.Type() == boundnodes.BT_DeclarationStmt {
        iterator := decl.(*boundnodes.BoundDeclarationStatementNode).Variable
        del := boundnodes.NewBoundDeleteStatementNode(stmt.Source(), iterator)

        stmts = append(stmts, del)
    }

    return boundnodes.NewBoundBlockStatementNode(stmt.Source(), stmts)
}

func rewriteLoopStatement(stmt *boundnodes.BoundLoopStatementNode) boundnodes.BoundStatementNode {
    // loop(<amount>) { <body> } 
    // -------------------------
    // from <i> <- 0 to <amount> {
    //   <body>
    // }

    // create iterator
    iterator := symbols.NewLocalSymbol("__iterator", compunit.GlobalDataTypeRegister["int"])

    // create zero initializer
    zeroLit := boundnodes.NewBoundLiteralExpressionNode(stmt.Source(), compunit.GlobalDataTypeRegister["int"], int32(0))

    // rewrite original upper bound
    upperBound := rewriteExpression(stmt.Amount)

    // rewrite original body
    body := rewriteStatement(stmt.Body)

    // create from-to statement
    fromtostmt := boundnodes.NewBoundFromToStatementNode(stmt.Source(), iterator, zeroLit, upperBound, body, stmt.BreakLabel(), stmt.ContinueLabel())

    return rewriteStatement(fromtostmt)
}

func rewriteBlockStatement(stmt *boundnodes.BoundBlockStatementNode) boundnodes.BoundStatementNode {
    stmts := []boundnodes.BoundStatementNode{}

    for _, v := range stmt.Statements {
        stmts = append(stmts, rewriteStatement(v))
    }

    return boundnodes.NewBoundBlockStatementNode(stmt.Source(), stmts)
}

func rewriteExpressionStatement(stmt *boundnodes.BoundExpressionStatementNode) boundnodes.BoundStatementNode {
    expr := rewriteExpression(stmt.Expression)
    return boundnodes.NewBoundExpressionStatementNode(stmt.Source(), expr)
}

func rewriteIfStatement(stmt *boundnodes.BoundIfStatementNode) boundnodes.BoundStatementNode {
    stmts := []boundnodes.BoundStatementNode{}

    // if (<cond>) {
    //   <body>
    // } 
    // -------------
    // gotoif <cond> .inner
    // goto .end
    // .inner:
    // <body>
    // .end:
    if !stmt.HasElse {

        inner := generateLabel()
        end := generateLabel()

        stmts = append(stmts, boundnodes.NewBoundGotoIfStatementNode(stmt.Source(), inner, stmt.Condition))
        stmts = append(stmts, boundnodes.NewBoundGotoStatementNode(stmt.Source(), end))
        stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.Source(), inner))
        stmts = append(stmts, rewriteStatement(stmt.Body))
        stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.SourceNode, end))

    // if (<cond>) {
    //   <body>
    // } else {
    //   <body>
    // }
    // -------------
    // gotoif <cond> .inner
    // goto .else
    // .inner:
    // <body>
    // goto .end
    // .else:
    // <body>
    // .end:
    } else {

        inner := generateLabel()
        els := generateLabel()
        end := generateLabel()

        stmts = append(stmts, boundnodes.NewBoundGotoIfStatementNode(stmt.Source(), inner, stmt.Condition))
        stmts = append(stmts, boundnodes.NewBoundGotoStatementNode(stmt.Source(), els))
        stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.Source(), inner))
        stmts = append(stmts, rewriteStatement(stmt.Body))
        stmts = append(stmts, boundnodes.NewBoundGotoStatementNode(stmt.Source(), end))
        stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.Source(), els))
        stmts = append(stmts, rewriteStatement(stmt.ElseBody))
        stmts = append(stmts, boundnodes.NewBoundLabelStatementNode(stmt.SourceNode, end))

    }

    return boundnodes.NewBoundBlockStatementNode(stmt.Source(), stmts)
}

// --------------------------------------------------------
// Expressions
// --------------------------------------------------------
func rewriteExpression(expr boundnodes.BoundExpressionNode) boundnodes.BoundExpressionNode {
    // god please end my suffering
    if expr.Type() == boundnodes.BT_LiteralExpr {
        return rewriteLiteralExpression(expr.(*boundnodes.BoundLiteralExpressionNode))
    } else if expr.Type() == boundnodes.BT_AssignmentExpr {
        return rewriteAssignmentExpression(expr.(*boundnodes.BoundAssignmentExpressionNode))
    } else if expr.Type() == boundnodes.BT_UnaryExpr {
        return rewriteUnaryExpression(expr.(*boundnodes.BoundUnaryExpressionNode))
    } else if expr.Type() == boundnodes.BT_BinaryExpr {
        return rewriteBinaryExpression(expr.(*boundnodes.BoundBinaryExpressionNode))
    } else if expr.Type() == boundnodes.BT_CallExpr {
        return rewriteCallExpression(expr.(*boundnodes.BoundCallExpressionNode))
    } else if expr.Type() == boundnodes.BT_NameExpr {
        return rewriteNameExpression(expr.(*boundnodes.BoundNameExpressionNode))
    } else if expr.Type() == boundnodes.BT_ConversionExpr {
        return rewriteConversionExpression(expr.(*boundnodes.BoundConversionExpressionNode))
    } else {
        error.Report(error.NewError(error.LWR, expr.Source().Position(), "Unable to rewrite expression '%s', no rewriter implemented! You should implement NOW!", expr.Type()))
        return expr
    }
}

func rewriteLiteralExpression(expr *boundnodes.BoundLiteralExpressionNode) boundnodes.BoundExpressionNode {
    return expr // no way
}

func rewriteAssignmentExpression(expr *boundnodes.BoundAssignmentExpressionNode) boundnodes.BoundExpressionNode {
    val := rewriteExpression(expr.Value)
    return boundnodes.NewBoundAssignmentExpressionNode(expr.Source(), expr.Variable, val)
}

func rewriteUnaryExpression(expr *boundnodes.BoundUnaryExpressionNode) boundnodes.BoundExpressionNode {
    operand := rewriteExpression(expr.Operand)
    return boundnodes.NewBoundUnaryExpressionNode(expr.Source(), expr.Operator, operand)
}

func rewriteBinaryExpression(expr *boundnodes.BoundBinaryExpressionNode) boundnodes.BoundExpressionNode {
    left := rewriteExpression(expr.Left)
    right := rewriteExpression(expr.Right)

    return boundnodes.NewBoundBinaryExpressionNode(expr.Source(), expr.Operator, left, right)
}

func rewriteCallExpression(expr *boundnodes.BoundCallExpressionNode) boundnodes.BoundExpressionNode {
    args := []boundnodes.BoundExpressionNode{}

    for _, v := range expr.Arguments {
        args = append(args, rewriteExpression(v))
    }

    return boundnodes.NewBoundCallExpressionNode(expr.Source(), expr.Function, args)
}

func rewriteNameExpression(expr *boundnodes.BoundNameExpressionNode) boundnodes.BoundExpressionNode {
    return expr
}

func rewriteConversionExpression(expr *boundnodes.BoundConversionExpressionNode) boundnodes.BoundExpressionNode {
    val := rewriteExpression(expr.Value)
    return boundnodes.NewBoundConversionExpressionNode(expr.Source(), val, expr.TargetType)
}
