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
       
        //prs.step(1)

        // fabricate a token of this kind to keep the compilation going
        return lexer.Token {
           Type: typ,
        }
    }

    // if it does match
    prs.step(1) // next token
    return prs.peek(-1) // return the one we were just at
}

func (prs *Parser) consumeWord(word string) lexer.Token {
    // get an id token
    id := prs.consume(lexer.TT_Identifier)

    // make sure it matches our word
    if id.Buffer != word {
        // report this error
        error.Report(error.NewError(error.PRS, prs.current().Position, "Expected keyword '%s', instead got: '%s'!", word, id.Buffer))
       
        // rewind
        prs.step(-1)

        // fabricate a token of this kind to keep the compilation going
        return lexer.Token {
           Type: lexer.TT_Identifier,
           Buffer: word,
        }
    }

    // cool beans
    return id
}

// real time travel shit
func (prs *Parser) rewind(tok lexer.Token) {
    for prs.current().Type != tok.Type &&
        prs.current().Position != tok.Position {

        prs.step(-1)
    }
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

    } else if prs.current().Type == lexer.TT_KW_Container {
        mem = prs.parseContainerMember()

    // anything else -> error
    } else {
        error.Report(error.NewError(error.PRS, prs.current().Position, "Expected member, instead got: '%s'!", prs.current().Type))
        prs.step(1)

        return
    }
   
    // if this isnt a function -> require a semicolon
    if mem.Type() != syntaxnodes.NT_Function && 
       mem.Type() != syntaxnodes.NT_Container {
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

    // consume function name OR constructor kw
    var id lexer.Token
    isConstructor := false

    if prs.current().Type == lexer.TT_KW_Constructor {
        id = prs.consume(lexer.TT_KW_Constructor)
        isConstructor = true
    } else {
        id = prs.consume(lexer.TT_Identifier)
    }

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

    return syntaxnodes.NewFunctionNode(kw, id, isConstructor, params, retType, hasReturnType, body)
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

func (prs *Parser) parseContainerMember() *syntaxnodes.ContainerNode {
    // consume 'container' keyword
    kw := prs.consume(lexer.TT_KW_Container)

    // consume container name 
    id := prs.consume(lexer.TT_Identifier)

    // consume '{'
    prs.consume(lexer.TT_OpenBraces)

    // consume as many members as we can
    fields := []*syntaxnodes.FieldClauseNode{}
    methods := []*syntaxnodes.FunctionNode{}
    for prs.current().Type != lexer.TT_CloseBraces && 
        prs.current().Type != lexer.TT_EOF {
        
        // is this a method?
        if prs.current().Type == lexer.TT_KW_Function {
            methods = append(methods, prs.parseFunctionMember())

        // if not -> probably a field lol
        } else {
            fields = append(fields, prs.parseFieldClause())
        }
    }

    // consume '}'
    cls := prs.consume(lexer.TT_CloseBraces)

    return syntaxnodes.NewContainerNode(kw, id, fields, methods, cls)
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

func (prs *Parser) parseFieldClause() *syntaxnodes.FieldClauseNode {
    // consume param name 
    id := prs.consume(lexer.TT_Identifier)

    // consume parm type
    typ := prs.parseTypeClause()

    // consume a semicolon
    prs.consume(lexer.TT_Semicolon)

    return syntaxnodes.NewFieldClauseNode(id, typ)
}

func (prs *Parser) parseTypeClause() *syntaxnodes.TypeClauseNode {
    var pack lexer.Token
    hasPackage := false

    // is there a package prefix?
    if prs.peek(1).Type == lexer.TT_Package {
        // consume the package name
        pack = prs.consume(lexer.TT_Identifier)

        // consume the ::
        prs.consume(lexer.TT_Package)
    }

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
    return syntaxnodes.NewTypeClauseNode(pack, hasPackage, id, subtypes)
}

func (prs *Parser) parseFieldAssignmentClause() *syntaxnodes.FieldAssignmentClauseNode {
    // consume the field name
    id := prs.consume(lexer.TT_Identifier)

    // consume the left arrow
    prs.consume(lexer.TT_LeftArrow)

    // consume the value
    val := prs.parseExpression()

    // ok cool
    return syntaxnodes.NewFieldAssignmentClauseNode(id, val)
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

        for prs.current().Type == lexer.TT_OpenBrackets ||
            prs.current().Type == lexer.TT_LeftArrow    ||
            prs.current().Type == lexer.TT_RightArrow   {

            // Is this actually an array index?
            if prs.current().Type == lexer.TT_OpenBrackets {
                left = prs.parseArrayIndexExpression(left)
            }

            // Is this actually an assignment?
            if prs.current().Type == lexer.TT_LeftArrow {
                left = prs.parseAssignmentExpression(left)
            }

            // Is this actually an access?
            if prs.current().Type == lexer.TT_RightArrow {
                left = prs.parseAccessExpression(left)
            }
        }
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
        // Call expression
        if prs.peek(1).Type == lexer.TT_OpenParenthesis ||
           prs.peek(1).Type == lexer.TT_Package {
            return prs.parseCallExpression()

        // Name expression
        } else {
            return prs.parseNameExpression()
        }

    // Parenthesized expressions
    } else if prs.current().Type == lexer.TT_OpenParenthesis {
        return prs.parseParenthesizedExpression()

    // Array creation
    } else if prs.current().Type == lexer.TT_KW_Make {
        return prs.parseMakeExpression()

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

func (prs *Parser) parseAssignmentExpression(expr syntaxnodes.ExpressionNode) *syntaxnodes.AssignmentExpressionNode {
    // consume '<-'
    prs.consume(lexer.TT_LeftArrow)

    // parse assignment value
    val := prs.parseExpression()

    // create new node
    return syntaxnodes.NewAssignmentExpressionNode(expr, val)
}

func (prs *Parser) parseCallExpression() *syntaxnodes.CallExpressionNode {
    var pack lexer.Token
    hasPackage := false

    if prs.peek(1).Type == lexer.TT_Package {
        pack = prs.consume(lexer.TT_Identifier)
        prs.consume(lexer.TT_Package)

        hasPackage = true
    }

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
    return syntaxnodes.NewCallExpressionNode(id, pack, hasPackage, args, cprm)
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

func (prs *Parser) parseMakeExpression() syntaxnodes.ExpressionNode {
    // consume the make kw
    kw := prs.consume(lexer.TT_KW_Make)

    // consume a package name if there is one
    var pack lexer.Token
    hasPack := false

    // if theres a :: token
    if prs.peek(1).Type == lexer.TT_Package {
        pack = prs.consume(lexer.TT_Identifier)
        prs.consume(lexer.TT_Package)
        hasPack = true
    }

    // consume the container name
    id := prs.consume(lexer.TT_Identifier)
    closing := id

    // is this next token an identifier?
    if prs.current().Type == lexer.TT_Identifier {
        if prs.current().Buffer == "array" {
            // aw man we fucked up
            // this is actually an array creation

            // use time travel to fix this
            prs.rewind(kw)
            return prs.parseMakeArrayExpression()
        }
    }

    // parse the initializer / constructor / nothing
    initializer := []*syntaxnodes.FieldAssignmentClauseNode{}
    hasInitializer := false

    args := []syntaxnodes.ExpressionNode{}
    hasConstructor := false

    // Constructor creation
    // make <container> (<args>);
    if prs.current().Type == lexer.TT_OpenParenthesis {
        prs.consume(lexer.TT_OpenParenthesis)

        // parse some args, yo
        for prs.current().Type != lexer.TT_CloseParenthesis {
            args = append(args, prs.parseExpression())

            // consume a comma, if there is one
            if prs.current().Type == lexer.TT_Comma {
                prs.consume(lexer.TT_Comma)

            // if theres none, we're probably done
            } else {
                break
            }
        }

        closing = prs.consume(lexer.TT_CloseParenthesis)
        hasConstructor = true

    // Initializer creation
    // make <container> { <fld> <- <val> }
    } else if prs.current().Type == lexer.TT_OpenBraces {
        prs.consume(lexer.TT_OpenBraces)

        // parse some field assignments, yo
        for prs.current().Type != lexer.TT_CloseBraces {
            initializer = append(initializer, prs.parseFieldAssignmentClause())

            // consume a comma, if there is one
            if prs.current().Type == lexer.TT_Comma {
                prs.consume(lexer.TT_Comma)

            // if theres none, we're probably done
            } else {
                break
            }
        } 

        closing = prs.consume(lexer.TT_CloseBraces)
        hasInitializer = true
    } 

    return syntaxnodes.NewMakeExpressionNode(kw, closing, id, pack, hasPack, initializer, hasInitializer, args, hasConstructor)
}

func (prs *Parser) parseMakeArrayExpression() *syntaxnodes.MakeArrayExpressionNode {
    // consume the make kw
    kw := prs.consume(lexer.TT_KW_Make)

    // consume the datatype
    typ := prs.parseTypeClause()

    // consume 'array' word
    prs.consumeWord("array")
    
    // consume either length or a literal
    var length syntaxnodes.ExpressionNode
    var initializer []syntaxnodes.ExpressionNode
    hasInitializer := false

    var closing lexer.Token

    // We got a literal
    if prs.current().Type == lexer.TT_OpenBraces {
        // consume {
        prs.consume(lexer.TT_OpenBraces)

        for prs.current().Type != lexer.TT_CloseBraces &&
            prs.current().Type != lexer.TT_EOF {
    
            initializer = append(initializer, prs.parseExpression())

            // Require a comma after every entry
            if prs.current().Type == lexer.TT_Comma {
                prs.consume(lexer.TT_Comma)
            } else {
                break
            }
        }

        // consume }
        closing = prs.consume(lexer.TT_CloseBraces)

        hasInitializer = true

    // We got a classic length generation
    } else {
        // consume (
        prs.consume(lexer.TT_OpenParenthesis)

        // consume length
        length = prs.parseExpression()

        // consume )
        closing = prs.consume(lexer.TT_CloseParenthesis)
    }

    return syntaxnodes.NewMakeArrayExpressionNode(kw, closing, typ, length, initializer, hasInitializer)
}

func (prs *Parser) parseArrayIndexExpression(expr syntaxnodes.ExpressionNode) *syntaxnodes.ArrayIndexExpressionNode {
    // consume [
    prs.consume(lexer.TT_OpenBrackets)

    // parse index
    idx := prs.parseExpression()

    // consume ]
    prs.consume(lexer.TT_CloseBrackets)

    return syntaxnodes.NewArrayIndexExpressionNode(expr, idx)
}

func (prs *Parser) parseAccessExpression(expr syntaxnodes.ExpressionNode) *syntaxnodes.AccessExpressionNode {
    // consume ->
    prs.consume(lexer.TT_RightArrow)

    // consume the field / method name
    id := prs.consume(lexer.TT_Identifier)
    cls := id

    // if this is a call -> parse some args
    args := []syntaxnodes.ExpressionNode{}
    isCall := false

    if prs.current().Type == lexer.TT_OpenParenthesis {
        // consume (
        prs.consume(lexer.TT_OpenParenthesis)

        // parse the arguments
        for prs.current().Type != lexer.TT_CloseParenthesis {
            args = append(args, prs.parseExpression())

            // consume a comma if we got one
            if prs.current().Type == lexer.TT_Comma {
                prs.consume(lexer.TT_Comma)
            } else {
                break
            }
        }

        // consume )
        cls = prs.consume(lexer.TT_CloseParenthesis)
        isCall = true
    }


    return syntaxnodes.NewAccessExpressionNode(expr, id, args, cls, isCall)
}
