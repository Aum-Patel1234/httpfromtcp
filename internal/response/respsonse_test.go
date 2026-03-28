package response_test

import (
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/response"
	"io"
	"testing"
)

func TestWriteHeaders(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w       io.Writer
		headers headers.Headers
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := response.WriteHeaders(tt.w, tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("WriteHeaders() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("WriteHeaders() succeeded unexpectedly")
			}
		})
	}
}
