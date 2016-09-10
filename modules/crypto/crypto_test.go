package crypto

import (
	"testing"

	"github.com/kildevaeld/notto/modules/promise"

	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/modules/global"
)

func TestCrypto(t *testing.T) {

	vm := notto.New()
	promise.Define(vm)
	global.Define(vm)
	e := Define(vm)
	if e != nil {
		t.Fatal(e)
	}
	_, e = vm.RunScript(`
        var crypto = require('crypto');
        crypto.randomBytes(120).then(function (b) {
            console.log(b)
        })
    `, "")

	if e != nil {
		t.Fatal(e)
	}

}
