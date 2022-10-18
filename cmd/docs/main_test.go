package docs

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DocsCommand(t *testing.T) {
	testCases := []struct {
		desc     string
		expected []byte
	}{
		{
			desc:     "Test datree docs",
			expected: []byte(`Opening https://hub.datree.io in your browser.`),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmd := New()
			out := bytes.NewBufferString("Opening https://hub.datree.io in your browser.")
			cmd.SetOut(out)
			cmd.Execute()

			got, err := ioutil.ReadAll(out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, string(tC.expected), string(got))
		})
	}
}
