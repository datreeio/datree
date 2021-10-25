package completion

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/datreeio/datree/bl/messager"
	"github.com/stretchr/testify/mock"
)

type printerMock struct {
	mock.Mock
}

type mockMessager struct {
	mock.Mock
}

func (p *printerMock) PrintMessage(messageText string, messageColor string) {
	p.Called(messageText, messageColor)
}

func (m *mockMessager) LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string) {
	go func() {
		messages <- &messager.VersionMessage{
			CliVersion:   "0.0.1",
			MessageText:  "message text",
			MessageColor: "White",
		}

		close(messages)
	}()

	m.Called(messages, cliVersion)
}

func Test_ExecuteCommand(t *testing.T) {
	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything, mock.Anything).Return(nil)

	printer := &printerMock{}
	printer.On("PrintMessage", mock.Anything, mock.Anything).Return(nil)

	cmd := New(&CompletionCommandContext{
		CliVersion: "1.2.3",
		Messager:   messager,
		Printer:    printer,
	})

	b := bytes.NewBufferString("bash") // Same for zsh, fish, powershell
	cmd.SetOut(b)
	cmd.SetArgs([]string{"bash"}) // Same for zsh, fish, powershell
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	// repalce bash with zsh, fish, powershell as per var b & SetArgs
	if string(out) != "bash" {
		t.Fatalf("expected \"%s\" got \"%s\"", "bash", string(out))
	}
}
