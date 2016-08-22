package notto

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"errors"

	"github.com/kildevaeld/ottoext/fetch"
	"github.com/kildevaeld/ottoext/loop"
	"github.com/kildevaeld/ottoext/promise"
	"github.com/kildevaeld/ottoext/timers"
	"github.com/robertkrimen/otto"
)

type Environ []string

func (self Environ) ToMap() map[string]string {
	env := make(map[string]string)
	for _, e := range self {
		a := strings.SplitN(e, "=", 2)
		env[a[0]] = a[1]
	}
	return env
}

func MapToEnviron(m map[string]string) Environ {
	var out Environ
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

type ProcessAttr struct {
	Environ Environ
	Argv    []string
	Stdout  io.Writer
	Stderr  io.Writer
	Cwd     string
}

// Globally registered modules
var globalModules map[string]ModuleLoader = make(map[string]ModuleLoader)

// Globally registered paths (paths to search for modules)
var globalPaths []string

type Notto struct {
	*otto.Otto
	// Modules that registered for current vm
	modules map[string]ModuleLoader
	// Location to search for modules
	paths []string
	// Once a module is required by vm, the exported value is cached for further
	// use.
	moduleCache map[string]otto.Value
	runLoop     *loop.Loop
	preScripts  []string
	processAttr *ProcessAttr
}

func (this *Notto) Runloop() *loop.Loop {
	return this.runLoop
}

func (this *Notto) Init() error {
	if err := timers.Define(this.Otto, this.runLoop); err != nil {
		return err
	}
	if err := promise.Define(this.Otto, this.runLoop); err != nil {
		return err
	}
	if err := fetch.Define(this.Otto, this.runLoop); err != nil {
		return err
	}

	return nil
}

func (self *Notto) SetProcessAttr(attr *ProcessAttr) {
	self.processAttr = attr
}

func (self *Notto) ProcessAttr() *ProcessAttr {
	if self.processAttr == nil {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		self.processAttr = &ProcessAttr{
			Environ: os.Environ(),
			Argv:    os.Args,
			Cwd:     cwd,
			Stderr:  os.Stderr,
			Stdout:  os.Stdout,
		}
	}

	return self.processAttr
}

func (this *Notto) AddPreScript(script string) {
	this.preScripts = append(this.preScripts, script)
}

func (this *Notto) RunScript(script, pwd string) (otto.Value, error) {
	var (
		val otto.Value
		err error
	)
	if val, err = CreateLoaderFromSource(script, pwd)(this); err != nil {
		return otto.NullValue(), err
	}
	if err = this.runLoop.Run(); err != nil {
		return otto.NullValue(), err
	}
	return val, nil
}

// Run a module or file
func (this *Notto) Run(name string, pwd string) (otto.Value, error) {
	if ok, _ := isFile(name); ok {
		name, _ = filepath.Abs(name)
	}
	var (
		val otto.Value
		err error
	)
	if val, err = this.Require(name, pwd); err != nil {
		return otto.NullValue(), err
	}
	if err = this.runLoop.Run(); err != nil {
		return otto.NullValue(), err
	}
	return val, err
}

// Require a module with cache
func (this *Notto) Require(id, pwd string) (otto.Value, error) {

	if cache, ok := this.moduleCache[id]; ok {
		return cache, nil
	}

	loader, ok := this.modules[id]
	if !ok {
		loader, ok = globalModules[id]
	}

	if loader != nil {
		value, err := loader(this)
		if err != nil {
			return otto.UndefinedValue(), err
		}

		this.moduleCache[id] = value
		return value, nil
	}

	filename, err := FindFileModule(id, pwd, append(this.paths, globalPaths...))
	if err != nil {
		return otto.UndefinedValue(), err
	}

	// resove id
	id = filename

	if cache, ok := this.moduleCache[id]; ok {
		return cache, nil
	}

	v, err := CreateLoaderFromFile(id)(this)

	if err != nil {
		return otto.UndefinedValue(), err
	}

	// cache
	this.moduleCache[id] = v

	return v, nil
}

// Must be thrown from within the otto context
func (this *Notto) Throw(name string, err interface{}) error {
	var v otto.Value
	switch s := err.(type) {
	case string:
		v = this.MakeCustomError(name, s)
	case error:
		v = this.MakeCustomError(name, s.Error())
	default:
		return errors.New("not a string or an error")
	}
	panic(v)
	return nil
}

// Register a new module to current vm.
func (this *Notto) AddModule(id string, loader ModuleLoader) {
	this.modules[id] = loader
}

func (this *Notto) ClearCache() {
	this.modules = make(map[string]ModuleLoader)
}

// Add paths to search for modules.
func (this *Notto) AddPath(paths ...string) {
	this.paths = append(this.paths, paths...)
}

// Register a global module
func AddModule(id string, m ModuleLoader) {
	globalModules[id] = m
}

// Register global path.
func AddPath(paths ...string) {
	globalPaths = append(globalPaths, paths...)
}

// Run module by name in the Notto module environment.
func Run(name string) (*Notto, otto.Value, error) {
	vm := New()
	v, err := vm.Run(name, ".")

	return vm, v, err
}

// Create a Notto vm instance.
func New() *Notto {
	m := &Notto{
		Otto:        otto.New(),
		modules:     make(map[string]ModuleLoader),
		paths:       nil,
		moduleCache: make(map[string]otto.Value),
	}
	m.runLoop = loop.New(m.Otto)
	return m
}
