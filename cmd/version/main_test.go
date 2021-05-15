package version

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/datreeio/datree/bl/messager"
	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/stretchr/testify/mock"
)

type mockMessager struct {
	mock.Mock
}

func (m *mockMessager) PopulateVersionMessageChan(cliVersion string) <-chan *messager.VersionMessage {
	args := m.Called(cliVersion)
	return args.Get(0).(<-chan *messager.VersionMessage)
}

func (m *mockMessager) HandleVersionMessage(messageChannel <-chan *messager.VersionMessage) {
	m.Called(messageChannel)
}

func Test_ExecuteCommand(t *testing.T) {
	messager := &mockMessager{}
	messager.On("PopulateVersionMessageChan", mock.Anything).Return(mockedMessagesChannel())

	cmd := NewCommand(&VersionCommandContext{
		CliVersion: "1.2.3",
		Messager:   messager,
	})

	// TODO: ask yishay
	b := bytes.NewBufferString("1.2.3")
	cmd.SetOut(b)
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "1.2.3" {
		t.Fatalf("expected \"%s\" got \"%s\"", "1.2.3", string(out))
	}
}

func mockedMessagesChannel() <-chan *cliClient.VersionMessage {
	mock := make(chan *cliClient.VersionMessage)
	mock <- &cliClient.VersionMessage{
		CliVersion:   "1.2.3",
		MessageText:  "version message mock",
		MessageColor: "green"}
	close(mock)

	return mock
}
