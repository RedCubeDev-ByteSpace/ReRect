// Parser - parser.go
// ---------------------------------------------------------------------
// Contains the main parser and syntax tree builder
// ---------------------------------------------------------------------
package parser

import (
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
	"bytespace.network/rerect/syntaxnodes"
)

// Parser struct
// -------------
type Parser struct {
    Source []lexer.Token
    SourceFileIdx int

    Index int
    Length int

    Members []syntaxnodes.MemberNode
}

// --------------------------------------------------------
// Helpers
// --------------------------------------------------------
func (prs *Parser) current() lexer.Token {
    return prs.peek(0)
}

func (prs *Parser) peek(offset int) lexer.Token {
    if offset + prs.Index >= prs.Length {
        // Fabricate EOF token when out of bounds
        return lexer.Token{
            Type: lexer.TT_EOF,
            Position: span.Span{
                File: prs.SourceFileIdx,
                FromIdx: prs.Length-1,
                ToIdx: prs.Length-1,
            },
        }
    }

    // otherwise return requested token
    return prs.Source[offset + prs.Index]
}

func (prs *Parser) step(size int) {
    prs.Index += size
}

func (prs *Parser) consume(typ lexer.TokenType) lexer.Token {
    // token doesnt match the expected type
    if prs.current().Type != typ {
        // report this error
        error.Report(error.NewError(error.PRS, prs.current().Position, "Expected token '%s', instead got: '%s'!", typ, prs.current().Type))
       
        prs.step(1)

        // fabricate a token of this kind to keep the compilation going
        return lexer.Token {
           Type: typ,
        }
    }

    // if it does match
    prs.step(1) // next token
    return prs.peek(-1) // return the one we were just at
}

// --------------------------------------------------------
// Parsing
// --------------------------------------------------------
func Parse(tokens []lexer.Token) []syntaxnodes.MemberNode {
    // if theres no tokens, theres nothing to parse
    if len(tokens) == 0 {
        return make([]syntaxnodes.MemberNode, 0)
    }

    // create parser instance
    prs := Parser {
        Source: tokens,
        SourceFileIdx: tokens[0].Position.File,
        Length: len(tokens),
        Members: make([]syntaxnodes.MemberNode, 0),
    }

    prs.parseMembers()

    return prs.Members
}

// --------------------------------------------------------
// Members
// --------------------------------------------------------
func (prs *Parser) parseMembers() {
    for prs.current().Type != lexer.TT_EOF {
        prs.parseMember()
    }
}

func (prs *Parser) parseMember() {
    var mem syntaxnodes.MemberNode

    // load <package> [include]
    if prs.current().Type == lexer.TT_KW_Load {
        mem = prs.parseLoadMember()
    
    // package <package>
    } else if prs.current().Type == lexer.TT_KW_Package {
        mem = prs.parsePackageMember()
    
    // function <name>(<args>) { ... }
    } else if prs.current().Type == lexer.TT_KW_Function {
        mem = prs.parseFunctionMember()

    // var <varname> <type>
    } else if prs.current().Type == lexer.TT_KW_Var {
        mem = prs.parseGlobalMember()

    // anything else -> error
    } else {
        error.Report(error.NewError(error.PRS, prs.current().Position, "Expected member, instead got: '%s'!", prs.current().Type))
        prs.step(1)

        return
    }
   
    // if this isnt a function -> require a semicolon
    if mem.Type() != syntaxnodes.NT_Function {
        prs.consume(lexer.TT_Semicolon)
    }

    prs.Members = append(prs.Members, mem)
}

func (prs *Parser) parseLoadMember() *syntaxnodes.LoadNode {
    // consume 'load'
    kw := prs.consume(lexer.TT_KW_Load)

    // consume library name
    lib := prs.consume(lexer.TT_Identifier)

    // optionally consume 'includel'
    var kwinclude lexer.Token
    var hasinclude bool

    if prs.current().Type == lexer.TT_KW_Include {
        kwinclude = prs.consume(lexer.TT_KW_Include)
        hasinclude = true
    }

    // create member node
    return syntaxnodes.NewLoadNode(kw, lib, kwinclude, hasinclude)
}

func (prs *Parser) parsePackageMember() *syntaxnodes.PackageNode {
    // consume 'package'
    kw := prs.consume(lexer.TT_KW_Package)

    // consume package name
    lib := prs.consume(lexer.TT_Identifier)

    // create member node
    return syntaxnodes.NewPackageNode(kw, lib)
}

func (prs *Parser) parseFunctionMember() *syntaxnodes.FunctionNode {
    // consume 'function'
    kw := prs.consume(lexer.TT_KW_Function)

    // consume function name
    id := prs.consume(lexer.TT_Identifier)

    // consume '('
    prs.consume(lexer.TT_OpenParenthesis)

    // parse parameters
    var params []*syntaxnodes.ParameterClauseNode
    for prs.current().Type != lexer.TT_CloseParenthesis {
        params = append(params, prs.parseParameterClause())
    }

    // consume ')'
    prs.consume(lexer.TT_CloseParenthesis)
   
    // parse return type
    var retType *syntaxnodes.TypeClauseNode
    hasReturnType := false
    if prs.current().Type != lexer.TT_OpenBraces && prs.current().Type != lexer.TT_Colon {
        retType = prs.parseTypeClause()
        hasReturnType = true
    }

    // parse the body
    var body syntaxnodes.StatementNode

    // the body can either be a single line, a la:
    // function a(): Print("hello"); 
    if prs.current().Type == lexer.TT_Colon {
        prs.consume(lexer.TT_Colon)
        body = prs.parseStatement()

    // or a traditional block statement
    // function b() { ... }
    } else {
        body = prs.parseBlockStatement()
    }

    return syntaxnodes.NewFunctionNode(kw, id, params, retType, hasReturnType, body)
}

func (prs *Parser) parseGlobalMember() *syntaxnodes.GlobalNode {
    // consume 'var' keyword
    kw := prs.consume(lexer.TT_KW_Var)

    // consume variable name
    id := prs.consume(lexer.TT_Identifier)

    // consume variable type 
    typ := prs.parseTypeClause()

    // create a new member node
    return syntaxnodes.NewGlobalNode(kw, id, typ)
}


// --------------------------------------------------------
// Clauses
// --------------------------------------------------------
func (prs *Parser) parseParameterClause() *syntaxnodes.ParameterClauseNode {
    // consume param name 
    id := prs.consume(lexer.TT_Identifier)

    // consume parm type
    typ := prs.parseTypeClause()

    return syntaxnodes.NewParameterClauseNode(id, typ)
}

func (prs *Parser) parseTypeClause() *syntaxnodes.TypeClauseNode {
    // consume type name
    id := prs.consume(lexer.TT_Identifier)

    var subtypes []*syntaxnodes.TypeClauseNode

    // check if theres are subtypes (look for '[')
    if prs.current().Type == lexer.TT_OpenBrackets {
        // consume '['
        prs.consume(lexer.TT_OpenBrackets)

        for prs.current().Type != lexer.TT_CloseBrackets {
            subtypes = append(subtypes, prs.parseTypeClause())

            // if we find a comma -> absorb it
            if prs.current().Type == lexer.TT_Comma {
                prs.consume(lexer.TT_Comma)

            // otherwise -> break
            } else {
                break
            }
        }

        // consume ']'
        prs.consume(lexer.TT_CloseBrackets)
    }

    // create new clause
    return syntaxnodes.NewTypeClauseNode(id, subtypes)
}

// --------------------------------------------------------
// Statements
// --------------------------------------------------------
func (prs *Parser) parseStatement() syntaxnodes.StatementNode {
    var stmt syntaxnodes.StatementNode 

    

    // var <name> [type] [<- <initializer>] 
    if prs.current().Type == lexer.TT_KW_Var {
        stmt = prs.parseDeclarationStatement()
    
    // return [val]
    } else if prs.current().Type == lexer.TT_KW_Return {
        stmt = prs.parseReturnStatement()

    // while (<cond>) { ... }
    } else if prs.current().Type == lexer.TT_KW_While {
        stmt = prs.parseWhileStatement()

    // From <iterator> <- <lower bound> to <higher bound>
    } else if prs.current().Type == lexer.TT_KW_From {
        stmt = prs.parseFromToStatement()
    
    // For(<decl>; <cond>; <action>) {...}
    } else if prs.current().Type == lexer.TT_KW_For {
        stmt = prs.parseForStatement()
    
    // Loop(<amount>) { ... }
    } else if prs.current().Type == lexer.TT_KW_Loop {
        stmt = prs.parseLoopStatement()

    // Break;
    } else if prs.current().Type == lexer.TT_KW_Break {
        stmt = prs.parseBreakStatement()

    // Continue;
    } else if prs.current().Type == lexer.TT_KW_Continue {
        stmt = prs.parseContinueStatement()
    
    // if (<cond>) { ... } [else { ... }]
    } else if prs.current().Type == lexer.TT_KW_If {
        stmt = prs.parseIfStatement()

    // { [statements] }
    } else if prs.current().Type == lexer.TT_OpenBraces {
        stmt = prs.parseBlockStatement()
    
    // <literally any expression>
    } else {
        stmt = prs.parseExpressionStatement()
    }

    // if this isnt a block statement
    if stmt.Type() == syntaxnodes.NT_ReturnStmt      ||
       stmt.Type() == syntaxnodes.NT_DeclarationStmt || 
       stmt.Type() == syntaxnodes.NT_BreakStmt || 
       stmt.Type() == syntaxnodes.NT_ContinueStmt || 
       stmt.Type() == syntaxnodes.NT_ExpressionStmt  {
        // require a semicolon
        prs.consume(lexer.TT_Semicolon)
    }

    return stmt
}

func (prs *Parser) parseDeclarationStatement() *syntaxnodes.DeclarationStatementNode {
    // consume 'var' keyword
    kw := prs.consume(lexer.TT_KW_Var)

    // consume variable name
    id := prs.consume(lexer.TT_Identifier)

    // (optional) consume explicit type
    var varType *syntaxnodes.TypeClauseNode
    hasExplicitType := false

    if prs.current().Type == lexer.TT_Identifier {
        varType = prs.parseTypeClause()
        hasExplicitType = true
    }

    // (optional) consume initializer
    var initializer syntaxnodes.ExpressionNode
    hasInitializer := false
    if prs.current().Type != lexer.TT_Semicolon {
        // consume assignment arrow
        prs.consume(lexer.TT_LeftArrow)

        // consume the value
        initializer = prs.parseExpression()
         
        hasInitializer = true
    }

    return syntaxnodes.NewDeclarationStatementNode(kw, id, varType, hasExplicitType, initializer, hasInitializer)
}

func (prs *Parser) parseReturnStatement() *syntaxnodes.ReturnStatementNode {
    // consume 'return' keyword
    kw := prs.consume(lexer.TT_KW_Return)

    // (option) consume return value
    var value syntaxnodes.ExpressionNode
    hasValue := false

    if prs.current().Type != lexer.TT_Semicolon {
        // parse retrun value
        value = prs.parseExpression()

        hasValue = true
    }

    return syntaxnodes.NewReturnStatementNode(kw, value, hasValue)
}

func (prs *Parser) parseWhileStatement() *syntaxnodes.WhileStatementNode {
    // consume 'while' keyword
    kw := prs.consume(lexer.TT_KW_While)

    // consume '('
    prs.consume(lexer.TT_OpenParenthesis)

    // parse condition
    condition := prs.parseExpression()

    // consume ')'
    prs.consume(lexer.TT_CloseParenthesis)

    // loop body
    body := prs.parseStatement()

    return syntaxnodes.NewWhileStatementNode(kw, condition, body)
}

func (prs *Parser) parseFromToStatement() *syntaxnodes.FromToStatementNode {
    // consume 'from' keyword
    kw := prs.consume(lexer.TT_KW_From)

    // consume iterator variable
    it := prs.consume(lexer.TT_Identifier)

    // consume '<-'
    prs.consume(lexer.TT_LeftArrow)

    // consume lower bound
    lb := prs.parseExpression()

    // consume 'to' keyword
    prs.consume(lexer.TT_KW_To)

    // consume upper bound
    up := prs.parseExpression()

    // consume loop body
    body := prs.parseStatement()

    return syntaxnodes.NewFromToStatementNode(kw, lb, it, up, body)
}

func (prs *Parser) parseForStatement() *syntaxnodes.ForStatementNode {
    // consume 'for' keyword
    kw := prs.consume(lexer.TT_KW_For)

    // consume '('
    prs.consume(lexer.TT_OpenParenthesis)

    init := prs.parseDeclarationStatement()
    prs.consume(lexer.TT_Semicolon)
    cond := prs.parseExpression()
    prs.consume(lexer.TT_Semicolon)
    action := prs.parseStatement()

    // consume ')'
    prs.consume(lexer.TT_CloseParenthesis)

    // parse loop body
    body := prs.parseStatement()

    return syntaxnodes.NewForStatementNode(kw, init, cond, action, body) 
}

func (prs *Parser) parseLoopStatement() *syntaxnodes.LoopStatementNode {
    // consume 'loop' keyword
    kw := prs.consume(lexer.TT_KW_Loop)

    // consume '('
    prs.consume(lexer.TT_OpenParenthesis)

    amount := prs.parseExpression()

    // consume ')'
    prs.consume(lexer.TT_CloseParenthesis)

    // loop body
    body := prs.parseStatement()

    return syntaxnodes.NewLoopStatementNode(kw, amount, body)
}

func (prs *Parser) parseBreakStatement() *syntaxnodes.BreakStatementNode {
    // consume 'break' keyword
    kw := prs.consume(lexer.TT_KW_Break)

    return syntaxnodes.NewBreakStatementNode(kw)
}

func (prs *Parser) parseContinueStatement() *syntaxnodes.ContinueStatementNode {
    // consume 'continue' keyword
    kw := prs.consume(lexer.TT_KW_Continue)

    return syntaxnodes.NewContinueStatementNode(kw)
}

func (prs *Parser) parseIfStatement() *syntaxnodes.IfStatementNode {
    // consume 'if' keyword
    kw := prs.consume(lexer.TT_KW_If)

    // consume '('
    prs.consume(lexer.TT_OpenParenthesis)

    // parse condition
    cond := prs.parseExpression()

    // consume ')'
    prs.consume(lexer.TT_CloseParenthesis)

    // parse if body 
    body := prs.parseStatement()

    // (optional) parse else statement
    var elseBody syntaxnodes.StatementNode
    hasElse := false

    if prs.current().Type == lexer.TT_KW_Else {
        // consume else kw
        prs.consume(lexer.TT_KW_Else)

        elseBody = prs.parseStatement()
        hasElse = true
    }

    return syntaxnodes.NewIfStatementNode(kw, cond, body, elseBody, hasElse)
}

func (prs *Parser) parseBlockStatement() *syntaxnodes.BlockStatementNode {
    // consume '{'
    op := prs.consume(lexer.TT_OpenBraces)

    var stmts []syntaxnodes.StatementNode
    for prs.current().Type != lexer.TT_CloseBraces {
        // parse statements
        stmts = append(stmts, prs.parseStatement())
    }

    // consume '}'
    cl := prs.consume(lexer.TT_CloseBraces)
    
    return syntaxnodes.NewBlockStatementNode(op, stmts, cl)
}

func (prs *Parser) parseExpressionStatement() *syntaxnodes.ExpressionStatementNode {
    // consume expression
    expr := prs.parseExpression()

    return syntaxnodes.NewExpressionStatementNode(expr)
}

// --------------------------------------------------------
// Expressions
// --------------------------------------------------------
func (prs *Parser) parseExpression() syntaxnodes.ExpressionNode {
    return prs.parseBinaryExpression(0)
}

func (prs *Parser) parseBinaryExpression(lastPrecedence int) syntaxnodes.ExpressionNode {
    var left syntaxnodes.ExpressionNode

    // is this a unary expression?
    unaryPrecendence := syntaxnodes.GetUnaryOperatorPrecedence(prs.current().Type)
    if unaryPrecendence != 0 && unaryPrecendence > lastPrecedence {
        // consume the operator
        op := prs.consume(prs.current().Type) 
        operand := prs.parseBinaryExpression(unaryPrecendence)

        // create new unary node
        return syntaxnodes.NewUnaryExpressionNode(operand, op)
        
    // otherwise: parse left side
    } else {
        left = prs.parsePrimaryExpression()
    }

    for {
        precedence := syntaxnodes.GetBinaryOperatorPrecedence(prs.current().Type)

        // if this isnt an operator or has less precedence
        // -> stop / hand control back over to parent
        if precedence == 0 || precedence <= lastPrecedence {
            break
        }

        operator := prs.consume(prs.current().Type)
        right := prs.parseBinaryExpression(precedence)

        // create binary expression and make it our new left
        left = syntaxnodes.NewBinaryExpressionNode(left, right, operator)
    }

    return left
}

func (prs *Parser) parsePrimaryExpression() syntaxnodes.ExpressionNode {
    // Literals
    if prs.current().Type == lexer.TT_String  || 
       prs.current().Type == lexer.TT_Integer ||
       prs.current().Type == lexer.TT_Float   ||
       prs.current().Type == lexer.TT_KW_True ||
       prs.current().Type == lexer.TT_KW_False {

        return prs.parseLiteralExpression()
    
    // Name, assignment, or call expression   
    } else if prs.current().Type == lexer.TT_Identifier {
        // Assignment expression
        if prs.peek(1).Type == lexer.TT_LeftArrow {
            return prs.parseAssignmentExpression()
        
        // Call expression
        } else if prs.peek(1).Type == lexer.TT_OpenParenthesis {
            return prs.parseCallExpression()

        // Name expression
        } else {
            return prs.parseNameExpression()
        }

    } else if prs.current().Type == lexer.TT_OpenParenthesis {
        return prs.parseParenthesizedExpression()

    // Dude i have no idea
    } else {
        error.Report(error.NewError(error.PRS, prs.current().Position, "Expected expression, got '%s'!", prs.current().Type))
        prs.step(1)

        return syntaxnodes.NewErrorExpressionNode(prs.current().Position)
    }
}

func (prs *Parser) parseLiteralExpression() *syntaxnodes.LiteralExpressionNode {
    // consume the literal
    lit := prs.consume(prs.current().Type)
    // create a new node
    return syntaxnodes.NewLiteralExpressionNode(lit)
}

func (prs *Parser) parseAssignmentExpression() *syntaxnodes.AssignmentExpressionNode {
    // consume the variable name
    id := prs.consume(lexer.TT_Identifier)

    // consume '<-'
    prs.consume(lexer.TT_LeftArrow)

    // parse assignment value
    val := prs.parseExpression()

    // create new node
    return syntaxnodes.NewAssignmentExpressionNode(id, val)
}

func (prs *Parser) parseCallExpression() *syntaxnodes.CallExpressionNode {
    // consume call expression
    id := prs.consume(lexer.TT_Identifier)

    // consume '('
    prs.consume(lexer.TT_OpenParenthesis)
    
    // arguments
    args := make([]syntaxnodes.ExpressionNode, 0)
    for prs.current().Type != lexer.TT_CloseParenthesis {
        // parse arg
        args = append(args, prs.parseExpression())

        if prs.current().Type == lexer.TT_Comma {
            prs.consume(lexer.TT_Comma)
        } else {
            break
        }
    }

    // consume ')'
    cprm := prs.consume(lexer.TT_CloseParenthesis)

    // create new node
    return syntaxnodes.NewCallExpressionNode(id, args, cprm)
}

func (prs *Parser) parseNameExpression() *syntaxnodes.NameExpressionNode {
    // consume name
    id := prs.consume(lexer.TT_Identifier)

    // create new node
    return syntaxnodes.NewNameExpressionNode(id)
}

func (prs *Parser) parseParenthesizedExpression() *syntaxnodes.ParenthesizedExpressionNode {
    // consume leading parenthsis
    start := prs.consume(lexer.TT_OpenParenthesis)

    // consume expression
    expr := prs.parseExpression()

    // consume trailing parenthesis
    end := prs.consume(lexer.TT_CloseParenthesis)

    // create a new node
    return syntaxnodes.NewParenthesizedExpressionNode(start, expr, end)
}
