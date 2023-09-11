package evalobjects

import (
	"fmt"
	"strconv"

	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/symbols"
)

func EvalConversion(val interface{}, to *symbols.TypeSymbol) (interface{}, bool) {
    // Casting anything to 'any'
    if to.Equal(compunit.GlobalDataTypeRegister["any"]) {
        return interface{}(val), true
    }

    // Casting to long
    if to.Equal(compunit.GlobalDataTypeRegister["long"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return val, true

        case int32:
            return int64(v), true

        case int16:
            return int64(v), true

        case int8:
            return int64(v), true

        // Cross cast
        // ----------
        case float64:
            return int64(v), true

        case float32:
            return int64(v), true

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 64)
            if err != nil {
                panic(err)
            }

            return int64(vl), true
        }
    }

    // Casting to int
    if to.Equal(compunit.GlobalDataTypeRegister["int"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return int32(v), true

        case int32:
            return v, true

        case int16:
            return int32(v), true

        case int8:
            return int32(v), true

        // Cross cast
        // ----------
        case float64:
            return int32(v), true

        case float32:
            return int32(v), true

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 32)
            if err != nil {
                panic(err)
            }

            return int32(vl), true
        }
    }
    
    // Casting to word
    if to.Equal(compunit.GlobalDataTypeRegister["word"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return int16(v), true

        case int32:
            return int16(v), true

        case int16:
            return v, true

        case int8:
            return int16(v), true

        // Cross cast
        // ----------
        case float64:
            return int16(v), true

        case float32:
            return int16(v), true

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 16)
            if err != nil {
                panic(err)
            }

            return int16(vl), true
        }
    }
    
    // Casting to byte
    if to.Equal(compunit.GlobalDataTypeRegister["byte"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case int64:
            return int8(v), true

        case int32:
            return int8(v), true

        case int16:
            return int8(v), true

        case int8:
            return v, true

        // Cross cast
        // ----------
        case float64:
            return int8(v), true

        case float32:
            return int8(v), true

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseInt(v, 10, 8)
            if err != nil {
                panic(err)
            }

            return int8(vl), true
        }
    }
    
    // Casting to double
    if to.Equal(compunit.GlobalDataTypeRegister["double"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case float64:
            return v, true

        case float32:
            return float64(v), true

        // Cross casts
        // -----------
        case int64:
            return float64(v), true

        case int32:
            return float64(v), true

        case int16:
            return float64(v), true

        case int8:
            return float64(v), true

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseFloat(v, 64)
            if err != nil {
                panic(err)
            }

            return float64(vl), true
        }
    }

    // Casting to float
    if to.Equal(compunit.GlobalDataTypeRegister["float"]) {
        switch v := val.(type) {

        // Up / Down casts
        // ---------------
        case float64:
            return float32(v), true

        case float32:
            return v, true

        // Cross casts
        // -----------
        case int64:
            return float32(v), true

        case int32:
            return float32(v), true

        case int16:
            return float32(v), true

        case int8:
            return float32(v), true

        // From string
        // -----------
        case string:
            vl, err := strconv.ParseFloat(v, 32)
            if err != nil {
                panic(err)
            }

            return float32(vl), true
        }
    }

    // Casting to bool
    if to.Equal(compunit.GlobalDataTypeRegister["bool"]) {
        switch v := val.(type) {
        case bool:
            return v, true

        case string:
            vl, err := strconv.ParseBool(v)
            if err != nil {
                panic(err)
            }

            return vl, true
        }
    }

    // Casting to string
    if to.Equal(compunit.GlobalDataTypeRegister["string"]) {
        switch v := val.(type) {

        // Integers
        // --------
        case int64:
            return fmt.Sprintf("%d", v), true

        case int32:
            return fmt.Sprintf("%d", v), true

        case int16:
            return fmt.Sprintf("%d", v), true

        case int8:
            return fmt.Sprintf("%d", v), true

        // Floats
        // ------
        case float64:
            return fmt.Sprintf("%f", v), true

        case float32:
            return fmt.Sprintf("%f", v), true

        // Booleans
        // --------
        case bool:
            if v {
                return "true", true
            } else {
                return "false", true
            }

        case *ArrayInstance:
            return fmt.Sprintf("[%s]", v.Type.Name()), true

        // Strings
        // -------
        case string:
            return v, true
        }
    }

    // Casting to array
    if to.TypeGroup == symbols.ARR {
        switch v := val.(type) {
        case *ArrayInstance:
            // only cast when the internal types match
            if v.Type.Equal(to) {
                return v, true
            }
        }
    }

    // Casting to container
    if to.TypeGroup == symbols.CONT {
        switch v := val.(type) {
        case *ContainerInstance:
            // only cast when the internal types match
            if v.Type.Equal(to) {
                return v, true
            }
        }
    }

    return nil, false
}
