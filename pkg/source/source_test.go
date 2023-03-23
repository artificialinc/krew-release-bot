package source_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rajatjindal/krew-release-bot/pkg/source"
)

func TestDownloadFileRetry(t *testing.T) {
	retries := 0
	handler := http.NewServeMux()
	handler.HandleFunc("/rajatjindal/kubectl-whoami/releases/download/v0.0.2/kubectl-whoami_v0.0.2_darwin_amd64.tar.gz", func(w http.ResponseWriter, r *http.Request) {
		retries++
		w.WriteHeader(http.StatusNotFound)
	})

	srv := httptest.NewServer(handler)
	defer srv.Close()

	d := source.NewHttpDownloader()

	_, err := d.DownloadFileWithName(srv.URL+"/rajatjindal/kubectl-whoami/releases/download/v0.0.2/kubectl-whoami_v0.0.2_darwin_amd64.tar.gz", "whoami")
	assert.NotNil(t, err)
	assert.Equal(t, 4, retries)
}
