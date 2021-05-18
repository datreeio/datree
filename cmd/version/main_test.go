package version

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/datreeio/datree/bl/messager"
	"github.com/stretchr/testify/mock"
)

type mockMessager struct {
	mock.Mock
}

func (m *mockMessager) LoadVersionMessages(messages chan *messager.VersionMessage, cliVersion string) {
	m.Called(messages, cliVersion)
}

func Test_ExecuteCommand(t *testing.T) {
	messager := &mockMessager{}

	cmd := New(&VersionCommandContext{
		CliVersion: "1.2.3",
		Messager:   messager,
	})

	// TODO: ask yishay wtf
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
