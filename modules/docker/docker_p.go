package docker

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/kildevaeld/notto"
	"github.com/kildevaeld/notto/loop"
	"github.com/robertkrimen/otto"
)

func mustValue(value otto.Value, err error) otto.Value {
	if err != nil {
		panic(err)
	}
	return value
}

func boolOr(o *otto.Object, prop string, d bool) bool {
	v, e := o.Get(prop)
	if e != nil {
		return d
	}
	if b, e := v.ToBoolean(); e != nil {
		return d
	} else {
		return b
	}
}

func boolOrFale(o otto.Value) bool {
	if o.IsBoolean() {
		b, _ := o.ToBoolean()
		return b
	}
	return false
}

func stringOr(o *otto.Object, prop string, or string) string {
	v, e := o.Get(prop)
	if e != nil {
		return or
	}
	if !v.IsString() {
		return or
	}
	if b, e := v.ToString(); e != nil {
		return or
	} else {
		return b
	}
}

type docker_p_task struct {
	id       int64
	err      error
	result   interface{}
	callback otto.Value
	name     string
}

func (self *docker_p_task) GetID() int64   { return self.id }
func (self *docker_p_task) SetID(id int64) { self.id = id }

func (self *docker_p_task) Execute(vm *otto.Otto, l *loop.Loop) error {

	var arguments []interface{}

	if self.err != nil {
		e := vm.MakeCustomError("DockerError", self.err.Error())

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}
	if self.result != nil {
		if v, e := vm.ToValue(self.result); e != nil {
			return e
		} else {
			arguments = append(arguments, v)
		}

	}

	if _, err := self.callback.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}

	return nil

}

func (self *docker_p_task) Cancel() {

}

type docker_p struct {
	client *dockerclient.Client
	vm     *notto.Notto
}

func contains(slice []string, n string) bool {
	for _, s := range slice {
		if s == n {
			return true
		}
	}
	return false
}

func getEnv(o *otto.Object) ([]string, error) {
	var array []string
	if env, err := o.Get("env"); err == nil {

		if env.Class() == "Array" {
			if cc, err := env.Object().Call("join", " "); err != nil {
				return nil, err
			} else {
				array = strings.Split(cc.String(), " ")
			}
		} else if env.Class() == "Object" {
			o := env.Object()

			for _, key := range o.Keys() {
				v := mustValue(o.Get(key)).String()
				array = append(array, fmt.Sprintf("%s=%s", key, v))
			}
		} else if !env.IsUndefined() && !env.IsNull() {
			return nil, errors.New("env must be an object or an array")
		}

	}
	return array, nil
}

func getStringSlice(o *otto.Object, prop string) ([]string, error) {

	if cmd, err := o.Get(prop); err == nil {

		var c string
		if cmd.Class() == "Array" {

			if cc, err := cmd.Object().Call("join", " "); err != nil {
				return nil, err
			} else {
				c = cc.String()
			}

		} else if cmd.Class() == "String" {
			c = cmd.String()
		}

		if c == "" {
			return nil, nil
		}

		return strings.Split(c, " "), nil
	}
	return nil, nil
}

func getCreateOptions(o *otto.Object) (dockerclient.CreateContainerOptions, error) {
	out := dockerclient.CreateContainerOptions{}
	out.Name = stringOr(o, "name", "")

	cfg := dockerclient.Config{}
	cfg.Image = stringOr(o, "image", "")

	var err error

	if cfg.Cmd, err = getStringSlice(o, "cmd"); err != nil {
		return out, err
	}
	if cfg.Cmd != nil && len(cfg.Cmd) == 0 {
		cfg.Cmd = nil
	}

	if cfg.Entrypoint, err = getStringSlice(o, "entrypoint"); err != nil {
		return out, err
	}

	if cfg.Entrypoint != nil && len(cfg.Entrypoint) == 0 {
		cfg.Entrypoint = nil
	}

	cfg.WorkingDir = stringOr(o, "workingDir", "")

	if cfg.Env, err = getEnv(o); err != nil {
		return out, err
	}

	cfg.AttachStderr = boolOr(o, "attachStderr", false)
	cfg.AttachStdin = boolOr(o, "attachStdin", false)
	cfg.AttachStdout = boolOr(o, "attachStdout", false)
	cfg.Tty = boolOr(o, "tty", false)
	cfg.OpenStdin = boolOr(o, "openStdin", false)
	cfg.StdinOnce = boolOr(o, "stdinOnce", false)

	out.Config = &cfg

	hcfg := &dockerclient.HostConfig{}

	hcfg.AutoRemove = boolOr(o, "autoRemove", false)
	hcfg.PublishAllPorts = boolOr(o, "publishAllPorts", false)
	hcfg.Privileged = boolOr(o, "previleged", false)
	if hcfg.Binds, err = getStringSlice(o, "binds"); err != nil {
		return out, err
	}

	if hcfg.Links, err = getStringSlice(o, "links"); err != nil {
		return out, err
	}

	if hcfg.Binds, err = getStringSlice(o, "volumes"); err != nil {
		return out, err
	}

	if v, e := getStringSlice(o, "publish"); e == nil {

		p := make(map[dockerclient.Port][]dockerclient.PortBinding)
		for _, k := range v {
			s := strings.Split(k, ":")
			if len(s) != 2 {
				return out, errors.New("")
			}
			p[dockerclient.Port(s[0]+"/tcp")] = []dockerclient.PortBinding{dockerclient.PortBinding{
				HostPort: s[1],
				HostIP:   "0.0.0.0",
			}}

		}
		//hcfg.PublishAllPorts = true
		hcfg.PortBindings = p

	}

	out.HostConfig = hcfg
	out.Config.ExposedPorts = map[dockerclient.Port]struct{}{
		"3306/tcp": {}}

	return out, nil

}

func getContextPath(path string) (io.Reader, error) {
	var err error
	if !filepath.IsAbs(path) {
		if path, err = filepath.Abs(path); err != nil {
			return nil, err
		}
	}

	t := time.Now()
	inputbuf := bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	defer tr.Close()

	dockerfileFound := false
	err = filepath.Walk(path, func(p string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if p == path {
			return nil
		}

		fp := p

		p = strings.Replace(p, path+"/", "", 1)

		if file.Name() == "Dockerfile" {
			if dockerfileFound {
				return errors.New("Cannot container 2 dockerfiles")
			}
			dockerfileFound = true
		}

		if file.IsDir() {
			return nil
		}

		tr.WriteHeader(&tar.Header{
			Name:       p,
			Size:       file.Size(),
			ModTime:    file.ModTime(),
			AccessTime: t,
			ChangeTime: t,
			Mode:       int64(file.Mode()),
		})

		if b, e := ioutil.ReadFile(fp); e != nil {
			return e
		} else {
			tr.Write(b)
		}

		return nil
	})

	if !dockerfileFound {
		return nil, errors.New("No dockerfiles in " + path)
	}

	return inputbuf, err
}

func getBuildOptions(o *otto.Object) (dockerclient.BuildImageOptions, error) {
	var (
		input io.Reader
		//output io.Writer
		err error
	)
	out := dockerclient.BuildImageOptions{}

	path := stringOr(o, "path", ".")

	if input, err = getContextPath(path); err != nil {
		return out, err
	}
	out.InputStream = input
	out.OutputStream = bytes.NewBuffer(nil)
	out.RmTmpContainer = boolOr(o, "rmTmpContainer", false)
	out.ForceRmTmpContainer = boolOr(o, "forceRmTmpContainer", false)
	out.Pull = boolOr(o, "pull", false)
	out.Name = stringOr(o, "name", "")

	return out, nil
}

func (self *docker_p) create(call otto.FunctionCall) (*dockerclient.Container, error) {
	obj := call.Argument(0).Object()

	o, err := getCreateOptions(obj)
	if err != nil {
		return nil, fmt.Errorf("create: %s", err)
	}

	pull := boolOr(obj, "pull", false)

	if !self.has_image(o.Config.Image) && pull {
		split := strings.Split(o.Config.Image, ":")
		o := dockerclient.PullImageOptions{Repository: split[0]}
		if len(split) == 2 {
			o.Tag = split[1]
		}
		self.client.PullImage(o, dockerclient.AuthConfiguration{})

	}

	if c, e := self.client.CreateContainer(o); e != nil {
		return nil, fmt.Errorf("create: %s - %s:  %s", o.Name, o.Config.Image, e)
	} else {
		return c, nil
	}
}

func (self *docker_p) Create(call otto.FunctionCall) otto.Value {

	sync := boolOrFale(call.Argument(1))
	if sync {
		container, err := self.create(call)
		if err != nil {
			self.vm.Throw("DockerError", err)
		}
		if val, err := self.vm.ToValue(container); err != nil {
			self.vm.Throw("TypeError", err)
		} else {
			return val
		}
	}
	task := self.getTask("create", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		obj := call.Argument(0).Object()

		o, err := getCreateOptions(obj)
		if err != nil {
			task.err = err
			return
		}

		pull := boolOr(obj, "pull", false)

		if !self.has_image(o.Config.Image) && pull {

			self.client.PullImage(dockerclient.PullImageOptions{Repository: o.Config.Image}, dockerclient.AuthConfiguration{})

		}

		if c, e := self.client.CreateContainer(o); e != nil {
			task.err = e
		} else {
			task.result = c
		}
	}()

	return otto.UndefinedValue()

}

func (self *docker_p) Build(call otto.FunctionCall) otto.Value {

	sync := boolOrFale(call.Argument(1))

	if sync {
		obj := call.Argument(0).Object()
		o, err := getBuildOptions(obj)
		if err != nil {
			self.vm.Throw("DockerError", err)
		}

		if e := self.client.BuildImage(o); e != nil {
			self.vm.Throw("DockerError", e)
		}
		return otto.UndefinedValue()
	}

	task := self.getTask("build", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		obj := call.Argument(0).Object()
		o, err := getBuildOptions(obj)
		if err != nil {
			task.err = err
			return
		}

		if e := self.client.BuildImage(o); e != nil {
			task.err = e
		}

	}()
	return otto.UndefinedValue()
}

func (self *docker_p) Attach(call otto.FunctionCall) otto.Value {

	obj := call.Argument(0).Object()

	o := dockerclient.AttachToContainerOptions{}

	o.Container = stringOr(obj, "container", "")

	o.InputStream = os.Stdin
	o.OutputStream = os.Stdout
	//o.ErrorStream = os.Stderr
	o.Stderr = true
	o.Stdout = true
	//o.Stdin = true
	//o.Stream = true
	c := make(chan struct{})
	o.Success = c
	o.RawTerminal = true
	go func() {
		<-c
		fmt.Printf("C")
	}()
	if err := self.client.AttachToContainer(o); err != nil {
		self.vm.Throw("DockerError", err)
	}

	return otto.UndefinedValue()
}

func (self *docker_p) Exec(call otto.FunctionCall) otto.Value {

	o := dockerclient.CreateExecOptions{}
	o.AttachStderr = true
	o.AttachStdin = true
	o.AttachStdout = true

	o.Container = "something"
	o.Tty = false
	o.Cmd = []string{"sh"}
	ex, e := self.client.CreateExec(o)
	if e != nil {
		panic(e)
	}
	if e := self.client.StartExec(ex.ID, dockerclient.StartExecOptions{
		Detach:       false,
		Tty:          false,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		//RawTerminal:  true,
	}); e != nil {
		panic(e)
	}
	return otto.UndefinedValue()
}

func (self *docker_p) check_args(call otto.FunctionCall) error {
	if !call.Argument(0).IsObject() || call.Argument(0).Class() != "Object" {
		return errors.New("must " + call.Argument(0).Class())
	}
	if !call.Argument(1).IsFunction() {
		return errors.New("function")
	}
	return nil
}

func (self *docker_p) start(name string, config *dockerclient.HostConfig) error {

	return self.client.StartContainer(name, config)
}

func (self *docker_p) Start(call otto.FunctionCall) otto.Value {

	name := mustValue(call.Argument(0).Object().Get("name")).String()
	sync := boolOrFale(call.Argument(1))
	if sync {
		o, e := getCreateOptions(call.Argument(0).Object())
		if e != nil {
			self.vm.Throw("DockerError", e)
		}

		if err := self.start(name, o.HostConfig); err != nil {
			self.vm.Throw("DockerError", err)
		}
		return otto.UndefinedValue()
	}

	task := self.getTask("start", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		task.err = self.client.StartContainer(name, &dockerclient.HostConfig{})
	}()

	return otto.UndefinedValue()
}

func (self *docker_p) getTask(name string, call otto.Value) *docker_p_task {
	task := &docker_p_task{
		callback: call,
		name:     name,
	}

	self.vm.Runloop().Add(task)
	return task
}

func (self *docker_p) Stop(call otto.FunctionCall) otto.Value {

	name := mustValue(call.Argument(0).Object().Get("name")).String()
	sync := boolOrFale(call.Argument(1))
	if sync {
		if err := self.client.StopContainer(name, 10); err != nil {
			self.vm.Throw("DockerError", err)
		}
		return otto.UndefinedValue()
	}

	task := self.getTask("stop", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		task.err = self.client.StopContainer(name, 10)

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) has_image(name string) bool {

	images, err := self.client.ListImages(dockerclient.ListImagesOptions{})
	if err != nil {
		return false
	}

	for _, i := range images {
		if name == i.ID {

			return true
		}
		for _, t := range i.RepoTags {
			index := strings.Index(t, ":")

			if t[:index] == name {

				return true
			}
		}
	}
	return false
}

func (self *docker_p) HasImage(call otto.FunctionCall) otto.Value {

	name := mustValue(call.Argument(0).Object().Get("name")).String()
	sync := boolOrFale(call.Argument(1))
	if sync {
		has := self.has_image(name)
		if has {
			return otto.TrueValue()
		} else {
			return otto.FalseValue()
		}
	}

	task := self.getTask("has_image", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		images, err := self.client.ListImages(dockerclient.ListImagesOptions{})
		if err != nil {
			task.err = err
			return
		}

		for _, i := range images {
			if name == i.ID {
				task.result = true
				return
			}
			for _, t := range i.RepoTags {
				index := strings.Index(t, ":")

				if t[:index] == name {
					task.result = true
					return
				}
			}
		}

		task.result = false

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) has_container(name string) bool {
	containers, err := self.client.ListContainers(dockerclient.ListContainersOptions{
		All: true,
	})
	if err != nil {

		return false
	}

	for _, i := range containers {
		if name == i.ID {

			return true
		}
		for _, t := range i.Names {
			if t[1:] == name || t == name {

				return true
			}
		}
	}
	return false
}

func (self *docker_p) HasContainer(call otto.FunctionCall) otto.Value {

	name := mustValue(call.Argument(0).Object().Get("name")).String()
	sync := boolOrFale(call.Argument(1))
	if sync {
		has := self.has_container(name)
		if has {
			return otto.TrueValue()
		} else {
			return otto.FalseValue()
		}
	}

	task := self.getTask("has_container", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		task.result = self.has_container(name)

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) RemoveContainer(call otto.FunctionCall) otto.Value {
	task := self.getTask("remove_container", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}
		o := call.Argument(0).Object()
		name := mustValue(o.Get("name")).String()
		force := boolOr(o, "force", false)
		removeVolumes := boolOr(o, "removeVolumes", false)
		task.err = self.client.RemoveContainer(dockerclient.RemoveContainerOptions{
			ID:            name,
			Force:         force,
			RemoveVolumes: removeVolumes,
		})

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) RemoveImage(call otto.FunctionCall) otto.Value {
	task := self.getTask("remove_image", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}
		o := call.Argument(0).Object()
		name := mustValue(o.Get("name")).String()
		force := boolOr(o, "force", false)
		prune := boolOr(o, "prune", true)
		task.err = self.client.RemoveImageExtended(name, dockerclient.RemoveImageOptions{
			Force:   force,
			NoPrune: !prune,
		})

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) IsRunning(call otto.FunctionCall) otto.Value {

	sync := boolOrFale(call.Argument(1))
	if sync {
		o := call.Argument(0).Object()
		name := mustValue(o.Get("name")).String()
		i, err := self.client.InspectContainer(name)
		if err != nil {
			return otto.FalseValue()
		} else if i.State.Running {
			return otto.TrueValue()
		}
		return otto.FalseValue()
	}

	task := self.getTask("is_running", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}
		o := call.Argument(0).Object()
		name := mustValue(o.Get("name")).String()

		i, err := self.client.InspectContainer(name)

		if err != nil {
			task.result = false
		} else {
			task.result = i.State.Running
		}

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) Inspect(call otto.FunctionCall) otto.Value {
	o := call.Argument(0).Object()
	name := mustValue(o.Get("name")).String()
	sync := boolOrFale(call.Argument(1))
	if sync {

		i, err := self.client.InspectContainer(name)
		if err != nil {
			return otto.NullValue()
		}

		if val, err := self.vm.ToValue(i); err != nil {
			self.vm.Throw("TypeError", err)
		} else {
			return val
		}
	}

	task := self.getTask("is_running", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		i, err := self.client.InspectContainer(name)

		if err != nil {
			task.result = false
		} else {
			task.result = i.State.Running
		}

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) ListContainers(call otto.FunctionCall) otto.Value {
	task := self.getTask("list_containers", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		containers, err := self.client.ListContainers(dockerclient.ListContainersOptions{})
		if err != nil {
			task.err = err
			return
		}

		task.result = containers

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) ListImages(call otto.FunctionCall) otto.Value {
	task := self.getTask("list_images", call.Argument(1))

	go func() {
		defer self.vm.Runloop().Ready(task)
		if task.err = self.check_args(call); task.err != nil {
			return
		}

		images, err := self.client.ListImages(dockerclient.ListImagesOptions{})
		if err != nil {
			task.err = err
			return
		}

		task.result = images

	}()

	return otto.UndefinedValue()
}

func (self *docker_p) Check(addr string, c int) bool {

	count := 0
	for {

		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return true
		}

		count++
		if count == c {
			return false
		}

		time.Sleep(1 * time.Second)
	}

	return false
}

type DockerOptions struct {
	Endpoint string
	Cert     string
	Key      string
	Ca       string
	Env      bool
}

func createDocker(vm *notto.Notto, o DockerOptions) (*docker_p, error) {
	var (
		c *dockerclient.Client
		e error
	)

	if o.Env {
		c, e = dockerclient.NewClientFromEnv()
	} else {
		if o.Cert == "" {
			c, e = dockerclient.NewClient(o.Endpoint)
		} else {
			c, e = dockerclient.NewTLSClient(o.Endpoint, o.Cert, o.Key, o.Ca)
		}
	}

	if e != nil {
		return nil, e
	}

	if c == nil {
		return nil, errors.New("no client")
	}

	return &docker_p{
		vm:     vm,
		client: c,
	}, nil
}
