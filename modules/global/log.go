package global

import (
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
			attr.Stdout.Write([]byte(strings.Join(str, " ")))
		}

		return otto.UndefinedValue()
	})

	o.Set("error", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		if attr.Stderr != nil {
			attr.Stderr.Write([]byte(strings.Join(str, " ")))
		}

		return otto.UndefinedValue()
	})

	o.Set("warn", func(call otto.FunctionCall) otto.Value {
		str := getStringList(call)
		if attr.Stderr != nil {
			attr.Stderr.Write([]byte(strings.Join(str, " ")))
		}

		return otto.UndefinedValue()
	})

	o.Set("dir", func(call otto.FunctionCall) otto.Value {
		o := call.Argument(0).Object()

		str := print_object(o, "", false)

		if attr.Stdout != nil {
			attr.Stdout.Write([]byte(str))
		}

		return otto.UndefinedValue()
	})

	return vm.Set("console", o)

	//return nil
}

func print_object(o *otto.Object, indent string, pretty bool) string {
	str := indent + "{"
	if pretty {
		str += "{\n}"
		indent = "  "
	} else {
		str += " "
	}

	l := len(o.Keys())
	for i, k := range o.Keys() {
		str += indent + k + ": "
		if v, e := o.Get(k); e != nil {
			return e.Error()
		} else {
			if v.IsBoolean() {
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
			}
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
