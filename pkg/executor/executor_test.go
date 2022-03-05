package executor

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestCommandRunner_RunCommand(t *testing.T) {
	type args struct {
		name string
		args []string
	}
	tests := []struct {
		name    string
		c       *CommandRunner
		args    args
		want    CommandOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandRunner{}
			got, err := c.RunCommand(tt.args.name, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandRunner.RunCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandRunner.RunCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

type commandRunnerMock struct {
	mock.Mock
}

func (c *commandRunnerMock) RunCommand(name string, arguments []string) (CommandOutput, error) {
	args := c.Called(name, arguments)
	return args.Get(0).(CommandOutput), args.Error(1)
}
func TestCommandRunner_ExecuteKustomizeBin(t *testing.T) {

	commandRunner := new(commandRunnerMock)
	commandRunner.On("RunCommand", "kustomize", []string{"build"}).Return(CommandOutput{})

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		c       *CommandRunner
		args    args
		want    []byte
		wantErr bool
		err     error
	}{
		{
			name: "should return error if kustomize or kubectl is not installed",
			c:    &CommandRunner{},
			args: args{
				args: []string{"bad_arg"},
			},
			want:    nil,
			wantErr: true,
			err:     errors.New("kubectl or kustomize is not installed"),
		},
		{
			name: "should return error if kustomize or kubectl is not installed",
			c:    CreateNewCommandRunner(),
			args: args{
				args: []string{"bad_arg"},
			},
			want:    nil,
			wantErr: true,
			err:     errors.New("kubectl or kustomize is not installed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandRunner{}
			got, err := c.ExecuteKustomizeBin(tt.args.args)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CommandRunner.ExecuteKustomizeBin() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if !reflect.DeepEqual(err, tt.err) {
					t.Errorf("CommandRunner.ExecuteKustomizeBin() = err %v, want %v", err, tt.err)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandRunner.ExecuteKustomizeBin() = %v, want %v", got, tt.want)
			}
		})
	}
}
