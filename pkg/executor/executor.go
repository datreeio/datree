package executor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// CommandOutput is RunCommand result object
type CommandOutput struct {
	ResultOutput bytes.Buffer
	ErrorOutput  bytes.Buffer
}

type CommandRunner struct {
}

func CreateNewCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

// BuildCommandDescription executes Cmd strcuts and returns a human-readable description of it
func (c *CommandRunner) BuildCommandDescription(dir string, name string, args []string) string {
	command := exec.Command(name, args...)
	return command.String()
}

// RunCommand executes Cmd struct with given named program and arguments and returns CommandOutput
func (c *CommandRunner) RunCommand(name string, args []string) (CommandOutput, error) {
	var commandOutputBuffer, commandErrorBuffer bytes.Buffer

	command := exec.Command(name, args...)

	command.Stdout = &commandOutputBuffer
	command.Stderr = &commandErrorBuffer

	err := command.Run()
	if err != nil {
		return CommandOutput{
			ResultOutput: commandOutputBuffer,
			ErrorOutput:  commandErrorBuffer,
		}, fmt.Errorf("command output:%s, err:%s", commandOutputBuffer.String(), commandErrorBuffer.String())
	}

	return CommandOutput{
		ResultOutput: commandOutputBuffer,
		ErrorOutput:  commandErrorBuffer,
	}, nil
}

func (c *CommandRunner) ExecuteKustomizeBin(args []string) ([]byte, error) {
	if c.commandExists("kustomize") {

		commandOutput, err := c.RunCommand("kustomize", append([]string{"build"}, args...))
		if err != nil {
			return nil, fmt.Errorf("kustomize build errored: %s",
				commandOutput.ErrorOutput.String())
		}

		return commandOutput.ResultOutput.Bytes(), nil
	} else if c.commandExists("kubectl") {

		commandOutput, err := c.RunCommand("kubectl", append([]string{"kustomize"}, args...))
		if err != nil {
			return nil, fmt.Errorf("kubectl kustomize errored: %s",
				commandOutput.ErrorOutput.String())
		}

		return commandOutput.ResultOutput.Bytes(), nil
	} else {
		return nil, errors.New("kubectl or kustomize is not installed")
	}
}

func (c *CommandRunner) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (c *CommandRunner) CreateTempFile(tempFilePrefix string, content []byte) (string, error) {
	tempFile, err := os.CreateTemp("", fmt.Sprintf("%s_*.yaml", tempFilePrefix))
	if err != nil {
		return "", err
	}

	_, err = tempFile.Write(content)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}
