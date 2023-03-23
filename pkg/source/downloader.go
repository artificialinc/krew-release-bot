package source

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/v29/github"
	"github.com/sirupsen/logrus"
)

const retries = 4

type Downloader interface {
	DownloadFileWithName(uri, name string) (string, error)
}

type HttpDownloader struct {
	client *http.Client
}

// WithClient sets the http client
func (h *HttpDownloader) WithClient(client *http.Client) *HttpDownloader {
	h.client = client
	return h
}

// NewHttpDownloader creates a new HttpDownloader
func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{
		client: http.DefaultClient,
	}
}

// DownloadFileWithName downloads a file with name
func (h *HttpDownloader) DownloadFileWithName(uri, name string) (string, error) {
	resp, err := getWithRetry(uri, h.client)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("downloading file %s failed. status code: %d, expected: %d", uri, resp.StatusCode, http.StatusOK)
	}
	return writeTempFile(resp.Body, name)
}

type GithubDownloader struct {
	client *http.Client
}

// WithClient sets the http client
func (g *GithubDownloader) WithClient(client *http.Client) *GithubDownloader {
	g.client = client
	return g
}

// NewGithubDownloader creates a new GithubDownloader
func NewGithubDownloader() *GithubDownloader {
	return &GithubDownloader{
		client: http.DefaultClient,
	}
}

// DownloadFileWithName downloads a file with name using github client
func (g *GithubDownloader) DownloadFileWithName(uri, name string) (string, error) {
	client := github.NewClient(g.client)

	owner, repo, tag, asset, err := parseGithubURI(uri)
	if err != nil {
		return "", err
	}

	release, _, err := client.Repositories.GetReleaseByTag(context.TODO(), owner, repo, tag)
	if err != nil {
		return "", err
	}
	for _, a := range release.Assets {
		if *a.Name == asset {
			r, _, err := client.Repositories.DownloadReleaseAsset(context.TODO(), owner, repo, *a.ID, http.DefaultClient)
			if err != nil {
				return "", err
			}
			return writeTempFile(r, name)
		}
	}
	return "", fmt.Errorf("failed to find asset %s in release %s", asset, tag)
}

func parseGithubURI(uri string) (string, string, string, string, error) {
	// Example: https://github.com/artificialinc/artificial-kubectl-plugins/releases/download/v0.0.1/kubectl-redis-v0.0.1-darwin-amd64.tar.gz
	r := regexp.MustCompile(`https://github.com/([^:\/\s]+)/([^:\/\s]+)/releases/download/([^:\/\s]+)/([^:\/\s]+)`)
	matches := r.FindStringSubmatch(uri)
	if len(matches) != 5 {
		return "", "", "", "", fmt.Errorf("failed to parse github uri %s", uri)
	}
	return matches[1], matches[2], matches[3], matches[4], nil
}

func writeTempFile(r io.ReadCloser, name string) (string, error) {
	defer r.Close()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	file := filepath.Join(dir, name)
	out, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return "", fmt.Errorf("failed to save file %s. error: %v", file, err)
	}

	logrus.Infof("downloaded file %s", file)
	return file, nil
}
