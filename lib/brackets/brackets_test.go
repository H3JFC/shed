package brackets

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

const (
	fortyCharVar    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	fortyOneCharVar = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
)

func init() {
	if len(fortyCharVar) != 40 {
		panic("fortyCharVar is not 40 characters long")
	}

	if len(fortyOneCharVar) != 41 {
		panic("fortyOneCharVar is not 41 characters long")
	}
}

func Test_parseBrackets(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input string
		want  []string
	}

	inputs := map[string]testcase{
		"standard-1": {
			input: "Hello, {{name}}! Welcome to {{place}}.",
			want:  []string{"name", "place"},
		},
		"standard-2": {
			input: "{{one}} some text {{two}} more text {{three}}",
			want:  []string{"one", "two", "three"},
		},
		"extra-spacing": {
			input: "{{ one }} some text {{two}} more text {{three}}",
			want:  []string{"one", "two", "three"},
		},
		"extra-spacing-with-desc": {
			input: "{{ one }} some text {{two | | description   }} more text {{three}}",
			want:  []string{"one", "two|| description", "three"},
		},
		"duplicates-1": {
			input: "{{one}} some text {{two}} more than {{one}} text {{three}}{{two}}",
			want:  []string{"one", "two", "three"},
		},
		"duplicates-2-fuller-description": {
			input: "{{one|foobar}} some text {{two|base}} more than {{one|foobarbaz}} text {{three|}}{{two}}",
			want:  []string{"one|foobarbaz", "two|base", "three|"},
		},
		"with-pipes": {
			input: "Hello {{world|foobar}} and {{universe|}}!",
			want:  []string{"world|foobar", "universe|"},
		},
		"just-brackets": {
			input: "{{first}}{{second}}{{third}}",
			want:  []string{"first", "second", "third"},
		},
		"no-blocks": {
			input: "No blocks here",
			want:  []string{},
		},
		"single-block": {
			input: "{{single_block}}",
			want:  []string{"single_block"},
		},
		"in-middle": {
			input: "Start {{middle}} end",
			want:  []string{"middle"},
		},
		"empty": {
			input: "{{}}",
			want:  []string{},
		},
		"empty-with-desc": {
			input: "{{|foo}}",
			want:  []string{},
		},
		"empty-with-empty-desc": {
			input: "{{|}}",
			want:  []string{},
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := parseBrackets(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("expected %d results, got %d, full obj: %v", len(tc.want), len(got), got)
			}

			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("at index %d, expected %q, got %q", i, tc.want[i], got[i])
				}
			}
		})
	}
}

func TestParseCommand(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input string
		want  string
	}

	inputs := map[string]testcase{
		"extra-space-at-ends": {
			input: "  Hello, {{name}}! Welcome to {{place}}.  ",
			want:  "Hello, {{name}}! Welcome to {{place}}.",
		},
		"extra-space-center": {
			input: "{{one}} some     text {{two}} more text {{three}}",
			want:  "{{one}} some text {{two}} more text {{three}}",
		},
		"extra-spacing-in-brackets": {
			input: "{{ one }} some text {{ two  }} more text {{three}}",
			want:  "{{one}} some text {{two}} more text {{three}}",
		},
		"extra-spacing-with-desc-1": {
			input: "{{ one | a normal description }} some text {{two}} more text {{three}}",
			want:  "{{one|a normal description}} some text {{two}} more text {{three}}",
		},
		"extra-spacing-with-desc-2": {
			input: "{{ one }} some text {{two | | description   }} more text {{three}}",
			want:  "{{one}} some text {{two|| description}} more text {{three}}",
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseCommand(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.want {
				t.Errorf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

func TestParseParameters_OK(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input string
		want  []Parameter
	}

	inputs := map[string]testcase{
		"standard-1": {
			input: "Hello, {{name}}! Welcome to {{place}}.",
			want:  []Parameter{{"name", ""}, {"place", ""}},
		},
		"standard-2": {
			input: "{{one}} some text {{two}} more text {{three}}",
			want:  []Parameter{{"one", ""}, {"two", ""}, {"three", ""}},
		},
		"character-limit-with-spaces": {
			input: "{{ " + fortyCharVar + " }} some text {{two}} more text {{three}}",
			want:  []Parameter{{fortyCharVar, ""}, {"two", ""}, {"three", ""}},
		},
		"extra-spacing": {
			input: "{{ one }} some text {{two}} more text {{three}}",
			want:  []Parameter{{"one", ""}, {"two", ""}, {"three", ""}},
		},
		"extra-spacing-with-desc": {
			input: "{{ one }} some text {{two | | description   }} more text {{three}}",
			want:  []Parameter{{"one", ""}, {"two", "| description"}, {"three", ""}},
		},
		"duplicates-1": {
			input: "{{one}} some text {{two}} more than {{one}} text {{three}}{{two}}",
			want:  []Parameter{{"one", ""}, {"two", ""}, {"three", ""}},
		},
		"duplicates-2-fuller-description": {
			input: "{{one|foobar}} some text {{two|base}} more than {{one|foobarbaz}} text {{three|}}{{two}}",
			want:  []Parameter{{"one", "foobarbaz"}, {"two", "base"}, {"three", ""}},
		},
		"with-pipes": {
			input: "Hello {{world|earth}} and {{universe|}}, {{universe2||}}!",
			want:  []Parameter{{"world", "earth"}, {"universe", ""}, {"universe2", "|"}},
		},
		"just-brackets": {
			input: "{{first}}{{second}}{{third}}",
			want:  []Parameter{{"first", ""}, {"second", ""}, {"third", ""}},
		},
		"no-blocks": {
			input: "No blocks here",
			want:  []Parameter{},
		},
		"single-block": {
			input: "{{single_block}}",
			want:  []Parameter{{"single_block", ""}},
		},
		"in-middle": {
			input: "Start {{middle}} end",
			want:  []Parameter{{"middle", ""}},
		},
		"empty": {
			input: "{{}}",
			want:  []Parameter{},
		},
		"empty-with-desc": {
			input: "{{|foo}}",
			want:  []Parameter{},
		},
		"empty-with-empty-desc": {
			input: "{{|}}",
			want:  []Parameter{},
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseParameters(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("expected %d results, got %d; obj: %v", len(tc.want), len(got), got)
			}

			for i := range got {
				if got[i].Name != tc.want[i].Name {
					t.Errorf("at index %d, expected name: %q got name: %q", i, tc.want[i].Name, got[i].Name)
				}

				if got[i].Description != tc.want[i].Description {
					t.Errorf("at index %d, expected description: %q got description: %q", i,
						tc.want[i].Description, got[i].Description)
				}
			}
		})
	}
}

func TestParseParameters_Err(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input string
		want  error
	}

	inputs := map[string]testcase{
		"exceeds-limit": {
			input: "Hello, {{" + fortyOneCharVar + "}}! Welcome to {{place}}.",
			want:  ErrParameterTooLong,
		},
		"space-in-param": {
			input: "Hello, {{foo bar}}! Welcome to {{place}}.",
			want:  ErrParameterContainsSpaces,
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseParameters(tc.input)
			if err == nil {
				t.Fatalf("expected an error, got nil")
			}

			if len(got) != 0 {
				t.Fatalf("expected empty result, got: %v", got)
			}

			if !errors.Is(err, tc.want) {
				t.Errorf("expected error: %v, got error: %v", tc.want, err)
			}
		})
	}
}

func TestParseParameters_ErrSymbol(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input string
		want  error
	}

	inputs := map[string]testcase{}

	nums := "0123456789"
	syms := "!@#$%^&*()-+=[]{};:'\",.<>?/\\`~"

	for _, ch := range nums {
		inputs["starts-invalid-char-"+string(ch)] = testcase{
			input: "Hello, {{ " + string(ch) + "valid_suffix }}! Welcome to {{place}}.",
			want:  ErrParameterStartsWithInvalidChar,
		}
	}

	for _, ch := range syms {
		inputs["contains-invalid-char-"+string(ch)] = testcase{
			input: "Hello, {{ " + "valid_" + string(ch) + "valid_suffix }}! Welcome to {{place}}.",
			want:  ErrContainsInvalidSymbols,
		}
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseParameters(tc.input)
			if len(got) != 0 {
				t.Fatalf("expected empty result, got: %v", got)
			}

			if err == nil {
				t.Fatalf("expected an error, got nil")
			}

			if !errors.Is(err, tc.want) {
				t.Errorf("expected error: %v, got error: %v", tc.want, err)
			}
		})
	}
}

func TestHydrateStringSafe(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input  string
		params ValuedParameters
		want   string
	}

	inputs := map[string]testcase{
		"standard-1": {
			input: "Hello, {{name}}! Welcome to {{place}}.",
			params: ValuedParameters{
				{"name", "Alice"},
				{"place", "Wonderland"},
			},
			want: "Hello, Alice! Welcome to Wonderland.",
		},
		"standard-2": {
			input: "{{one}} some text {{two}} more text {{three}}",
			params: ValuedParameters{
				{"one", "1"},
				{"three", "three"},
			},
			want: "1 some text {{two}} more text three",
		},
		"with-pipes": {
			input: "Hello {{world|earth}}, {{universe|}} and {{universe2||}}!",
			params: ValuedParameters{
				{"world", "globe"},
				{"universe", "cosmos"},
				{"universe2", "reality"},
			},
			want: "Hello globe, cosmos and reality!",
		},
		"just-brackets": {
			input: "{{first}}{{second}}{{third}}",
			params: ValuedParameters{
				{"first", "first"},
				{"second", "second"},
				{"third", "third"},
			},
			want: "firstsecondthird",
		},
		"no-blocks": {
			input: "No blocks here",
			params: ValuedParameters{
				{"unused", "value"},
			},
			want: "No blocks here",
		},
		"single-block": {
			input: "{{single_block}}",
			params: ValuedParameters{
				{"single_block", "just one"},
				{"unused", "value"},
			},
			want: "just one",
		},
		"in-middle": {
			input: "Start {{middle}} end",
			params: ValuedParameters{
				{"middle", "center"},
			},
			want: "Start center end",
		},
		"empty": {
			input: "{{}}",
			params: ValuedParameters{
				{"unused", "value"},
			},
			want: "",
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := HydrateStringSafe(tc.input, tc.params)
			if got != tc.want {
				t.Errorf("expected string: %q got string: %q", tc.want, got)
			}
		})
	}
}

func TestHydrateString_NoErr(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input  string
		params ValuedParameters
		want   string
	}

	inputs := map[string]testcase{
		"standard-1": {
			input: "Hello, {{name}}! Welcome to {{place}}.",
			params: ValuedParameters{
				{"name", "Alice"},
				{"place", "Wonderland"},
			},
			want: "Hello, Alice! Welcome to Wonderland.",
		},
		"with-pipes": {
			input: "Hello {{world|earth}}, {{universe|}} and {{universe2||}}!",
			params: ValuedParameters{
				{"world", "globe"},
				{"universe", "cosmos"},
				{"universe2", "reality"},
			},
			want: "Hello globe, cosmos and reality!",
		},
		"just-brackets": {
			input: "{{first}}{{second}}{{third}}",
			params: ValuedParameters{
				{"first", "first"},
				{"second", "second"},
				{"third", "third"},
			},
			want: "firstsecondthird",
		},
		"no-blocks": {
			input: "No blocks here",
			params: ValuedParameters{
				{"unused", "value"},
			},
			want: "No blocks here",
		},
		"single-block": {
			input: "{{single_block}}",
			params: ValuedParameters{
				{"single_block", "just one"},
				{"unused", "value"},
			},
			want: "just one",
		},
		"in-middle": {
			input: "Start {{middle}} end",
			params: ValuedParameters{
				{"middle", "center"},
			},
			want: "Start center end",
		},
		"empty": {
			input: "{{}}",
			params: ValuedParameters{
				{"unused", "value"},
			},
			want: "",
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := HydrateString(tc.input, tc.params)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.want {
				t.Errorf("expected string: %q got string: %q", tc.want, got)
			}
		})
	}
}

func TestHydrateString_Err(t *testing.T) {
	t.Parallel()

	type testcase struct {
		input  string
		params ValuedParameters
		want   error
	}

	inputs := map[string]testcase{
		"standard-2": {
			input: "{{one}} some text {{two}} more text {{three}}",
			params: ValuedParameters{
				{"one", "1"},
				{"three", "three"},
			},
			want: ErrMissingParameters,
		},
	}

	for name, tc := range inputs {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := HydrateString(tc.input, tc.params)
			if got != "" {
				t.Fatalf("expected empty string, got: %q", got)
			}

			if err == nil {
				t.Fatalf("expected an error, got: %v", err)
			}

			if !errors.Is(err, tc.want) {
				t.Errorf("expected error: %v, got error: %v", tc.want, err)
			}
		})
	}
}

func TestHydrateStringFromJSON_NoErr(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input           string
		jsonValueParams string
		want            string
	}

	tests := map[string]testcase{
		"standard": {
			input:           "Hello, {{name}}! Welcome to {{place}}.",
			jsonValueParams: `{"name":"Alice","place":"Wonderland"}`,
			want:            "Hello, Alice! Welcome to Wonderland.",
		},
		"with-description-pipes": {
			input:           "Hello {{world|earth}}, {{universe|cosmos}} and {{universe2|reality}}!",
			jsonValueParams: `{"world":"globe","universe":"multiverse","universe2":"dimension"}`,
			want:            "Hello globe, multiverse and dimension!",
		},
		"single-parameter": {
			input:           "curl -XGET {{url}}",
			jsonValueParams: `{"url":"https://example.com"}`,
			want:            "curl -XGET https://example.com",
		},
		"multiple-parameters": {
			input: "curl -XPOST --data '{{data}}' -H {{header}} {{url}}",
			jsonValueParams: `{"data":"{\"foo\":\"bar\"}","header":"Authorization:Bearer token",` +
				`"url":"https://api.example.com"}`,
			want: "curl -XPOST --data '{\"foo\":\"bar\"}' -H Authorization:Bearer token https://api.example.com",
		},
		"empty-json": {
			input:           "No parameters here",
			jsonValueParams: `{}`,
			want:            "No parameters here",
		},
		"unused-parameters": {
			input:           "Only {{used}} parameter",
			jsonValueParams: `{"used":"this one","unused":"ignored"}`,
			want:            "Only this one parameter",
		},
		"special-characters": {
			input:           "Path: {{path}}",
			jsonValueParams: `{"path":"/home/user/my documents/file.txt"}`,
			want:            "Path: /home/user/my documents/file.txt",
		},
		"empty-value": {
			input:           "Value: {{empty}}",
			jsonValueParams: `{"empty":""}`,
			want:            "Value: ",
		},
		"whitespace-in-json": {
			input:           "{{first}} {{second}}",
			jsonValueParams: `{ "first" : "one" , "second" : "two" }`,
			want:            "one two",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := HydrateStringFromJSON(tc.input, tc.jsonValueParams)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.want {
				t.Errorf("expected: %q, got: %q", tc.want, got)
			}
		})
	}
}

func TestHydrateStringFromJSON_WithEmptyString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input string
		want  string
	}{
		"with-parameters": {
			input: "curl -XGET {{url}} -H {{header}}",
			want:  "curl -XGET {{url}} -H {{header}}",
		},
		"no-parameters": {
			input: "echo hello world",
			want:  "echo hello world",
		},
		"with-description-pipes": {
			input: "{{name|user name}} says {{message|greeting message}}",
			want:  "{{name|user name}} says {{message|greeting message}}",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Empty string should work the same as "{}"
			got, err := HydrateStringFromJSON(tc.input, "")
			if err != nil {
				t.Fatalf("unexpected error with empty string: %v", err)
			}

			if got != tc.want {
				t.Errorf("expected: %q, got: %q", tc.want, got)
			}

			// Verify it's the same as using "{}"
			gotWithEmptyJSON, err := HydrateStringFromJSON(tc.input, "{}")
			if err != nil {
				t.Fatalf("unexpected error with empty JSON: %v", err)
			}

			if got != gotWithEmptyJSON {
				t.Errorf("empty string and empty JSON should produce same result: %q vs %q", got, gotWithEmptyJSON)
			}
		})
	}
}

func TestHydrateStringFromJSON_Err(t *testing.T) { //nolint:funlen
	t.Parallel()

	type testcase struct {
		input           string
		jsonValueParams string
		wantErr         error
	}

	tests := map[string]testcase{
		"invalid-json-syntax": {
			input:           "{{param}}",
			jsonValueParams: `{invalid}`,
			wantErr:         ErrParsingValueParams,
		},
		"invalid-json-missing-quotes": {
			input:           "{{param}}",
			jsonValueParams: `{param:value}`,
			wantErr:         ErrParsingValueParams,
		},
		"invalid-json-trailing-comma": {
			input:           "{{param}}",
			jsonValueParams: `{"param":"value",}`,
			wantErr:         ErrParsingValueParams,
		},
		"invalid-json-single-quotes": {
			input:           "{{param}}",
			jsonValueParams: `{'param':'value'}`,
			wantErr:         ErrParsingValueParams,
		},
		"invalid-json-array": {
			input:           "{{param}}",
			jsonValueParams: `["param","value"]`,
			wantErr:         ErrParsingValueParams,
		},
		"plain-text": {
			input:           "{{param}}",
			jsonValueParams: "not json",
			wantErr:         ErrParsingValueParams,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := HydrateStringFromJSON(tc.input, tc.jsonValueParams)
			if got != "" {
				t.Fatalf("expected empty string, got: %q", got)
			}

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected error: %v, got error: %v", tc.wantErr, err)
			}
		})
	}
}

func TestParameters_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		params   Parameters
		expected string
	}{
		{
			name: "sorted output",
			params: Parameters{
				{Name: "zebra", Description: "last"},
				{Name: "apple", Description: "first"},
				{Name: "monkey", Description: "middle"},
			},
			expected: `[{"name":"apple","description":"first"},{"name":"monkey","description":"middle"},{"name":"zebra","description":"last"}]`, //nolint:lll
		},
		{
			name:     "empty slice",
			params:   Parameters{},
			expected: `[]`,
		},
		{
			name:     "nil slice",
			params:   nil,
			expected: `[]`,
		},
		{
			name: "single element",
			params: Parameters{
				{Name: "single", Description: "only one"},
			},
			expected: `[{"name":"single","description":"only one"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(got) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestParameters_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		json     string
		expected Parameters
		wantErr  bool
	}{
		{
			name: "unsorted input gets sorted",
			json: `[{"name":"zebra","description":"last"},{"name":"apple","description":"first"}]`,
			expected: Parameters{
				{Name: "apple", Description: "first"},
				{Name: "zebra", Description: "last"},
			},
		},
		{
			name:     "empty array",
			json:     `[]`,
			expected: Parameters{},
		},
		{
			name:     "null",
			json:     `null`,
			expected: nil,
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got Parameters

			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("Unmarshal() length = %d, want %d", len(got), len(tt.expected))
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Unmarshal()[%d] = %+v, want %+v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestParameters_DeterministicMarshaling(t *testing.T) {
	t.Parallel()

	// Same data in different orders should produce identical JSON
	params1 := Parameters{
		{Name: "z", Description: "3"},
		{Name: "a", Description: "1"},
		{Name: "m", Description: "2"},
	}
	params2 := Parameters{
		{Name: "a", Description: "1"},
		{Name: "m", Description: "2"},
		{Name: "z", Description: "3"},
	}

	json1, err1 := json.Marshal(params1)
	json2, err2 := json.Marshal(params2)

	if err1 != nil || err2 != nil {
		t.Fatalf("Marshal errors: %v, %v", err1, err2)
	}

	if string(json1) != string(json2) {
		t.Errorf("Marshaling not deterministic:\n%s\n%s", json1, json2)
	}
}

func TestValuedParameters_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		params   ValuedParameters
		expected string
	}{
		{
			name: "sorted output",
			params: ValuedParameters{
				{Name: "zebra", Value: "last"},
				{Name: "apple", Value: "first"},
				{Name: "monkey", Value: "middle"},
			},
			expected: `[{"name":"apple","value":"first"},{"name":"monkey","value":"middle"},{"name":"zebra","value":"last"}]`,
		},
		{
			name:     "empty slice",
			params:   ValuedParameters{},
			expected: `[]`,
		},
		{
			name:     "nil slice",
			params:   nil,
			expected: `[]`,
		},
		{
			name: "single element",
			params: ValuedParameters{
				{Name: "single", Value: "only one"},
			},
			expected: `[{"name":"single","value":"only one"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if string(got) != tt.expected {
				t.Errorf("Marshal() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestValuedParameters_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		json     string
		expected ValuedParameters
		wantErr  bool
	}{
		{
			name: "unsorted input gets sorted",
			json: `[{"name":"zebra","value":"last"},{"name":"apple","value":"first"}]`,
			expected: ValuedParameters{
				{Name: "apple", Value: "first"},
				{Name: "zebra", Value: "last"},
			},
		},
		{
			name:     "empty array",
			json:     `[]`,
			expected: ValuedParameters{},
		},
		{
			name:     "null",
			json:     `[]`,
			expected: ValuedParameters{},
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got ValuedParameters

			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("Unmarshal() length = %d, want %d", len(got), len(tt.expected))
			}

			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Unmarshal()[%d] = %+v, want %+v", i, got[i], tt.expected[i])
				}
			}
		})
	}
}

func TestValuedParameters_DeterministicMarshaling(t *testing.T) {
	t.Parallel()

	// Same data in different orders should produce identical JSON
	params1 := ValuedParameters{
		{Name: "z", Value: "3"},
		{Name: "a", Value: "1"},
		{Name: "m", Value: "2"},
	}
	params2 := ValuedParameters{
		{Name: "a", Value: "1"},
		{Name: "m", Value: "2"},
		{Name: "z", Value: "3"},
	}

	json1, err1 := json.Marshal(params1)
	json2, err2 := json.Marshal(params2)

	if err1 != nil || err2 != nil {
		t.Fatalf("Marshal errors: %v, %v", err1, err2)
	}

	if string(json1) != string(json2) {
		t.Errorf("Marshaling not deterministic:\n%s\n%s", json1, json2)
	}
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("Parameters", func(t *testing.T) {
		t.Parallel()

		original := Parameters{
			{Name: "z", Description: "last"},
			{Name: "a", Description: "first"},
			{Name: "m", Description: "middle"},
		}

		// Marshal
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		// Unmarshal
		var unmarshaled Parameters
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		// Marshal again
		data2, err := json.Marshal(unmarshaled)
		if err != nil {
			t.Fatalf("Second marshal error: %v", err)
		}

		// Should be identical
		if string(data) != string(data2) {
			t.Errorf("Round trip not deterministic:\n%s\n%s", data, data2)
		}
	})

	t.Run("ValuedParameters", func(t *testing.T) {
		t.Parallel()

		original := ValuedParameters{
			{Name: "z", Value: "last"},
			{Name: "a", Value: "first"},
			{Name: "m", Value: "middle"},
		}

		// Marshal
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}

		// Unmarshal
		var unmarshaled ValuedParameters
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}

		// Marshal again
		data2, err := json.Marshal(unmarshaled)
		if err != nil {
			t.Fatalf("Second marshal error: %v", err)
		}

		// Should be identical
		if string(data) != string(data2) {
			t.Errorf("Round trip not deterministic:\n%s\n%s", data, data2)
		}
	})
}

func TestParameters_ToMap_Basic(t *testing.T) {
	t.Parallel()

	p := Parameters{
		{Name: "param1", Description: "desc1"},
		{Name: "param2", Description: "desc2"},
	}

	got := p.ToMap()
	want := map[string]string{
		"param1": "desc1",
		"param2": "desc2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToMap() = %v, want %v", got, want)
	}
}

func TestParameters_ToMap_Empty(t *testing.T) {
	t.Parallel()

	var p Parameters // nil slice

	got := p.ToMap()
	if len(got) != 0 {
		t.Errorf("ToMap() expected empty map, got %v", got)
	}

	p = Parameters{} // empty but non-nil

	got = p.ToMap()
	if len(got) != 0 {
		t.Errorf("ToMap() expected empty map, got %v", got)
	}
}

func TestParameters_Names_Basic(t *testing.T) {
	t.Parallel()

	p := Parameters{
		{Name: "z", Description: "last"},
		{Name: "a", Description: "first"},
	}

	got := p.Names()
	want := []string{"z", "a"} // preserves original order, not sorted

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestParameters_Names_Empty(t *testing.T) {
	t.Parallel()

	var p Parameters

	got := p.Names()
	want := []string{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestParameters_Description_Found(t *testing.T) {
	t.Parallel()

	p := Parameters{
		{Name: "a", Description: "first"},
		{Name: "b", Description: "second"},
	}

	desc, err := p.Description("a")
	if err != nil {
		t.Fatalf("Description() returned unexpected error: %v", err)
	}

	if desc != "first" { // nolint:goconst
		t.Errorf("Description() = %q, want %q", desc, "first")
	}
}

func TestParameters_Description_NotFound(t *testing.T) {
	t.Parallel()

	p := Parameters{
		{Name: "a", Description: "first"},
	}

	_, err := p.Description("missing")
	if err == nil {
		t.Fatal("Description() expected error, got nil")
	}

	if !errors.Is(err, ErrParameterNotFound) {
		t.Fatalf("Description() error = %v, want ErrParameterNotFound", err)
	}
}

func TestParameters_Replace(t *testing.T) { // nolint:funlen,gocognit,cyclop
	t.Parallel()

	t.Run("replace-existing-parameter", func(t *testing.T) {
		t.Parallel()

		params := Parameters{
			{Name: "alpha", Description: "first"},
			{Name: "beta", Description: "second"},
			{Name: "gamma", Description: "third"},
		}

		params.Replace("beta", "updated second")

		if len(params) != 3 {
			t.Errorf("expected length 3, got %d", len(params))
		}

		desc, err := params.Description("beta")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "updated second" {
			t.Errorf("expected 'updated second', got '%s'", desc)
		}
	})

	t.Run("append-new-parameter", func(t *testing.T) {
		t.Parallel()

		params := Parameters{
			{Name: "alpha", Description: "first"},
			{Name: "gamma", Description: "third"},
		}

		params.Replace("beta", "new second")

		if len(params) != 3 {
			t.Errorf("expected length 3, got %d", len(params))
		}

		desc, err := params.Description("beta")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "new second" {
			t.Errorf("expected 'new second', got '%s'", desc)
		}
	})

	t.Run("maintains-sorted-order-after-replace", func(t *testing.T) {
		t.Parallel()

		params := Parameters{
			{Name: "alpha", Description: "first"},
			{Name: "beta", Description: "second"},
		}

		params.Replace("beta", "updated")

		names := params.Names()
		if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
			t.Errorf("expected sorted order [alpha, beta], got %v", names)
		}
	})

	t.Run("maintains-sorted-order-after-append", func(t *testing.T) {
		t.Parallel()

		params := Parameters{
			{Name: "alpha", Description: "first"},
			{Name: "gamma", Description: "third"},
		}

		params.Replace("beta", "new second")

		names := params.Names()
		expected := []string{"alpha", "beta", "gamma"}

		if len(names) != 3 {
			t.Fatalf("expected length 3, got %d", len(names))
		}

		for i, name := range expected {
			if names[i] != name {
				t.Errorf("at index %d: expected '%s', got '%s'", i, name, names[i])
			}
		}
	})

	t.Run("replace-on-empty-slice", func(t *testing.T) {
		t.Parallel()

		params := Parameters{}

		params.Replace("alpha", "first")

		if len(params) != 1 {
			t.Errorf("expected length 1, got %d", len(params))
		}

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "first" {
			t.Errorf("expected 'first', got '%s'", desc)
		}
	})

	t.Run("replace on nil slice", func(t *testing.T) {
		t.Parallel()

		var params Parameters

		params.Replace("alpha", "first")

		if len(params) != 1 {
			t.Errorf("expected length 1, got %d", len(params))
		}

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "first" {
			t.Errorf("expected 'first', got '%s'", desc)
		}
	})

	t.Run("multiple-replaces-maintain-order", func(t *testing.T) {
		t.Parallel()

		params := Parameters{
			{Name: "delta", Description: "fourth"},
		}

		params.Replace("alpha", "first")
		params.Replace("gamma", "third")
		params.Replace("beta", "second")
		params.Replace("alpha", "updated first")

		names := params.Names()
		expected := []string{"alpha", "beta", "delta", "gamma"}

		if len(names) != 4 {
			t.Fatalf("expected length 4, got %d", len(names))
		}

		for i, name := range expected {
			if names[i] != name {
				t.Errorf("at index %d: expected '%s', got '%s'", i, name, names[i])
			}
		}

		desc, _ := params.Description("alpha")
		if desc != "updated first" {
			t.Errorf("expected 'updated first', got '%s'", desc)
		}
	})

	t.Run("replace with empty description", func(t *testing.T) {
		t.Parallel()

		params := Parameters{
			{Name: "alpha", Description: "first"},
		}

		params.Replace("alpha", "")

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "" {
			t.Errorf("expected empty string, got '%s'", desc)
		}
	})
}

func TestParameters_MergeName(t *testing.T) { // nolint:funlen,cyclop,gocognit,maintidx
	t.Parallel()

	t.Run("merge-existing-parameter-updates-description", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
			{Name: "beta", Description: "second"},
		}
		other := &Parameters{
			{Name: "alpha", Description: "updated first"},
			{Name: "gamma", Description: "third"},
		}

		params.MergeName(other, "alpha")

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "updated first" {
			t.Errorf("expected 'updated first', got '%s'", desc)
		}
	})

	t.Run("merge-non-existing-parameter-adds-it", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
		}
		other := &Parameters{
			{Name: "beta", Description: "second"},
		}

		params.MergeName(other, "beta")

		if len(*params) != 2 {
			t.Errorf("expected length 2, got %d", len(*params))
		}

		desc, err := params.Description("beta")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "second" {
			t.Errorf("expected 'second', got '%s'", desc)
		}
	})

	t.Run("merge-name-not-in-other-does-nothing", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
		}
		other := &Parameters{
			{Name: "beta", Description: "second"},
		}

		params.MergeName(other, "gamma")

		if len(*params) != 1 {
			t.Errorf("expected length 1, got %d", len(*params))
		}

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "first" {
			t.Errorf("expected 'first', got '%s'", desc)
		}
	})

	t.Run("nil-receiver-does-nothing", func(t *testing.T) {
		t.Parallel()

		var params *Parameters

		other := &Parameters{
			{Name: "alpha", Description: "first"},
		}

		// Should not panic
		params.MergeName(other, "alpha")

		if params != nil {
			t.Error("expected params to remain nil")
		}
	})

	t.Run("nil-other-does-nothing", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
		}

		var other *Parameters

		params.MergeName(other, "alpha")

		if len(*params) != 1 {
			t.Errorf("expected length 1, got %d", len(*params))
		}

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "first" {
			t.Errorf("expected 'first', got '%s'", desc)
		}
	})

	t.Run("merge-into-empty-parameters", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{}
		other := &Parameters{
			{Name: "alpha", Description: "first"},
		}

		params.MergeName(other, "alpha")

		if len(*params) != 1 {
			t.Errorf("expected length 1, got %d", len(*params))
		}

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "first" {
			t.Errorf("expected 'first', got '%s'", desc)
		}
	})

	t.Run("merge-from-empty-other-does-nothing", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
		}
		other := &Parameters{}

		params.MergeName(other, "alpha")

		if len(*params) != 1 {
			t.Errorf("expected length 1, got %d", len(*params))
		}

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "first" {
			t.Errorf("expected 'first', got '%s'", desc)
		}
	})

	t.Run("maintains-sorted-order-after-merge", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
			{Name: "delta", Description: "fourth"},
		}
		other := &Parameters{
			{Name: "beta", Description: "second"},
		}

		params.MergeName(other, "beta")

		names := params.Names()
		expected := []string{"alpha", "beta", "delta"}

		if len(names) != 3 {
			t.Fatalf("expected length 3, got %d", len(names))
		}

		for i, name := range expected {
			if names[i] != name {
				t.Errorf("at index %d: expected '%s', got '%s'", i, name, names[i])
			}
		}
	})

	t.Run("merge-with-empty-description", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
		}
		other := &Parameters{
			{Name: "alpha", Description: ""},
		}

		params.MergeName(other, "alpha")

		desc, err := params.Description("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if desc != "" {
			t.Errorf("expected empty string, got '%s'", desc)
		}
	})

	t.Run("does-not-affect-other-parameters", func(t *testing.T) {
		t.Parallel()

		params := &Parameters{
			{Name: "alpha", Description: "first"},
			{Name: "beta", Description: "second"},
			{Name: "gamma", Description: "third"},
		}
		other := &Parameters{
			{Name: "beta", Description: "updated second"},
		}

		params.MergeName(other, "beta")

		// Check that beta was updated
		desc, _ := params.Description("beta")
		if desc != "updated second" {
			t.Errorf("expected 'updated second' for beta, got '%s'", desc)
		}

		// Check that alpha was not changed
		desc, _ = params.Description("alpha")
		if desc != "first" {
			t.Errorf("expected 'first' for alpha, got '%s'", desc)
		}

		// Check that gamma was not changed
		desc, _ = params.Description("gamma")
		if desc != "third" {
			t.Errorf("expected 'third' for gamma, got '%s'", desc)
		}
	})
}

func TestThreeWayMerge_NewParameterTakesPriority(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "priority desc"},
	}
	before := &Parameters{}
	updated := &Parameters{
		{Name: "alpha", Description: "short"},
	}

	priority.ThreeWayMerge(before, updated)

	desc, _ := priority.Description("alpha")
	if desc != "priority desc" {
		t.Errorf("expected 'priority desc', got '%s'", desc)
	}
}

func TestThreeWayMerge_NewParameterTakesUpdated(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "short"},
	}
	before := &Parameters{}
	updated := &Parameters{
		{Name: "alpha", Description: "updated longer desc"},
	}

	priority.ThreeWayMerge(before, updated)

	desc, _ := priority.Description("alpha")
	if desc != "updated longer desc" {
		t.Errorf("expected 'updated longer desc', got '%s'", desc)
	}
}

func TestThreeWayMerge_BothChangedTakesLonger(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "priority changed"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "original"},
	}
	updated := &Parameters{
		{Name: "alpha", Description: "updated changed longer"},
	}

	priority.ThreeWayMerge(before, updated)

	desc, _ := priority.Description("alpha")
	if desc != "updated changed longer" {
		t.Errorf("expected 'updated changed longer', got '%s'", desc)
	}
}

func TestThreeWayMerge_BothChangedTakesPriority(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "priority changed longer"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "original"},
	}
	updated := &Parameters{
		{Name: "alpha", Description: "updated short"},
	}

	priority.ThreeWayMerge(before, updated)

	desc, _ := priority.Description("alpha")
	if desc != "priority changed longer" {
		t.Errorf("expected 'priority changed longer', got '%s'", desc)
	}
}

func TestThreeWayMerge_OnlyPriorityChanged(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "priority changed"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "original"},
	}
	updated := &Parameters{
		{Name: "alpha", Description: "original"},
	}

	priority.ThreeWayMerge(before, updated)

	desc, _ := priority.Description("alpha")
	if desc != "priority changed" { //nolint:goconst
		t.Errorf("expected 'priority changed', got '%s'", desc)
	}
}

func TestThreeWayMerge_OnlyUpdatedChanged(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "original"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "original"},
	}
	updated := &Parameters{
		{Name: "alpha", Description: "updated changed"},
	}

	priority.ThreeWayMerge(before, updated)

	desc, _ := priority.Description("alpha")
	if desc != "updated changed" {
		t.Errorf("expected 'updated changed', got '%s'", desc)
	}
}

func TestThreeWayMerge_NothingChanged(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "same"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "same"},
	}
	updated := &Parameters{
		{Name: "alpha", Description: "same"},
	}

	priority.ThreeWayMerge(before, updated)

	desc, _ := priority.Description("alpha")
	if desc != "same" {
		t.Errorf("expected 'same', got '%s'", desc)
	}
}

func TestThreeWayMerge_NilPriority(t *testing.T) {
	t.Parallel()

	var priority *Parameters

	before := &Parameters{{Name: "alpha", Description: "before"}}
	updated := &Parameters{{Name: "alpha", Description: "updated"}}

	// Should not panic
	priority.ThreeWayMerge(before, updated)

	if priority != nil {
		t.Error("expected priority to remain nil")
	}
}

func TestThreeWayMerge_NilBefore(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "priority"},
	}

	var before *Parameters

	updated := &Parameters{
		{Name: "alpha", Description: "updated longer"},
	}

	priority.ThreeWayMerge(before, updated)

	// Treats as new parameter since before is nil
	desc, _ := priority.Description("alpha")
	if desc != "updated longer" {
		t.Errorf("expected 'updated longer', got '%s'", desc)
	}
}

func TestThreeWayMerge_NilUpdated(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "priority changed"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "original"},
	}

	var updated *Parameters

	priority.ThreeWayMerge(before, updated)

	// Priority changed but updated doesn't exist, keep priority
	desc, _ := priority.Description("alpha")
	if desc != "priority changed" {
		t.Errorf("expected 'priority changed', got '%s'", desc)
	}
}

func TestThreeWayMerge_MultipleParameters(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "priority changed"},
		{Name: "beta", Description: "same"},
		{Name: "gamma", Description: "original"},
		{Name: "delta", Description: "new"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "original"},
		{Name: "beta", Description: "same"},
		{Name: "gamma", Description: "original"},
	}
	updated := &Parameters{
		{Name: "alpha", Description: "original"},
		{Name: "beta", Description: "updated changed"},
		{Name: "gamma", Description: "updated changed longer"},
		{Name: "delta", Description: "new but longer description"},
	}

	priority.ThreeWayMerge(before, updated)

	// alpha: priority changed, updated didn't -> keep priority
	desc, _ := priority.Description("alpha")
	if desc != "priority changed" {
		t.Errorf("alpha: expected 'priority changed', got '%s'", desc)
	}

	// beta: priority didn't change, updated changed -> take updated
	desc, _ = priority.Description("beta")
	if desc != "updated changed" {
		t.Errorf("beta: expected 'updated changed', got '%s'", desc)
	}

	// gamma: both changed, updated is longer -> take updated
	desc, _ = priority.Description("gamma")
	if desc != "updated changed longer" {
		t.Errorf("gamma: expected 'updated changed longer', got '%s'", desc)
	}

	// delta: new parameter, updated is longer -> take updated
	desc, _ = priority.Description("delta")
	if desc != "new but longer description" {
		t.Errorf("delta: expected 'new but longer description', got '%s'", desc)
	}
}

func TestThreeWayMerge_UpdatedMissingParameter(t *testing.T) {
	t.Parallel()

	priority := &Parameters{
		{Name: "alpha", Description: "changed"},
	}
	before := &Parameters{
		{Name: "alpha", Description: "original"},
	}
	updated := &Parameters{}

	priority.ThreeWayMerge(before, updated)

	// Priority changed but updated doesn't have it -> keep priority
	desc, _ := priority.Description("alpha")
	if desc != "changed" {
		t.Errorf("expected 'changed', got '%s'", desc)
	}
}

func TestThreeWayMerge_EmptyPriority(t *testing.T) {
	t.Parallel()

	priority := &Parameters{}
	before := &Parameters{{Name: "alpha", Description: "before"}}
	updated := &Parameters{{Name: "alpha", Description: "updated"}}

	priority.ThreeWayMerge(before, updated)

	// Priority is empty, nothing to merge
	if len(*priority) != 0 {
		t.Errorf("expected empty priority, got %d items", len(*priority))
	}
}
