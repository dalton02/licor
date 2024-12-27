package validator

import (
	"reflect"
)

type MapJson struct {
	Mapa      map[int][]string
	Level     int
	Type      string
	MapaQuery map[string]string
}

func getMapaField(level int, field string, tabela map[int][]string) bool {
	if stringsNivel, ok := tabela[level]; ok {
		for _, str := range stringsNivel {
			if str == field {
				return true
			}
		}
	}
	return false
}

func CheckPropretys[T any](data T, mapaJson MapJson) (string, bool) {

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	errorMessage := ""
	hasError := false
	if t.Kind() == reflect.Map {
		return "", hasError
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		typeValidation := field.Tag.Get("validator")

		toValidate := CommonToArray(typeValidation)
		hasOp := hasOptional(toValidate)

		if hasOp && mapaJson.MapaQuery[field.Tag.Get("query")] == "" && mapaJson.Type == "query" {
			continue
		}

		if hasOp && !getMapaField(mapaJson.Level, field.Tag.Get("json"), mapaJson.Mapa) && mapaJson.Type != "query" {
			continue
		}

		firstErr := true
		currentErrMessage := ""
		for i := 0; i < len(toValidate); i++ {
			test := string(testCase(toValidate[i], value, field, mapaJson))
			if firstErr && test != "" {
				currentErrMessage += "must be " + test
				firstErr = false
				hasError = true
			} else if test != "" {
				currentErrMessage += "," + test
				hasError = true
			}
		}
		if errorMessage != "" && currentErrMessage != "" {
			errorMessage += ";"
			hasError = true
		}
		if currentErrMessage != "" {
			nameCamp := field.Tag.Get("json")
			if nameCamp == "" {
				nameCamp = field.Tag.Get("query")
			}
			errorMessage += nameCamp + ":" + currentErrMessage
			hasError = true
		}

	}
	return errorMessage, hasError

}

func testCase(field string, value reflect.Value, primitive reflect.StructField, mapaJson MapJson) string {
	switch field {
	case "numericString":
		return isNumericString(value)
	case "strongPassword":
		return isStrongPassword(value)
	case "email":
		return isEmail(value)
	case "dateString":
		return isDate(value)
	case "required":
		return testPrimitive(primitive, value, mapaJson)
	case "optional":
		return testPrimitive(primitive, value, mapaJson)
	default:
		return ""
	}
}

func testPrimitive(primitive reflect.StructField, value reflect.Value, mapaJson MapJson) string {
	isSet := value.IsValid() && !value.IsZero()
	switch primitive.Type.Kind() {
	case reflect.Int:
		_, test := value.Interface().(int)
		if test && isSet {
			return ""
		}
		return "a number"
	case reflect.String:
		_, test := value.Interface().(string)
		if test && isSet {
			return ""
		}
		return "a string"
	case reflect.Float64:
		_, test := value.Interface().(float64)
		if test && isSet {
			return ""
		}
		return "a float"
	case reflect.Bool:
		_, test := value.Interface().(bool)
		if test && isSet {
			return ""
		}
		return "a boolean"
	case reflect.Struct:
		mapaJson.Level++
		err, has := CheckPropretys(value.Interface(), mapaJson)
		if has {
			mapaJson.Level--
			return "(" + err + ")"
		}
		return ""
	case reflect.Slice:
		return testSlices(primitive, value, mapaJson)
	default:
		return "unsupported type"
	}
}

func testSlices(primitive reflect.StructField, value reflect.Value, mapaJson MapJson) string {
	isSet := value.IsValid() && value.Len() > 0
	switch primitive.Type.Elem().Kind() {
	case reflect.String:
		if isSet {
			_, test := value.Interface().([]string)
			if test {
				return ""
			}
		}
		return "a []string with at least one element"
	case reflect.Int:
		if isSet {
			_, test := value.Interface().([]int)
			if test {
				return ""
			}
		}
		return "a []int with at least one element"
	case reflect.Float64:
		if isSet {
			_, test := value.Interface().([]float64)
			if test {
				return ""
			}
		}
		return "a []float64 with at least one element"
	case reflect.Float32:
		if isSet {
			_, test := value.Interface().([]float32)
			if test {
				return ""
			}
		}
		return "a []float32 with at least one element"
	case reflect.Bool:
		if isSet {
			_, test := value.Interface().([]bool)
			if test {
				return ""
			}
		}
		return "a []bool with at least one element"
	case reflect.Int64:
		if isSet {
			_, test := value.Interface().([]int64)
			if test {
				return ""
			}
		}
		return "a []int64 with at least one element"
	case reflect.Struct:
		if isSet {
			mapaJson.Level++
			for i := 0; i < value.Len(); i++ {
				msg, hasError := CheckPropretys(value.Index(i).Interface(), mapaJson)
				if hasError {
					mapaJson.Level--
					return "Array{" + msg + "}"
				}
			}
			mapaJson.Level--
		} else {
			valuePtr := reflect.New(value.Type().Elem())
			valueElem := valuePtr.Elem() // Isso obtém o valor do ponteiro
			mapaJson.Level--
			return "" + showPropretys(valueElem) + ""
		}
	default:
		return "unsupported type in slice"
	}
	return ""
}

func showPropretys(data reflect.Value) string {
	t := data.Type()
	var msg string = "[]{"
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		isToBe := CommonToArray(string(field.Tag.Get("validator")))
		msg += field.Tag.Get("json") + ":" + field.Type.Name()
		if len(isToBe) > 0 {
			msg += "=>("
		}
		for i := 0; i < len(isToBe); i++ {
			msg += isToBe[i]
			if i == len(isToBe)-1 {
				msg += ")"
			} else {
				msg += ","
			}

		}
		if i != t.NumField()-1 {
			msg += ";"
		}
	}
	msg += "}"
	return msg

}

func hasOptional(tipos []string) bool {
	for i := 0; i < len(tipos); i++ {
		if tipos[i] == "optional" {
			return true
		}
	}
	return false
}

func ExtractKeysByLevel(data map[string]interface{}, level int, result map[int][]string) MapJson {
	for key, value := range data {
		result[level] = append(result[level], key)
		if reflect.TypeOf(value).Kind() == reflect.Map {
			subMap := value.(map[string]interface{})
			ExtractKeysByLevel(subMap, level+1, result)
		}
	}
	mapJ := MapJson{
		Mapa:  result,
		Level: 1,
		Type:  "json",
	}
	return mapJ
}
func QueryMap(querys map[string]string) MapJson {
	mapJ := MapJson{
		Mapa:      make(map[int][]string),
		Level:     1,
		Type:      "query",
		MapaQuery: querys,
	}
	return mapJ
}

func CommonToArray(data string) []string {
	var array []string
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == ',' {
			newString := data[start:i]
			array = append(array, newString)
			start = i + 1
		}
	}

	if start < len(data) {
		newString := data[start:]
		array = append(array, newString)
	}

	return array

}
