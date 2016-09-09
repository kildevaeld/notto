package global

import (
	"errors"
	"fmt"
	"strings"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

func getStringList(call otto.FunctionCall) []string {
	var out []string
	for _, a := range call.ArgumentList {
		if a.IsString() {
			if s, e := a.ToString(); e == nil {
				out = append(out, s)
			}
		} else if a.IsNumber() {
			if s, e := a.ToInteger(); e == nil {
				out = append(out, fmt.Sprintf("%d", s))
			}
		} else {
			out = append(out, a.String())
		}

	}
	return out
}

func Define(vm *notto.Notto) error {
	if v, err := vm.Get("process"); err != nil {
		return err
	} else if !v.IsUndefined() {
		return nil
	}

	attr := vm.ProcessAttr()

	o, e := vm.Object("({})")
	if e != nil {
		return e
	}

	o.Set("log", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		if attr.Stdout != nil {
			attr.Stdout.Write([]byte(strings.Join(str, " ") + "\n"))
		}

		return otto.UndefinedValue()
	})

	o.Set("error", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		if attr.Stderr != nil {
			attr.Stderr.Write([]byte(strings.Join(str, " ") + "\n"))
		}

		return otto.UndefinedValue()
	})

	o.Set("warn", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		if attr.Stderr != nil {
			attr.Stderr.Write([]byte(strings.Join(str, " "+"\n")))
		}

		return otto.UndefinedValue()
	})

	o.Set("dir", func(call otto.FunctionCall) otto.Value {
		o := call.Argument(0) //.Object()

		str := print(vm, o, "", true, true)

		if attr.Stdout != nil {
			attr.Stdout.Write([]byte(str + "\n"))
		}

		return otto.UndefinedValue()
	})

	return vm.Set("console", o)

	//return nil
}

func print(vm *notto.Notto, o otto.Value, indent string, pretty bool, first bool) string {
	//str := ""
	//fmt.Printf("CLASS %v Value: %#v\n", o.Class(), o.IsObject())
	switch o.Class() {
	case "Object":
		return print_object(vm, o, indent, pretty, first)
	case "Function":
		return "Function"
	case "Number", "Boolean", "Date", "RegExp":
		return o.String()
	case "String":
		return "\"" + o.String() + "\""
	case "Array":
		//return ""
		return print_array(vm, o, indent, pretty, first)
	default:
		return ""
		//return o.Class() + o.
	}

	/*if o.IsFunction() {
		str = "[Function]"
	} else if o.IsPrimitive() {
		str = o.String()
	} else if o.IsObject() {

		if o.Class() == "Array" {
			if v, e := o.Export(); e != nil {
				return e
			}
			slice, ok := v.([]interface{})
			if !ok {
				return errors.New("should not happen")
			}

		}
	}*/
}

func print_array(vm *notto.Notto, value otto.Value, indent string, pretty bool, first bool) string {
	var v interface{}
	var e error
	if v, e = value.Export(); e != nil {
		panic(e)
	}
	slice, ok := v.([]map[string]interface{})
	if !ok {
		//return print_object(vm, value, indent, pretty, true)
		panic(errors.New(fmt.Sprintf("%#v", v)))
	}

	return "Array"
	str := "["
	if !first && pretty {

		str += "\n" + indent
	} else if !first {
		indent += " "
	}
	l := len(slice)
	for i, v := range slice {
		ov, _ := vm.ToValue(v)
		str = print(vm, ov, indent, pretty, false)
		//i, _ := strconv.Atoi(k)
		if i < l-1 {
			str += ","
		}
		if pretty {
			str += "\n"
		} else {
			str += " "
		}
	}

	str += "]"

	if pretty {
		str += "\n"
	}

	first = false

	return str
}

func print_object(vm *notto.Notto, value otto.Value, indent string, pretty bool, first bool) string {
	o := value.Object()
	str := indent + "{"
	if pretty {
		str += "\n"
		indent = indent + "  "
	} else {
		str += " "
	}

	l := len(o.Keys())
	for i, k := range o.Keys() {
		str += indent + k + ": "
		if v, e := o.Get(k); e != nil {
			vm.Throw("TypeError", e)
			//return e.Error()
		} else {
			/*if v.IsBoolean() {
				b, _ := v.ToBoolean()
				str += fmt.Sprintf("%v", b)
			} else if v.IsNull() {
				str += "null"
			} else if v.IsString() {
				str += "'" + v.String() + "'"
			} else if v.IsNumber() {
				str += v.String()
			} else if v.IsObject() {
				str += print_object(v.Object(), indent+"  ", pretty)
			} else {
				return ""
			}*/
			str += print(vm, v, indent+"  ", pretty, first)
		}
		if i < l-1 {
			str += ","
		}
		if pretty {
			str += "\n"
		} else {
			str += " "
		}

	}
	str += "}\n"
	return str
}
