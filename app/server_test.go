package app

import (
	"github.com/ashtanko/octo-server/store/mockstore"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestStartServerSuccess(t *testing.T) {
	store := mockstore.New()
	s, err := NewServer(store)
	require.NoError(t, err)
	serverErr := s.Start(8000)
	client := &http.Client{}
	err = checkEndpoint(t, client, "http://localhost:"+"8000"+"/", http.StatusNotFound)
	require.NoError(t, err)
	s.Shutdown()
	require.NoError(t, serverErr)
}

func checkEndpoint(t *testing.T, client *http.Client, url string, expectedStatus int) error {
	res, err := client.Get(url)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != expectedStatus {
		t.Errorf("Response code was %d; want %d", res.StatusCode, expectedStatus)
	}

	return nil
}
