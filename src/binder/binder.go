// Binder - binder.go
// --------------------------------------------------------
// The binder is a crucial part of the compilation process
// Its job is to do a semantic analysis on the parsed tree
// --------------------------------------------------------
package binder

import (
	"fmt"
	"slices"
	"strconv"

	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/lexer"
	packageprocessor "bytespace.network/rerect/package_processor"
	"bytespace.network/rerect/span"
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// --------------------------------------------------------
// Trait indexing
// --------------------------------------------------------
func IndexTraitTypes(file *packageprocessor.CompilationFile) {
    // look through all member nodes
    for _, v := range file.Members {
        // we're only looking for trait nodes
        if v.Type() != syntaxnodes.NT_Trait {
            continue
        }

        trtMem := v.(*syntaxnodes.TraitNode)

        // register a type symbol for this trait 
        typ := symbols.NewTypeSymbol(trtMem.TraitName.Buffer, []*symbols.TypeSymbol{}, symbols.TRT, 0, nil)

        // register a trait symbol for this trait 
        trt := symbols.NewTraitSymbol(file.Package, trtMem.TraitName.Buffer, typ)

        // register the type in the package
        ok := file.Package.TryRegisterTrait(trt) 

        if !ok {
            error.Report(error.NewError(error.BND, trtMem.TraitName.Position, "Cannot register trait '%s'! A symbol with that name already exists!", trt.TraitName))
            continue
        }

        file.Traits = append(file.Traits, trt)
        file.TraitSrc[trt] = trtMem
    }
}

func IndexTraitContents(file *packageprocessor.CompilationFile) {

    // work through all traits
    for _, trt := range file.Traits {
        // get the source node for this trait
        src := file.TraitSrc[trt]

        // create a new collection for the trait fields and methods
        fields := []*symbols.FieldSymbol{}
        meths := []*symbols.FunctionSymbol{}

        // bind all fields
        for _, v := range src.Fields {
            // resolve the field type
            typ := LookupTypeClause(v.FieldType, file.Package)

            // WAIT A MINUTE, DID WE HAVE A FIELD WITH THIS NAME ALREADY???
            if slices.Contains(trt.Symbols, v.FieldName.Buffer) {
                // jes -> DIE!!!! >:)
                error.Report(error.NewError(error.BND, v.Position(), "Cannot register field '%s'! A symbol with that name already exists!", v.FieldName.Buffer))
                continue
            }

            // nah, we good
            sym := symbols.NewTraitFieldSymbol(trt, v.FieldName.Buffer, typ)

            // add it to the list
            fields = append(fields, sym)
            trt.Symbols = append(trt.Symbols, sym.Name())
        }

        // bind all meths (methods, of course)
        for _, fncMem := range src.Methods {

            // is this a constructor?
            if fncMem.IsConstructor {
                // cringe
                error.Report(error.NewError(error.BND, fncMem.FunctionName.Position, "Illegal constructor outside of a container!"))
                continue
            }

            // create parameter symbols
            prms := []*symbols.ParameterSymbol{}
            for i, prm := range fncMem.Parameters {
                prms = append(prms, symbols.NewParameterSymbol(
                    prm.ParameterName.Buffer,
                    i,
                    LookupTypeClause(prm.ParameterType, file.Package),
                ))
            }

            ret := LookupTypeClause(fncMem.ReturnType, file.Package)

            // register a function symbol for this method
            fnc := symbols.NewMethodSymbol(
                file.Package,
                trt.TraitType,
                fncMem.FunctionName.Buffer,
                ret,
                prms,
            )

            // if this is just a declaration -> mark this is needing to be called virtually
            if !fncMem.HasBody {
                fnc.NeedsVirtualCallToContainer = true
            }

            // okay but like, is this legal?
            if slices.Contains(trt.Symbols, fnc.Name()) {
                error.Report(error.NewError(error.BND, fncMem.FunctionName.Position, "Cannot register method '%s'! A symbol with that name already exists!", fnc.Name()))
                continue
            }

            // if it is -> register the name in this container
            trt.Symbols = append(trt.Symbols, fnc.FuncName)
            meths = append(meths, fnc)

            // register it globally in this package
            file.Package.TryRegisterFunction(fnc) 

            // if this is a declaration -> do not register is as an actual function
            if !fncMem.HasBody {
                continue
            }

            // otherwise -> register method as a global function (because im lazy and they are treated equal anyways :) ) 
            file.Functions = append(file.Functions, fnc)
            file.FunctionBodiesSrc[fnc] = fncMem.Body
        }

        // store the container contents in the container
        trt.Fields = fields
        trt.Methods = meths
    }
}

// --------------------------------------------------------
// Container indexing
// --------------------------------------------------------
func IndexContainerTypes(file *packageprocessor.CompilationFile) {
    // look through all member nodes
    for _, v := range file.Members {
        // we're only looking for container nodes
        if v.Type() != syntaxnodes.NT_Container {
            continue
        }

        cntMem := v.(*syntaxnodes.ContainerNode)

        // register a type symbol for this container
        typ := symbols.NewTypeSymbol(cntMem.ContainerName.Buffer, []*symbols.TypeSymbol{}, symbols.CONT, 0, nil)

        // register a container symbol for this container 
        cnt := symbols.NewContainerSymbol(file.Package, cntMem.ContainerName.Buffer, typ)

        // now: look up the traits we got
        for _, v := range cntMem.Traits {
            // look the trait up

            var trt *symbols.TraitSymbol

            if v.HasPackage {
                // look up the package
                pack := LookupPackageInPackage(file.Package, v.Package.Buffer)

                // aw man no package
                if pack == nil {
                    error.Report(error.NewError(error.BND, v.Package.Position, "Could not find a package called '%s'", v.Package.Buffer))
                    continue

                // aw man a package
                } else {
                    trt = LookupTraitInPackage(v.TraitName.Buffer, pack)

                    if trt == nil {
                        error.Report(error.NewError(error.BND, v.Position(), "Could not find a trait called '%s' in package '%s'!", v.TraitName.Buffer, v.Package.Buffer))
                        continue
                    }
                } 

            // if theres no package specified -> try in the current one
            } else {
                trt = LookupTrait(v.TraitName.Buffer, file.Package)

                if trt == nil {
                    error.Report(error.NewError(error.BND, v.Position(), "Could not find a trait called '%s'!", v.TraitName.Buffer))
                    continue
                }
            }

            // if we got to here then we found a trait
            cnt.Traits = append(cnt.Traits, trt)
        }

        // register the type in the package
        ok := file.Package.TryRegisterContainer(cnt) 

        if !ok {
            error.Report(error.NewError(error.BND, cntMem.ContainerName.Position, "Cannot register container '%s'! A symbol with that name already exists!", cnt.ContainerName))
            continue
        }

        file.Containers = append(file.Containers, cnt)
        file.ContainerSrc[cnt] = cntMem
    }
}

func IndexContainerContents(file *packageprocessor.CompilationFile) {

    // work through all containers
    for _, cnt := range file.Containers {
        // get the original source node of this container
        src := file.ContainerSrc[cnt]

        // create a new fields collection and a flag for keeping track of constructors
        fields := []*symbols.FieldSymbol{}
        meths := []*symbols.FunctionSymbol{}
        hasConstructor := false

        // bind all fields
        for _, v := range src.Fields {
            // resolve the field type
            typ := LookupTypeClause(v.FieldType, file.Package)

            // WAIT A MINUTE, DID WE HAVE A FIELD WITH THIS NAME ALREADY???
            if slices.Contains(cnt.Symbols, v.FieldName.Buffer) {
                // jes -> DIE!!!! >:)
                error.Report(error.NewError(error.BND, v.Position(), "Cannot register field '%s'! A symbol with that name already exists!", v.FieldName.Buffer))
                continue
            }

            // nah, we good
            sym := symbols.NewFieldSymbol(cnt, v.FieldName.Buffer, typ)

            // add it to the list
            fields = append(fields, sym)
            cnt.Symbols = append(cnt.Symbols, sym.FieldName)
        }

        // bind all meths (methods, of course)
        for _, fncMem := range src.Methods {

            // Do we have multiple constructors?
            if fncMem.IsConstructor && hasConstructor {
                error.Report(error.NewError(error.BND, fncMem.FunctionName.Position, "Only one constructor per container is allowed!"))
                continue
            }

            // is this a function declaration?
            if !fncMem.HasBody {
                // cringe
                error.Report(error.NewError(error.BND, fncMem.Position(), "Illegal function declaration outside of a trait!"))
                continue
            }

            // create parameter symbols
            prms := []*symbols.ParameterSymbol{}
            for i, prm := range fncMem.Parameters {
                prms = append(prms, symbols.NewParameterSymbol(
                    prm.ParameterName.Buffer,
                    i,
                    LookupTypeClause(prm.ParameterType, file.Package),
                ))
            }

            ret := LookupTypeClause(fncMem.ReturnType, file.Package)

            // if this is a constructor -> we found one
            if fncMem.IsConstructor {
                // is this legal doe?
                if !ret.Equal(compunit.GlobalDataTypeRegister["void"]) {
                    error.Report(error.NewError(error.BND, fncMem.ReturnType.Position(), "Constructor is required to be of type void!"))
                    continue
                }

                hasConstructor = true
            }

            // register a function symbol for this method
            fnc := symbols.NewMethodSymbol(
                file.Package,
                cnt.ContainerType,
                fncMem.FunctionName.Buffer,
                ret,
                prms,
            )

            // okay but like, is this legal?
            if slices.Contains(cnt.Symbols, fnc.Name()) {
                error.Report(error.NewError(error.BND, fncMem.FunctionName.Position, "Cannot register method '%s'! A symbol with that name already exists!", fnc.Name()))
                continue
            }

            // if it is -> register it globally in this package
            file.Package.TryRegisterFunction(fnc) 

            // register method as a global function (because im lazy and they are treated equal anyways :) ) 
            file.Functions = append(file.Functions, fnc)
            file.FunctionBodiesSrc[fnc] = fncMem.Body

            // also register the name in this container
            cnt.Symbols = append(cnt.Symbols, fnc.FuncName)

            // also store it in our local meths array (for trait checking)
            meths = append(meths, fnc)

            // if its a constructor -> also register it in the container symbol
            if fncMem.IsConstructor {
                cnt.Constructor = fnc
            }
        }

        // If all of that was successful...

        // ------------------------------------------------
        // Apply all traits
        // ------------------------------------------------

        // First:
        // Implement all fields
        // --------------------
        for i, trt := range cnt.Traits {
            trtSrc := src.Traits[i]

            // implement fields one by one
            for _, fld := range trt.Fields {

                // does this container already have a field with this name?
                var containerField *symbols.FieldSymbol
                for _, f := range fields {
                    if f.FieldName == fld.FieldName {
                        containerField = f
                    }
                }

                // if so...
                if containerField != nil {

                    // do the datatypes match?
                    if !containerField.VarType().Equal(fld.VarType()) {

                        // was this field added by another trait? (this is just for nicer error messages)
                        if containerField.HasParentTrait {
                            error.Report(error.NewError(error.BND, trtSrc.Position(), "Unable to apply trait '%s'! A field with the name '%s' has already been added by trait '%s' with a different datatype!", trt.Name(), fld.Name(), containerField.ParentTrait.Name()))
                            continue

                        // otherwise: a less complicated error message
                        } else {
                            error.Report(error.NewError(error.BND, trtSrc.Position(), "Unable to apply trait '%s'! The container '%s' already defines a field called '%s' with a different type!", trt.Name(), cnt.Name(), fld.Name()))
                            continue
                        }
                    }

                    // if the datatypes match -> everything is cool
                    // we dont need to add another field, the field required by this trait has already been added
                    continue
                }

                // ooootherwise -> we need to create this field

                // WAIT A MINUTE, DID WE HAVE A SYMBOL WITH THIS NAME ALREADY???
                if slices.Contains(cnt.Symbols, fld.Name()) {
                    // jes -> DIE!!!! >:)
                    error.Report(error.NewError(error.BND, trtSrc.Position(), "Cannot register field '%s' of trait '%s'! A symbol with that name already exists!", fld.Name(), trt.Name()))
                    continue
                }

                // nah, we good
                sym := symbols.NewFieldSymbol(cnt, fld.Name(), fld.FieldType)

                // add the trait in here (in case another trait also defines this field)
                sym.HasParentTrait = true
                sym.ParentTrait = trt

                // add it to the list
                fields = append(fields, sym)
                cnt.Symbols = append(cnt.Symbols, sym.FieldName)
            } 
        }

        // Secondly: 
        // Implement all pre-defined Methods
        // ---------------------------------
        // These are methods given by the trait, which already have an implementation.
        // Because of this the container is not allowed to define a function with the same name, even if its
        // contents are identical / equvilant.
        // These methods are standardised between trait implementers.
        // --------------------------------------------------------------------------------------------------

        for i, trt := range cnt.Traits {
            trtSrc := src.Traits[i]

            // include these methods one by one 
            for _, meth := range trt.Methods {
                // we only care for already implemented methods
                if meth.NeedsVirtualCallToContainer {
                    continue
                }

                // does this container already have a method with this name?
                isConflicting := false
                for _, f := range meths {
                    if f.FuncName == meth.FuncName {
                        // ILLEGAL!!!

                        // where did this method come from?
                        // did another trait add it? (this is just for more helpful error messages)
                        if f.SourceTrait != nil {
                            error.Report(error.NewError(error.BND, trtSrc.Position(), "Cannot add method '%s' of trait '%s'! A method with the same name has already been added by trait '%s'!", meth.Name(), trt.Name(), f.SourceTrait.Name()))
                        } else {
                            error.Report(error.NewError(error.BND, trtSrc.Position(), "Cannot add method '%s' of trait '%s'! The container already implements a method with that name!", meth.Name(), trt.Name()))
                        }

                        isConflicting = true
                        break
                    }
                }

                if isConflicting {
                    continue
                }

                // if everything is fine -> we need to import this method

                // WAIT A MINUTE, DID WE HAVE A SYMBOL WITH THIS NAME ALREADY???
                if slices.Contains(cnt.Symbols, meth.Name()) {
                    // yea :(
                    error.Report(error.NewError(error.BND, trtSrc.Position(), "Cannot register method '%s' of trait '%s'! A symbol with that name already exists!", meth.Name(), trt.Name()))
                    continue
                }

                // if everything is looking good -> import!!!!!
                // --------------------------------------------

                // create a copy of our original trait method symbol, but change the source type
                fnc := symbols.NewMethodSymbol(
                    meth.ParentPackage,
                    cnt.ContainerType,
                    meth.FuncName,
                    meth.ReturnType,
                    meth.Parameters,
                )

                // remember from which trait this method came (for better error reporting)
                fnc.SourceTrait = trt

                // mark this symbol as a redirection to another method and add a ref to that method
                fnc.NeedsVirtualCallToTrait = true
                fnc.TraitSourceMethod = meth

                // register it globally in this package
                file.Package.TryRegisterFunction(fnc) 

                // also register the name in this container
                cnt.Symbols = append(cnt.Symbols, fnc.FuncName)

                // also store it in our local meths array (for trait checking)
                meths = append(meths, fnc)
            }
        }

        // Third but not third: (it is third)
        // Make sure that all trait-declared methods have been implemented
        // (this is the easiest step because all it does is complain)

        for i, trt := range cnt.Traits {
            trtSrc := src.Traits[i]

            // include these methods one by one 
            for _, meth := range trt.Methods {
                // we only care about declarations
                if meth.NeedsVirtualCallToTrait {
                    continue
                }

                // did this container implement the declaration?
                var fnc *symbols.FunctionSymbol
                for _, f := range meths {
                    if f.FuncName == meth.FuncName {
                        // we found something
                        fnc = f
                        break
                    }
                }

                // if we did not find an implementation -> complain
                if fnc == nil {
                    error.Report(error.NewError(error.BND, trtSrc.Position(), "Container '%s' did not implement method '%s' which is required by trait '%s'!", cnt.Name(), meth.Name(), trt.Name()))
                    continue
                }

                // if we found a method with the correct name -> make sure the signatures match up
                // -------------------------------------------------------------------------------
                
                if !meth.ReturnType.Equal(fnc.ReturnType) {
                    error.Report(error.NewError(error.BND, trtSrc.Position(), "Container '%s' did not implement method '%s' correctly. Trait '%s' requires a return type of '%s',got '%s' instead!", cnt.Name(), meth.Name(), trt.Name(), meth.ReturnType.Name(), fnc.ReturnType.Name()))
                    continue
                }

                if len(meth.Parameters) != len(fnc.Parameters) {
                    error.Report(error.NewError(error.BND, trtSrc.Position(), "Container '%s' did not implement method '%s' correctly. Trait '%s' requires %d parameters, got %d instead!", cnt.Name(), meth.Name(), trt.Name(), len(meth.Parameters), len(fnc.Parameters)))
                    continue
                }

                for i := range meth.Parameters {
                    if !meth.Parameters[i].VarType().Equal(fnc.Parameters[i].VarType()) {
                        error.Report(error.NewError(error.BND, trtSrc.Position(), "Container '%s' did not implement method '%s' correctly. Trait '%s' requires the parameter at index %d to be of type '%s', got '%s' instead!", cnt.Name(), meth.Name(), trt.Name(), i, meth.Parameters[i].VarType().Name(), fnc.Parameters[i].VarType().Name()))
                        break
                    }
                }

                // otherwise we good
            }
        }



        // store the container contents in the container
        cnt.Fields = fields
        cnt.Methods = meths
    }
}

// --------------------------------------------------------
// Function indexing
// --------------------------------------------------------
func IndexFunctions(file *packageprocessor.CompilationFile) {
    // look through all member nodes
    for _, v := range file.Members {
        // we're only looking for function nodes
        if v.Type() != syntaxnodes.NT_Function {
            continue
        }

        fncMem := v.(*syntaxnodes.FunctionNode)

        // is this a constructor?
        if fncMem.IsConstructor {
            // cringe
            error.Report(error.NewError(error.BND, fncMem.FunctionName.Position, "Illegal constructor outside of a container!"))
            continue
        }

        // is this a function declaration?
        if !fncMem.HasBody {
            // cringe
            error.Report(error.NewError(error.BND, fncMem.Position(), "Illegal function declaration outside of a trait!"))
            continue
        }

        // create parameter symbols
        prms := []*symbols.ParameterSymbol{}
        for i, prm := range fncMem.Parameters {
            prms = append(prms, symbols.NewParameterSymbol(
                prm.ParameterName.Buffer,
                i,
                LookupTypeClause(prm.ParameterType, file.Package),
            ))
        }

        // register a function symbol for this function
        fnc := symbols.NewFunctionSymbol(
            file.Package,
            fncMem.FunctionName.Buffer,
            LookupTypeClause(fncMem.ReturnType, file.Package),
            prms,
        )

        ok := file.Package.TryRegisterFunction(fnc) 

        if !ok {
            error.Report(error.NewError(error.BND, fncMem.FunctionName.Position, "Cannot register function '%s'! A function with that name already exists!", fnc.FuncName))
            continue
        }

        file.Functions = append(file.Functions, fnc)
        file.FunctionBodiesSrc[fnc] = fncMem.Body
    }
}

// --------------------------------------------------------
// Global indexing
// --------------------------------------------------------
func IndexGlobals(file *packageprocessor.CompilationFile) {
    // look through all member nodes
    for _, v := range file.Members {
        // we're only looking for function nodes
        if v.Type() != syntaxnodes.NT_Global {
            continue
        }

        glbMem := v.(*syntaxnodes.GlobalNode)

        // register a global symbol for this function
        glb := symbols.NewGlobalSymbol(file.Package, glbMem.GlobalName.Buffer, LookupTypeClause(glbMem.VarType, file.Package))
        ok := file.Package.TryRegisterGlobal(glb) 

        if !ok {
            error.Report(error.NewError(error.BND, glbMem.GlobalName.Position, "Cannot register global '%s'! A global with that name already exists!", glb.Name()))
            continue
        }

        file.Globals = append(file.Globals, glb)
    }
}

// --------------------------------------------------------
// Binding
// --------------------------------------------------------
type Binder struct {
    CurrentPackage *symbols.PackageSymbol
    CurrentType *symbols.TypeSymbol
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

func BindFunctions(file *packageprocessor.CompilationFile) {
    for _, sym := range file.Functions {
        // create a new binder
        bin := Binder{
            CurrentPackage: file.Package,
            CurrentFunction: sym,
            CurrentScope: NewScope(nil),
        }

        // if this is a method -> register the current type
        if sym.FunctionKind == symbols.FT_METH {
            bin.CurrentType = sym.MethodSource

            // is this method part of a container or a trait?
            // ----------------------------------------------

            // also register all fields as variables
            if bin.CurrentType.TypeGroup == symbols.CONT {
                for _, v := range bin.CurrentType.Container.Fields {
                    bin.CurrentScope.RegisterVariable(v)
                } 

                // and register an instance variable ("this")
                bin.CurrentScope.RegisterVariable(symbols.NewInstanceSymbol(bin.CurrentType))

            } else if bin.CurrentType.TypeGroup == symbols.TRT {
                for _, v := range bin.CurrentType.Trait.Fields {
                    bin.CurrentScope.RegisterVariable(v)
                } 

                // and register an instance variable ("this")
                bin.CurrentScope.RegisterVariable(symbols.NewInstanceSymbol(bin.CurrentType))
            }

        }

        // register the package globals as variables
        for _, v := range file.Globals {
            bin.CurrentScope.RegisterVariable(v)
        }

        // create sub-scope so globals can be overwritten
        bin.EnterNewScope()

        // register the function parameters as variables
        for _, v := range sym.Parameters {
            bin.CurrentScope.RegisterVariable(v)
        }


        file.FunctionBodies[sym] = bin.bindStatement(file.FunctionBodiesSrc[sym])
    }
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
        typ = LookupTypeClause(stmt.VarType, bin.CurrentPackage)
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
    if expr.Type() != boundnodes.BT_CallExpr       && 
       expr.Type() != boundnodes.BT_AccessCallExpr &&
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

    } else if expr.Type() == syntaxnodes.NT_MakeArrayExpr {
        return bin.bindMakeArrayExpression(expr.(*syntaxnodes.MakeArrayExpressionNode))
    
    } else if expr.Type() == syntaxnodes.NT_ArrayIndexExpr {
        return bin.bindArrayIndexExpression(expr.(*syntaxnodes.ArrayIndexExpressionNode))

    } else if expr.Type() == syntaxnodes.NT_AccessExpr {
        return bin.bindAccessExpression(expr.(*syntaxnodes.AccessExpressionNode))

    } else if expr.Type() == syntaxnodes.NT_MakeExpr {
        return bin.bindMakeExpression(expr.(*syntaxnodes.MakeExpressionNode))

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
        val, err := strconv.ParseInt(expr.Literal.Buffer, 10, 32) 
        
        if err != nil {
            error.Report(error.NewError(error.BND, expr.Position(), "Could not convert '%s' to an integer!", expr.Literal.Buffer))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        value = int32(val)
        typ = compunit.GlobalDataTypeRegister["int"]

    } else if expr.Literal.Type == lexer.TT_Float {
        val, err := strconv.ParseFloat(expr.Literal.Buffer, 32)
        
        if err != nil {
            error.Report(error.NewError(error.BND, expr.Position(), "Could not convert '%s' to a float!", expr.Literal.Buffer))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        value = float32(val)
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
    // bind the source expression
    exp := bin.bindExpression(expr.Expression)

    // make sure we're allowed to assign to this type of expression
    if exp.Type() != boundnodes.BT_NameExpr       &&
       exp.Type() != boundnodes.BT_ArrayIndexExpr &&
       exp.Type() != boundnodes.BT_AccessFieldExpr {

        error.Report(error.NewError(error.BND, expr.Position(), "Cannot assign to expression of type '%s'!", expr.Expression.Type()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // bind assignment value
    val := bin.bindExpression(expr.Value)

    // make sure the data types match
    val = bin.bindConversion(val, exp.ExprType(), false)

    // cool
    return boundnodes.NewBoundAssignmentExpressionNode(expr, exp, val)
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
    if !expr.HasPackage && len(expr.Parameters) == 1 {
        // are we calling a type name?
        typ := LookupType(expr.Identifier.Buffer, expr.Identifier.Position, bin.CurrentPackage, true)
        
        // if so -> bind a conversion
        if typ != nil {
            exp := bin.bindExpression(expr.Parameters[0])
            return bin.bindConversion(exp, typ, true)
        }
    }

    // is this actually a cast but to a type from a different package?
    if expr.HasPackage && len(expr.Parameters) == 1 {
        // see if this is actually a container
        cnt := bin.LookupContainerInPackage(expr.Package.Buffer, expr.Identifier.Buffer)

        // if we got something -> bind a conversion
        if cnt != nil {
            exp := bin.bindExpression(expr.Parameters[0])
            return bin.bindConversion(exp, cnt.ContainerType, true)
        }
    }

    // otherwise -> bind a call
    // ------------------------

    // lookup the function
    var fnc *symbols.FunctionSymbol
    if expr.HasPackage {
        fnc = bin.LookupFunctionInPackage(expr.Package.Buffer, expr.Identifier.Buffer)
    } else {
        fnc = bin.LookupFunction(expr.Identifier.Buffer)
    }

    if fnc == nil {
        if !expr.HasPackage {
            error.Report(error.NewError(error.BND, expr.Identifier.Position, "Could not find function '%s'!", expr.Identifier.Buffer))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        } else {
            error.Report(error.NewError(error.BND, expr.Identifier.Position.SpanBetween(expr.Package.Position), "Could not find function '%s::%s'!", expr.Package.Buffer, expr.Identifier.Buffer))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }
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
        args[i] = bin.bindConversion(args[i], fnc.Parameters[i].VarType(), false)
    }

    // ok cool
    return boundnodes.NewBoundCallExpressionNode(expr, fnc, args)
}

func (bin *Binder) bindNameExpression(expr *syntaxnodes.NameExpressionNode) boundnodes.BoundExpressionNode {
    // is this a global name expression?
    if expr.HasPackage {
        // only globals in this package are accessible
        if bin.CurrentPackage.Name() != expr.PackageName.Buffer {
            error.Report(error.NewError(error.BND, expr.Position(), "Unable to resolve global '%s' in package '%s': only globals in the current package ('%s') are allowed to be accessed!", expr.Identifier.Buffer, expr.PackageName.Buffer, bin.CurrentPackage.Name()))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        // if we *are* in this package
        // look up the global
        glb := bin.LookupGlobal(expr.Identifier.Buffer)

        // did we find one?
        if glb == nil {
            error.Report(error.NewError(error.BND, expr.Position(), "Could not find global called '%s'!", expr.Identifier.Buffer))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        // ok cool
        return boundnodes.NewBoundNameExpressionNode(expr, glb)
    }

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

func (bin *Binder) bindMakeArrayExpression(expr *syntaxnodes.MakeArrayExpressionNode) boundnodes.BoundExpressionNode {
    // resolve the array type
    typ := LookupTypeClause(expr.ArrType, bin.CurrentPackage)

    // create an array type for it
    arrtyp := createArrayType(typ)

    // bind either the length or the initializer
    var length boundnodes.BoundExpressionNode
    var initializer []boundnodes.BoundExpressionNode

    if !expr.HasInitializers {
        // bind the length expression
        length = bin.bindExpression(expr.Length)
        
        // make sure its an int
        length = bin.bindConversion(length, compunit.GlobalDataTypeRegister["int"], false)

    } else {
        // bind each element of the initializer
        for _, v := range expr.Initializers {
            entry := bin.bindExpression(v)

            // make sure the types match
            entry = bin.bindConversion(entry, typ, false)

            // add it to the list
            initializer = append(initializer, entry)
        }
    }

    return boundnodes.NewBoundMakeArrayExpressionNode(expr, arrtyp, length, initializer, expr.HasInitializers)
}

func (bin *Binder) bindArrayIndexExpression(expr *syntaxnodes.ArrayIndexExpressionNode) boundnodes.BoundExpressionNode {
    // bind the source of the array
    src := bin.bindExpression(expr.Expression)

    // make sure the src is an array
    if src.ExprType().TypeGroup != symbols.ARR {
        error.Report(error.NewError(error.BND, expr.Expression.Position(), "Indexing is only allowed on array types, got '%s'!", src.ExprType().Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // bind the index
    idx := bin.bindExpression(expr.Index)

    // make sure the index is an int
    idx = bin.bindConversion(idx, compunit.GlobalDataTypeRegister["int"], false)

    // ok cool
    return boundnodes.NewBoundArrayIndexExpressionNode(expr, src, idx)
}

func (bin *Binder) bindAccessExpression(expr *syntaxnodes.AccessExpressionNode) boundnodes.BoundExpressionNode {
    if expr.IsCall {
        return bin.bindAccessCallExpression(expr)
    }

    // bind the source expression
    src := bin.bindExpression(expr.Expression)

    // damn but like, is this even a container?
    if src.ExprType().TypeGroup != symbols.CONT && src.ExprType().TypeGroup != symbols.TRT {
        error.Report(error.NewError(error.BND, expr.Identifier.Position, "Unable to access a field on non container or trait type '%s'!", src.ExprType().Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    var fld *symbols.FieldSymbol

    if src.ExprType().TypeGroup == symbols.TRT {
        // look up the field in a trait
        fld = LookupFieldInTrait(expr.Identifier.Buffer, src.ExprType().Trait) 
    } else {
        // look up the field in a container
        fld = LookupFieldInContainer(expr.Identifier.Buffer, src.ExprType().Container) 
    }

    // did we actually find something?
    if fld == nil {
        error.Report(error.NewError(error.BND, expr.Identifier.Position, "Did not find field '%s' in container or trait type '%s'!", expr.Identifier.Buffer, src.ExprType().Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // ok cool
    return boundnodes.NewBoundAccessFieldExpressionNode(expr, src, fld)
}

func (bin *Binder) bindAccessCallExpression(expr *syntaxnodes.AccessExpressionNode) boundnodes.BoundExpressionNode {
    // bind the source expression
    src := bin.bindExpression(expr.Expression)

    // lookup this method
    meth := bin.LookupMethod(expr.Identifier.Buffer, src.ExprType())

    // did we find something?
    if meth == nil {
        error.Report(error.NewError(error.BND, expr.Identifier.Position, "Could not find method '%s' for type '%s'!", expr.Identifier.Buffer, src.ExprType().Name()))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // was the right amount of arguments given?
    if len(meth.Parameters) != len(expr.Arguments) {
        error.Report(error.NewError(error.BND, expr.Position(), "Method '%s' expects %d arguments, got: %d!", meth.FuncName, len(meth.Parameters), len(expr.Arguments)))
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    // bind all args
    args := []boundnodes.BoundExpressionNode{}
    for _, v := range expr.Arguments {
        args = append(args, bin.bindExpression(v))
    }

    // make sure the datatypes match up
    for i := range meth.Parameters {
        args[i] = bin.bindConversion(args[i], meth.Parameters[i].VarType(), false)
    }

    // ok cool
    return boundnodes.NewBoundAccessCallExpressionNode(expr, src, meth, args)
}

func (bin *Binder) bindMakeExpression(expr *syntaxnodes.MakeExpressionNode) boundnodes.BoundExpressionNode {
    // look up the container we'll be instantiating
    var cnt *symbols.ContainerSymbol

    // if a package was given the lookup needs to be slightly different
    if expr.HasPackage {
        cnt = bin.LookupContainerInPackage(expr.Package.Buffer, expr.Container.Buffer)

    // otherwise, just do a normal lookup
    } else {
        cnt = LookupContainer(expr.Container.Buffer, bin.CurrentPackage)
    }

    // did we find it?
    if cnt == nil {
        error.Report(error.NewError(error.BND, expr.Container.Position, "Unable to find a container called '%s'!", expr.Container.Buffer)) 
        return boundnodes.NewBoundErrorExpressionNode(expr)
    }

    initializer := make(map[*symbols.FieldSymbol]boundnodes.BoundExpressionNode)
    args := []boundnodes.BoundExpressionNode{}

    // bind the initializer, if it exists
    if expr.HasInitializer {
        for _, v := range expr.Initializer {
            field := LookupFieldInContainer(v.FieldName.Buffer, cnt)

            if field == nil {
                error.Report(error.NewError(error.BND, v.FieldName.Position, "Unable to find field '%s' in container '%s'!", v.FieldName.Buffer, cnt.Name())) 
                return boundnodes.NewBoundErrorExpressionNode(expr)
            }

            val := bin.bindExpression(v.Value)

            // make sure the types match up
            val = bin.bindConversion(val, field.VarType(), false)

            // add it to the list
            initializer[field] = val
        }
    }

    // bind the constructor, if it exists
    if expr.HasConstructor {
        // does the container even have a constructor??
        if cnt.Constructor == nil {
            error.Report(error.NewError(error.BND, expr.Position(), "Unable to call constructor: container '%s' does not have a constructor!", cnt.Name()))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        // otherwise -> make sure the call is correct
        if len(cnt.Constructor.Parameters) != len(expr.ConstructorArguments) {
            error.Report(error.NewError(error.BND, expr.Position(), "Unable to call constructor: constructor for container '%s' expects %d arguments, got %d!", cnt.Name(), len(cnt.Constructor.Parameters), len(expr.ConstructorArguments)))
            return boundnodes.NewBoundErrorExpressionNode(expr)
        }

        // bind all args, make sure the types match up
        for i, v := range expr.ConstructorArguments {
            val := bin.bindExpression(v)
            val = bin.bindConversion(val, cnt.Constructor.Parameters[i].VarType(), false)
            args = append(args, val)
        }
    }

    // create the node
    return boundnodes.NewBoundMakeExpressionNode(expr, cnt, initializer, expr.HasInitializer, args, expr.HasConstructor)
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
    
    // identity conversion -> just return original value
    if con == boundnodes.CT_Identity {
        return expr
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
func LookupType(name string, pos span.Span, pck *symbols.PackageSymbol, canfail bool) *symbols.TypeSymbol {
    // lookup primitives
    typ, ok := compunit.GlobalDataTypeRegister[name]
    if ok {
        return typ
    }

    // lookup containers
    cnt := LookupContainer(name, pck)
    if cnt != nil {
        return cnt.ContainerType
    }

    // lookup traits
    trt := LookupTrait(name, pck)
    if trt != nil {
        return trt.TraitType
    }

    // if this allowed to fail -> do that
    if canfail {
        return nil
    }

    // otherwise -> DIE!!!!! (but like, gently, no crashing here :) ) 
    error.Report(error.NewError(error.BND, pos, "Unknown data type '%s'!", name))
    return compunit.GlobalDataTypeRegister["error"]
}

func LookupTypeClause(typ *syntaxnodes.TypeClauseNode, pack *symbols.PackageSymbol) *symbols.TypeSymbol {
   
    // if the type clause does not exists -> void return type
    if typ == nil {
        return compunit.GlobalDataTypeRegister["void"]
    }

    // if this is an array type, we will need to construct it
    if typ.TypeName.Buffer == "array" {
        // make sure we have exactly one subtype 
        if len(typ.SubTypes) != 1 {
            error.Report(error.NewError(error.BND, typ.Position(), "Data type '%s' takes exactly one subtype, got: %d!", typ.TypeName.Buffer, len(typ.SubTypes)))
            return compunit.GlobalDataTypeRegister["error"]
        }

        // if we do -> resolve it
        subtype := LookupTypeClause(typ.SubTypes[0], pack)

        // create a new type symbol
        arrsym := createArrayType(subtype) 
        return arrsym
    }

    // if the type clause has a package prefix -> this is def a container
    // (packages cant just contain random primitives)
    if typ.HasPackageName {
        // look up the package
        pck := LookupPackageInPackage(pack, typ.PackageName.Buffer)

        // did it work?
        if pck == nil {
            error.Report(error.NewError(error.BND, typ.PackageName.Position, "Could not find package '%s'!", typ.PackageName.Buffer))
            return compunit.GlobalDataTypeRegister["error"]
        }

        // try looking up a trait first
        trt := LookupTraitInPackage(typ.TypeName.Buffer, pck)
        if trt != nil {
            // we found something? great success!
            return trt.TraitType
        }

        // look up the container
        cnt := LookupContainerInPackage(typ.TypeName.Buffer, pck)

        // did it work?
        if cnt == nil {
            error.Report(error.NewError(error.BND, typ.TypeName.Position, "Could not find type '%s' in '%s'!", typ.TypeName.Buffer, typ.PackageName.Buffer))
            return compunit.GlobalDataTypeRegister["error"]
        }

        // ok cool
        return cnt.ContainerType
    }

    // otherwise -> look up the type
    return LookupType(typ.TypeName.Buffer, typ.Position(), pack, false)
}

func createArrayType(subtype *symbols.TypeSymbol) *symbols.TypeSymbol {
    return symbols.NewTypeSymbol(subtype.Name() + " Array", []*symbols.TypeSymbol{subtype}, symbols.ARR, 0, nil)
}

// --------------------------------------------------------
// Trait Lookup
// --------------------------------------------------------
func LookupTrait(name string, pck *symbols.PackageSymbol) *symbols.TraitSymbol {
    trt := LookupTraitInPackage(name, pck)
    if trt != nil {
        return trt
    }

    // lookup containers in loaded packages
    for _, packname := range pck.IncludedPackages {
        pack := pck.LoadedPackages[packname]
        trt := LookupTraitInPackage(name, pack)
        if trt != nil {
            return trt
        }
    }

    // got nothin man
    return nil
}

func LookupTraitInPackage(name string, pack *symbols.PackageSymbol) *symbols.TraitSymbol {
    for _, v := range pack.Traits {
        if v.TraitName == name {
            return v
        }
    }

    return nil
}

// --------------------------------------------------------
// Container Lookup
// --------------------------------------------------------
func (bin *Binder) LookupContainerInPackage(pck string, cnt string) *symbols.ContainerSymbol {
    pack := bin.LookupPackage(pck)

    if pack == nil {
        return nil
    }

    return LookupContainerInPackage(cnt, pack)
}

func LookupContainer(name string, pck *symbols.PackageSymbol) *symbols.ContainerSymbol {
    cnt := LookupContainerInPackage(name, pck)
    if cnt != nil {
        return cnt
    }

    // lookup containers in loaded packages
    for _, packname := range pck.IncludedPackages {
        pack := pck.LoadedPackages[packname]

        cnt := LookupContainerInPackage(name, pack)
        if cnt != nil {
            return cnt
        }
    }

    // got nothin man
    return nil
}

func LookupContainerInPackage(name string, pack *symbols.PackageSymbol) *symbols.ContainerSymbol {
    for _, v := range pack.Containers {
        if v.ContainerName == name {
            return v
        }
    }

    return nil
}

// --------------------------------------------------------
// Field lookup
// --------------------------------------------------------
func LookupFieldInContainer(name string, cnt *symbols.ContainerSymbol) *symbols.FieldSymbol {
    for _, v := range cnt.Fields {
        if v.FieldName == name {
            return v
        }
    }

    return nil
}

func LookupFieldInTrait(name string, trt *symbols.TraitSymbol) *symbols.FieldSymbol {
    for _, v := range trt.Fields {
        if v.FieldName == name {
            return v
        }
    }

    return nil
}

// --------------------------------------------------------
// Global lookup
// --------------------------------------------------------
func (bin *Binder) LookupGlobal(name string) *symbols.GlobalSymbol {
    for _, v := range bin.CurrentPackage.Globals {
        if v.Name() == name {
            return v
        }
    }

    return nil
}

// --------------------------------------------------------
// Package lookup
// --------------------------------------------------------
func (bin *Binder) LookupPackage(name string) *symbols.PackageSymbol {
    // wait, is this us?
    if bin.CurrentPackage.Name() == name {
        return bin.CurrentPackage
    }

    // look the package up
    for _, p := range bin.CurrentPackage.LoadedPackages {
        if p.Name() == name {
            return p
        }
    }
    
    return nil
}

func LookupPackageInPackage(pack *symbols.PackageSymbol, name string) *symbols.PackageSymbol {
    // wait, is this us?
    if pack.Name() == name {
        return pack 
    }

    // look the package up
    for _, p := range pack.LoadedPackages {
        if p.Name() == name {
            return p
        }
    }
    
    return nil
}

// --------------------------------------------------------
// Function Lookup
// --------------------------------------------------------
func (bin *Binder) LookupFunction(name string) *symbols.FunctionSymbol {
    // if we're currently in a type -> look up methods first
    if bin.CurrentType != nil {
        meth := bin.LookupMethod(name, bin.CurrentType)
        if meth != nil {
            return meth
        }
    }

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

func (bin *Binder) LookupFunctionInPackage(pack string, name string) *symbols.FunctionSymbol {
    pck := bin.LookupPackage(pack)

    // did we find something?
    if pck == nil {
        return nil
    }

    return LookupFunctionInPackage(pck, name)
}

func LookupFunctionInPackage(pck *symbols.PackageSymbol, name string) *symbols.FunctionSymbol {
    for _, v := range pck.Functions {
        if v.FunctionKind == symbols.FT_FUNC && v.FuncName == name {
            return v
        }
    }

    return nil
}

// --------------------------------------------------------
// Method Lookup
// --------------------------------------------------------
func (bin *Binder) LookupMethod(name string, typ *symbols.TypeSymbol) *symbols.FunctionSymbol {
    // look in local package first 
    fnc := LookupMethodInPackage(bin.CurrentPackage, name, typ)

    if fnc != nil {
        return fnc
    }

    // if we didnt find anything -> start looking through included packages
    for _, pck := range bin.CurrentPackage.LoadedPackages {
        fnc := LookupMethodInPackage(pck, name, typ)
        if fnc != nil {
            return fnc
        }
    }

    // we got nothin man
    return nil
}

func LookupMethodInPackage(pck *symbols.PackageSymbol, name string, typ *symbols.TypeSymbol) *symbols.FunctionSymbol {
    for _, v := range pck.Functions {
        if v.FunctionKind == symbols.FT_METH && v.FuncName == name {

            // make sure this method applies for this type
            // -------------------------------------------

            // this method applies to all types
            if v.MethodKind == symbols.MT_ALL {
                return v
            }

            // this method applies to all types of a group
            if v.MethodKind == symbols.MT_GROUP && v.MethodSource.TypeGroup == typ.TypeGroup {
                return v
            }


            // this method only applies to one specific type
            if v.MethodKind == symbols.MT_STRICT && v.MethodSource.Equal(typ) {
                return v
            }

        }
    }

    return nil
}
