package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
)

var colors = []color.Attribute{
	color.FgBlue,
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgCyan,
	color.FgMagenta,
	color.FgWhite,
}

func main() {
	shell := flag.String("x", "bash", "Default shell when using commands via arguments")
	file := flag.String("f", "parx.yml", "Path to the parx.yml file")
	flag.Parse()

	rawArgs := flag.Args()
	var config Config
	if len(rawArgs) == 0 {
		f, err := os.Open(*file)
		if err != nil {
			fmt.Printf("Could not read config file: %v", err)
			os.Exit(1)
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Printf("Could not read config file: %v", err)
			os.Exit(1)
		}

		err = yaml.Unmarshal(data, &config)
		if err != nil {
			fmt.Printf("Could not unmarshal config data to struct: %v", err)
			os.Exit(1)
		}
	} else {
		// parse from arguments, name and prefix is auto generated
		var processes []*Process
		for i, arg := range rawArgs {
			p := &Process{
				Name:    fmt.Sprintf("process_%d", i),
				Shell:   *shell,
				Command: arg,
				Env:     make(map[string]string),
			}
			processes = append(processes, p)
		}
		config.Processes = processes
	}

	wg := sync.WaitGroup{}
	for i, process := range config.Processes {
		// generate color for this process
		c := colors[i%len(colors)]
		process.color = c

		wg.Add(1)
		go func(process *Process) {
			defer wg.Done()
			cmd := process.ToExecCommand()

			err := cmd.Run()
			if err != nil {
				fmt.Printf("%s exited with an error: %v\n", process.Name, err)
			} else {
				fmt.Printf("%s exited\n", process.Name)
			}
		}(process)
	}
	wg.Wait()
}

type Config struct {
	Processes []*Process `json:"processes"`
}

type Process struct {
	Name    string            `json:"name"`
	Shell   string            `json:"shell"`
	Command string            `json:"command"`
	Env     map[string]string `json:"env"`

	color color.Attribute
	cmd   *exec.Cmd
}

func (p *Process) ToExecCommand() *exec.Cmd {
	if p.cmd != nil {
		return p.cmd
	}
	if p.Shell == "" {
		p.Shell = "bash"
	}

	p.cmd = exec.Command(p.Shell, "-c", p.Command)
	p.cmd.Env = os.Environ()
	for key, val := range p.Env {
		p.cmd.Env = append(p.cmd.Env, fmt.Sprintf("%s=%s", key, val))
	}

	prefix := fmt.Sprintf("%s | ", p.Name)
	p.cmd.Stdout = NewPrefixedWriter(prefix, p.color, os.Stdout)
	p.cmd.Stderr = NewPrefixedWriter(prefix, p.color, os.Stderr)

	return p.cmd
}

type PrefixedWriter struct {
	Prefix string

	colorize func(...interface{}) string
	target   io.Writer
	buf      bytes.Buffer
}

func NewPrefixedWriter(prefix string, prefixColor color.Attribute, target io.Writer) *PrefixedWriter {
	return &PrefixedWriter{
		Prefix:   prefix,
		colorize: color.New(prefixColor).SprintFunc(),
		target:   target,
	}
}

func (p *PrefixedWriter) Write(payload []byte) (int, error) {
	p.buf.Reset()

	p.buf.WriteString(p.colorize(p.Prefix))
	p.buf.Write(payload)

	n, err := p.target.Write(p.buf.Bytes())
	if err != nil {
		if n < len(payload) {
			n = len(payload)
		}
		return n, err
	}
	return len(payload), nil
}
