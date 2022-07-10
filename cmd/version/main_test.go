package version

import (
	"bytes"
	"io"
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

func (m *mockMessager) LoadVersionMessages(cliVersion string) chan *messager.VersionMessage {
	messages := make(chan *messager.VersionMessage, 1)
	go func() {
		messages <- &messager.VersionMessage{
			CliVersion:   "0.0.1",
			MessageText:  "message text",
			MessageColor: "Highlight",
		}

		close(messages)
	}()

	m.Called(cliVersion)
	return messages
}

func Test_ExecuteCommand(t *testing.T) {
	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything).Return(nil)

	printer := &printerMock{}
	printer.On("PrintMessage", mock.Anything, mock.Anything).Return(nil)

	cmd := New(&VersionCommandContext{
		CliVersion: "1.2.3",
		Messager:   messager,
		Printer:    printer,
	})

	b := bytes.NewBufferString("1.2.3")
	cmd.SetOut(b)
	cmd.Execute()
	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "1.2.3" {
		t.Fatalf("expected \"%s\" got \"%s\"", "1.2.3", string(out))
	}
}
