package docs

import (
	"bytes"
	"io"
	"testing"

	"github.com/datreeio/datree/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type openBrowserMock struct {
	mock.Mock
}

func (o *openBrowserMock) OpenURL(url string) error {
	args := o.Called(url)
	return args.Error(0)
}

func Test_DocsCommand(t *testing.T) {
	testCases := []struct {
		desc     string
		url      string
		expected []byte
	}{
		{
			desc:     "Test official URL",
			url:      "https://hub.datree.io",
			expected: []byte(`Opening https://hub.datree.io in your browser.`),
		},
		{
			desc:     "Test empty URL",
			url:      "",
			expected: []byte(`Opening https://hub.datree.io in your browser.`),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			openBrowser := &openBrowserMock{}
			openBrowser.On("OpenURL", mock.Anything).Return(nil)

			cmd := New(&DocsCommandContext{
				BrowserCtx: utils.OpenBrowserContext{
					UrlOpener: openBrowser,
				},
				URL: tC.url,
			})

			out := bytes.NewBufferString("Opening https://hub.datree.io in your browser.")
			cmd.SetOut(out)
			err := cmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			got, err := io.ReadAll(out)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, string(tC.expected), string(got))
		})
	}
}
