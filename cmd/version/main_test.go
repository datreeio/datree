package version

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func Test_ExecuteCommand(t *testing.T) {
	cmd := NewVersionCommand("1.2.3")
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
