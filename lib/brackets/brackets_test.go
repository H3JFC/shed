package brackets

import (
	"encoding/json"
	"errors"
	"testing"
)

const (
	fortyCharVar    = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	fortyOneCharVar = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
)

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
			expected: `[{"name":"apple","description":"first"},{"name":"monkey","description":"middle"},{"name":"zebra","description":"last"}]`,
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
