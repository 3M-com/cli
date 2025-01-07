package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/cli/cli/v2/pkg/cmd/attestation/test/data"
	"github.com/golang/snappy"
)

type mockHttpClient struct {
	mutex             sync.RWMutex
	currNumCalls      int
	failAfterNumCalls int
}

func (m *mockHttpClient) Get(url string) (*http.Response, error) {
	m.mutex.Lock()
	m.currNumCalls++
	m.mutex.Unlock()

	if m.failAfterNumCalls > 0 && m.currNumCalls >= m.failAfterNumCalls {
		return &http.Response{
			StatusCode: 500,
		}, fmt.Errorf("failed to fetch with %s", url)
	}

	var compressed []byte
	compressed = snappy.Encode(compressed, data.SigstoreBundleRaw)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(compressed)),
	}, nil
}
