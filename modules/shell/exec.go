package shell

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/kildevaeld/notto"
	"github.com/robertkrimen/otto"
)

type result struct {
	stderr []byte
	stdout []byte
	code   int
}

func execCall(vm *notto.Notto, rootDir string) func(call otto.FunctionCall) otto.Value {

	return func(call otto.FunctionCall) otto.Value {

		if len(call.ArgumentList) == 0 {
			return otto.UndefinedValue()
		}

		if !call.Argument(0).IsString() {
			return otto.UndefinedValue()
		}
		silence := true
		if call.Argument(1).IsBoolean() {
			silence, _ = call.Argument(1).ToBoolean()
		}

		var s string
		var e error
		var r *result
		if s, e = call.Argument(0).ToString(); e != nil {
			panic(e)
		}

		if r, e = execImpl(vm, s, rootDir, silence); e != nil {
			/*err := vm.MakeCustomError("ExecError", e.Error())
			panic(err)*/
			vm.Throw("ExecError", e)
		}

		//m := make(map[string]interface{})

		ob, _ := vm.Object("({})")
		ob.Set("stderr", strings.TrimSpace(string(r.stderr)))
		ob.Set("stdout", strings.TrimSpace(string(r.stdout)))
		ob.Set("pipe", pipeFn(string(r.stdout)))

		return ob.Value()
	}
}

func execImpl(vm *notto.Notto, cmd, root string, silence bool) (*result, error) {

	split := strings.Split(cmd, " ")

	c := exec.Command(split[0], split[1:]...)

	sbuf := bytes.NewBuffer(nil)
	ebuf := bytes.NewBuffer(nil)

	var stdout io.Writer = sbuf
	var stderr io.Writer = ebuf

	if !silence {
		attr := vm.ProcessAttr()
		stdout = io.MultiWriter(sbuf, attr.Stdout)
		stderr = io.MultiWriter(ebuf, attr.Stderr)
	}

	c.Stderr = stderr
	c.Stdout = stdout

	c.Dir = root
	c.Env = vm.ProcessAttr().Environ

	if runtime.GOOS == "linux" && runtime.GOOS == "darwin" {
		u, err := user.Current()

		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)
		if err == nil {
			c.SysProcAttr = &syscall.SysProcAttr{
				Chroot: root,
				//Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
				//Setpgid:    true,
				Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
			}
		}
	}

	e := c.Run()
	if ee, ok := e.(*exec.ExitError); ok {

		if !ee.Success() {
			return &result{sbuf.Bytes(), sbuf.Bytes(), 0}, errors.New(string(ebuf.Bytes()))
		}
	} else if e != nil {
		return nil, e
	}
	return &result{ebuf.Bytes(), sbuf.Bytes(), 0}, nil

}
