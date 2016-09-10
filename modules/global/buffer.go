package global

import (
	"fmt"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

type buffer struct {
	buf []byte
	vm  *notto.Notto
}

var bufferSource = `

var Buffer = (function () {
	function Buffer (_private_buffer) {
    this._private = _private_buffer;
};

Buffer.prototype.toString = function (encoding) {
	if (encoding == null) {
		var m = this._private.ToSlice()
		var t = false;
		if (m.length > 100) {
			m = m.slice(0, 100)
			t = true;
		}
		return '<Buffer ' + m.join(' ') + (t ? " ... " : "")  + '>'
	}
    return this._private.ToString(encoding);
}

Object.defineProperty(Buffer.prototype, "length", {
    get: function () {
        return this._private.Length();
    }
});

return Buffer;
})()

`

func (self *buffer) ToString(call otto.FunctionCall) otto.Value {
	str := string(self.buf)

	switch call.Argument(0).String() {
	case "utf8", "utf-8":
		str = string(self.buf)
	case "hex":
		str = fmt.Sprintf("%x", self.buf)

	}

	v, e := self.vm.ToValue(str)
	if e != nil {
		self.vm.Throw("BufferError", "")
	}
	return v
}

func (self *buffer) Length() int64 {
	return int64(len(self.buf))
}

func (self *buffer) ToSlice() []uint8 {
	s := make([]uint8, self.Length())
	for i, b := range self.buf {
		s[i] = b
	}
	return s
}

func NewBuffer(vm *notto.Notto, bs []byte) otto.Value {

	//
	buf := &buffer{bs, vm}

	v, _ := vm.Object("({})")
	v.Set("__private", buf)
	v.Set("toString", buf.ToString)
	v.Set("length", len(bs))
	/*v, e := vm.Call("new Buffer", nil)
	if e != nil {
		vm.Throw("BufferError", e)
	}*/
	return v.Value()

}
