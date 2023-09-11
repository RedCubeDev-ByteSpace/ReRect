// Go packages - sys.go
// --------------------------------------------------------
// ah yes, Rects sys library, its good to be back
// --------------------------------------------------------
package gopackages

import (
	"bytespace.network/rerect/compunit"
	evalobjects "bytespace.network/rerect/eval_objects"
	"bytespace.network/rerect/symbols"
	"fmt"
)

// LoadExample must be added to load() in load.go (and recompile)
func LoadExample() {
	packName := "example"
	example := registerPackage("example")

	registerFunction(packName, symbols.NewVMFunctionSymbol(
		// The package this function belongs to
		example,

		// Function name
		"Add",

		// Return type
		compunit.GlobalDataTypeRegister["int"],

		// Function parameters
		[]*symbols.ParameterSymbol{
			symbols.NewParameterSymbol("a", 0, compunit.GlobalDataTypeRegister["int"]),
			symbols.NewParameterSymbol("b", 0, compunit.GlobalDataTypeRegister["int"]),
		},

		// Pointer to function
		Add,
	))

	hotdogTypeSymbol := symbols.NewTypeSymbol("Hotdog", []*symbols.TypeSymbol{}, symbols.CONT, 0, nil)
	hotdogContainer := symbols.NewContainerSymbol(example, "Hotdog", hotdogTypeSymbol)
	hotdogTypeSymbol.Container = hotdogContainer // "doubly linked" more like "doubly ludicrous" >:(

	registerContainer(packName, hotdogContainer)

	hotdogContainer.Fields = append(hotdogContainer.Fields, symbols.NewFieldSymbol(hotdogContainer, "name", compunit.GlobalDataTypeRegister["string"]))
	hotdogContainer.Constructor = symbols.NewVMMethodSymbol(
		example,
		symbols.MT_GROUP,
		hotdogTypeSymbol,
		"Constructor",
		compunit.GlobalDataTypeRegister["void"],
		[]*symbols.ParameterSymbol{
			symbols.NewParameterSymbol("name", 0, compunit.GlobalDataTypeRegister["string"]),
		},
		Hotdog_Constructor,
	)

	registerFunction(
		packName,
		symbols.NewVMMethodSymbol(
			example,
			symbols.MT_GROUP,
			hotdogTypeSymbol,
			"Dance",
			compunit.GlobalDataTypeRegister["void"],
			[]*symbols.ParameterSymbol{},
			Hotdog_Dance,
		),
	)

	registerFunction(
		packName,
		symbols.NewVMMethodSymbol(
			example,
			symbols.MT_GROUP,
			hotdogTypeSymbol,
			"Debug",
			compunit.GlobalDataTypeRegister["string"],
			[]*symbols.ParameterSymbol{},
			Hotdog_Debug,
		),
	)
}

func Add(args []any) any {
	a := args[0].(int32)
	b := args[1].(int32)

	return a + b
}

func Hotdog_Dance(instance any, args []any) any {
	if instance == nil {
		return nil
	}

	con := instance.(*evalobjects.ContainerInstance)

	for i := 0; i < 10; i++ {
		fmt.Printf("Hotdog(%s): dancin'\n", con.Fields["name"])
	}

	return nil
}

func Hotdog_Constructor(instance any, args []any) any {
	if instance == nil {
		return nil
	}

	con := instance.(*evalobjects.ContainerInstance)
	con.Fields["name"] = args[0].(string)

	return nil
}

func Hotdog_Debug(instance any, args []any) any {
	if instance == nil {
		return nil
	}

	con := instance.(*evalobjects.ContainerInstance)

	return fmt.Sprintf("Hotdog: {\n  Type: %v\n  Fields: %v\n}\n", *con.Type, con.Fields)
}
