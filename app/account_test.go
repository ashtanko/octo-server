package app

import (
	"bytes"
	"encoding/json"
	"github.com/ashtanko/octo-server/store/mockstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleAccountCreate(t *testing.T) {
	store := mockstore.New()
	s, err := NewServer(store)
	require.NoError(t, err)
	serverErr := s.Start(8080)
	require.NoError(t, serverErr)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]float64{
				"account_balance": 0,
			},
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			err := json.NewEncoder(b).Encode(tc.payload)
			require.NoError(t, err)
			req, _ := http.NewRequest(http.MethodPost, "/account/create", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
