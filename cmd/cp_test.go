package cmd

import (
	"testing"
)

func TestValidateJSON(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		wantErr bool
	}{
		"empty object": {
			input:   "{}",
			wantErr: false,
		},
		"single param": {
			input:   `{"path":"/home/user"}`,
			wantErr: false,
		},
		"multiple params": {
			input:   `{"name":"John","title":"Mr."}`,
			wantErr: false,
		},
		"with spaces": {
			input:   `{"path": "/home/user", "name": "John"}`,
			wantErr: false,
		},
		"empty string value": {
			input:   `{"key":""}`,
			wantErr: false,
		},
		"value with special chars": {
			input:   `{"url":"https://example.com:8080/path?query=value"}`,
			wantErr: false,
		},
		"invalid json - missing quotes": {
			input:   `{path:/home/user}`,
			wantErr: true,
		},
		"invalid json - trailing comma": {
			input:   `{"name":"John",}`,
			wantErr: true,
		},
		"invalid json - single quotes": {
			input:   `{'name':'John'}`,
			wantErr: true,
		},
		"invalid json - not an object": {
			input:   `["name","John"]`,
			wantErr: true,
		},
		"invalid json - unquoted value": {
			input:   `{"name":John}`,
			wantErr: true,
		},
		"empty string": {
			input:   "",
			wantErr: true,
		},
		"plain text": {
			input:   "not json",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := validateJSON(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("validateJSON() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("validateJSON() unexpected error: %v", err)
			}
		})
	}
}
