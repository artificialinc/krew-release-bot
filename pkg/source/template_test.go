package source

import (
	"encoding/json"
	"testing"

	"io/ioutil"

	"github.com/google/go-github/v29/github"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestFixIndent(t *testing.T) {
	testcases := []struct {
		name   string
		input  string
		indent int
		output string
	}{
		{
			name: "indent with 4 spaces",
			input: `uri: some-secure-uri
sha256: some-sha256`,
			indent: 4,
			output: `uri: some-secure-uri
    sha256: some-sha256`,
		},
		{
			name: "fix legacy 4 space indent with 6 spaces",
			input: `uri: some-secure-uri
    sha256: some-sha256`,
			indent: 6,
			output: `uri: some-secure-uri
      sha256: some-sha256`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			output := indent(tc.indent, tc.input)
			assert.Equal(t, tc.output, output)
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	testcases := []struct {
		name     string
		file     string
		expected string
	}{
		{
			name:     "needs 6 space indentation",
			file:     "data/needs-6-space-indentation.yaml",
			expected: `data/needs-6-space-indentation-expected.yaml`,
		},
		{
			name:     "needs 4 space indentation",
			file:     "data/needs-4-space-indentation.yaml",
			expected: `data/needs-4-space-indentation-expected.yaml`,
		},
		{
			name:     "line start with dash",
			file:     "data/line-start-with-dash.yaml",
			expected: `data/line-start-with-dash-expected.yaml`,
		},
	}

	values := ReleaseRequest{
		TagName: "v0.0.2",
	}

	setup := func() {
		gock.New("https://github.com").
			Get("/rajatjindal/kubectl-whoami/releases/download/v0.0.2/kubectl-whoami_v0.0.2_darwin_amd64.tar.gz").
			Reply(200).
			BodyString("my-plugin-binary")
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer gock.Off()

			output, err := RenderTemplate(NewHttpDownloader(), tc.file, values)
			if err != nil {
				panic(err)
			}

			expectedOut, err := ioutil.ReadFile(tc.expected)
			if err != nil {
				panic(err)
			}

			assert.Equal(t, string(expectedOut), string(output))
		})
	}
}

func TestRenderTemplateGithubDownloader(t *testing.T) {
	testcases := []struct {
		name     string
		file     string
		expected string
		err      error
	}{
		{
			name:     "needs 6 space indentation",
			file:     "data/needs-6-space-indentation.yaml",
			expected: `data/needs-6-space-indentation-expected.yaml`,
		},
		{
			name:     "needs 4 space indentation",
			file:     "data/needs-4-space-indentation.yaml",
			expected: `data/needs-4-space-indentation-expected.yaml`,
		},
		{
			name:     "line start with dash",
			file:     "data/line-start-with-dash.yaml",
			expected: `data/line-start-with-dash-expected.yaml`,
		},
	}

	values := ReleaseRequest{
		TagName: "v0.0.2",
	}

	setup := func() {
		release := github.RepositoryRelease{
			Assets: []github.ReleaseAsset{
				{
					Name: github.String("kubectl-whoami_v0.0.2_darwin_amd64.tar.gz"),
					ID:   github.Int64(1),
				},
			},
		}
		b, err := json.Marshal(release)
		if err != nil {
			panic(err)
		}
		gock.New("https://api.github.com").
			Get("repos/rajatjindal/kubectl-whoami/releases/tags/v0.0.2").
			Reply(200).
			BodyString(string(b))

		gock.New("https://api.github.com").
			Get("/repos/rajatjindal/kubectl-whoami/releases/assets/1").
			Reply(200).
			BodyString("my-plugin-binary")
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			setup()
			defer gock.Off()

			output, err := RenderTemplate(NewGithubDownloader(), tc.file, values)
			assert.Equal(t, tc.err, err)

			expectedOut, err := ioutil.ReadFile(tc.expected)
			if err != nil {
				panic(err)
			}

			assert.Equal(t, string(expectedOut), string(output))
		})
	}
}
