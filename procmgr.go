package phpstart

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// TODO can we do this without yml?
// TODO clean up code

// Procs is the existing list of process names and commands to run
type Procs struct {
	Processes map[string]Proc
}

// Proc is a single process to run
type Proc struct {
	Command string
	Args    []string
}

func NewProc(command string, args []string) Proc {
	return Proc{Command: command, Args: args}
}

func NewProcs() Procs {
	return Procs{
		Processes: map[string]Proc{},
	}
}

func (procs Procs) Add(procName string, newProc Proc) {
	procs.Processes[procName] = newProc
}

//TODO make this WriteFile
func (procs Procs) WriteFile(path string) error {
	bytes, err := yaml.Marshal(procs)
	if err != nil {
		//untested
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func ReadProcs(path string) (Procs, error) {
	procs := Procs{}

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return Procs{
			Processes: map[string]Proc{},
		}, nil
	} else if err != nil {
		return Procs{}, fmt.Errorf("failed to open proc.yml: %w", err)
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return Procs{}, err
	}

	err = yaml.UnmarshalStrict(contents, &procs)
	if err != nil {
		return Procs{}, fmt.Errorf("invalid proc.yml contents:\n %q: %w", contents, err)
	}

	return procs, nil
}
