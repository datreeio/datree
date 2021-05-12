package version

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/datreeio/datree/pkg/cliClient"
	"github.com/stretchr/testify/mock"
)

type mockVersionMessageClient struct {
	mock.Mock
}

func (m *mockVersionMessageClient) GetVersionMessage(cliVersion string) (*cliClient.VersionMessage, error) {
	args := m.Called(cliVersion)
	return args.Get(0).(*cliClient.VersionMessage), nil
}

func Test_ExecuteCommand(t *testing.T) {
	versionMessageClient := &mockVersionMessageClient{}
	versionMessageClient.On("GetVersionMessage", mock.Anything).Return(
		&cliClient.VersionMessage{
			CliVersion:   "1.2.3",
			MessageText:  "version message mock",
			MessageColor: "green"},
	)

	cmd := NewVersionCommand(&VersionCommandContext{
		CliVersion:           "1.2.3",
		VersionMessageClient: versionMessageClient,
	})
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
