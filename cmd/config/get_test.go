package config

import (
	"bytes"
	"io"
	"testing"

	"github.com/datreeio/datree/pkg/localConfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const TOKEN_VALUE = "tokenValue"

func (lc *LocalConfigMock) Get(key string) string {
	lc.Called(key)
	return TOKEN_VALUE
}

func TestGetCommand(t *testing.T) {
	messager := &mockMessager{}
	messager.On("LoadVersionMessages", mock.Anything, mock.Anything)

	printerMock := &PrinterMock{}
	printerMock.On("PrintMessage", mock.Anything, mock.Anything)

	localConfigMock := &LocalConfigMock{}
	localConfigMock.On("GetLocalConfiguration").Return(&localConfig.LocalConfig{Token: "previousToken"}, nil)
	localConfigMock.On("Get", mock.Anything)
	ctx := &ConfigCommandContext{
		Messager:    messager,
		CliVersion:  "1.2.3",
		Printer:     printerMock,
		LocalConfig: localConfigMock,
	}

	testGetTokenCommand(t, ctx, localConfigMock)
	testExecuteGetCommand(t, ctx)
}

func testGetTokenCommand(t *testing.T, ctx *ConfigCommandContext, localConfigMock *LocalConfigMock) {
	err := get(ctx, "token")
	assert.Equal(t, nil, err)
	localConfigMock.AssertCalled(t, "Get", "token")
}

func testExecuteGetCommand(t *testing.T, ctx *ConfigCommandContext) {
	cmd := NewGetCommand(ctx)
	cmd.SetArgs([]string{"token"})
	b := bytes.NewBufferString(TOKEN_VALUE)
	cmd.SetOut(b)
	err := cmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != TOKEN_VALUE {
		t.Fatalf("expected \"%s\" got \"%s\"", TOKEN_VALUE, string(out))
	}
}
