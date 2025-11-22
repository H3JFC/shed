package brackets

import (
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"slices"
	"strings"
)

var (
	ErrMissingParameters              = errors.New("missing parameters")
	ErrParameterTooLong               = errors.New("parameter too long")
	ErrParameterStartsWithInvalidChar = errors.New("parameter starts with invalid character")
	ErrParameterContainsSpaces        = errors.New("parameter contains spaces")
	ErrContainsInvalidSymbols         = errors.New("parameter contains invalid symbols")
)

const (
	maxParts       = 2
	characterLimit = 40
	symbols        = "!@#$%^&*()-+=[]{};:'\",.<>?/\\|`~"
)

var symbolSet map[rune]struct{}

func init() {
	symbolSet = make(map[rune]struct{})
	for _, r := range symbols {
		symbolSet[r] = struct{}{}
	}
}

type Parameter struct {
	Name        string
	Description string
}

type ValuedParameter struct {
	Name  string
	Value string
}

type (
	Parameters       []Parameter
	ValuedParameters []ValuedParameter
)

func (vp ValuedParameters) Value(name string) (string, bool) {
	for _, p := range vp {
		if p.Name == name {
			return p.Value, true
		}
	}

	return "", false
}

func (vp ValuedParameters) MissingSubset(p Parameters) Parameters {
	var missing Parameters

	for i := range p {
		n := p[i].Name
		if _, exists := vp.Value(n); !exists {
			missing = append(missing, p[i])
		}
	}

	return missing
}

func ValuedParametersFromMap(m map[string]string) ValuedParameters {
	vp := make(ValuedParameters, 0, len(m))

	for k, v := range m {
		vp = append(vp, ValuedParameter{Name: k, Value: v})
	}

	return vp
}

func ValuedParametersFromJSON(jsonStr string) (ValuedParameters, error) {
	var m map[string]string

	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return ValuedParametersFromMap(m), nil
}

func ParseParameters(input string) (Parameters, error) {
	ss := parseBrackets(input)

	params := Map(slices.Values(ss), func(s string) Parameter {
		parts := strings.SplitN(s, "|", maxParts)

		if len(parts) == 1 {
			return Parameter{Name: parts[0], Description: ""}
		}

		return Parameter{Name: parts[0], Description: parts[1]}
	})

	pp := slices.Collect(params)

	for i := range pp {
		if len(pp[i].Name) == 0 {
			panic("parameter name cannot be empty")
		}

		firstChar := rune(pp[i].Name[0])
		if firstChar >= '0' && firstChar <= '9' {
			return nil, fmt.Errorf("%w: %s", ErrParameterStartsWithInvalidChar, pp[i].Name)
		}

		for _, r := range pp[i].Name {
			if _, exists := symbolSet[r]; exists {
				return nil, fmt.Errorf("%w: %s", ErrContainsInvalidSymbols, pp[i].Name)
			}
		}

		if len(pp[i].Name) > characterLimit {
			return nil, fmt.Errorf("%w: %s", ErrParameterTooLong, pp[i].Name)
		}

		if strings.Contains(pp[i].Name, " ") {
			return nil, fmt.Errorf("%w: %s", ErrParameterContainsSpaces, pp[i].Name)
		}
	}

	return pp, nil
}

func HydrateString(input string, vp ValuedParameters) (string, error) {
	out := HydrateStringSafe(input, vp)

	p, err := ParseParameters(input)
	if err != nil {
		return "", err
	}

	missing := vp.MissingSubset(p)
	if len(missing) > 0 {
		missingNames := Map(slices.Values(missing), func(param Parameter) string {
			return param.Name
		})

		return "", fmt.Errorf("%w: %v", ErrMissingParameters, missingNames)
	}

	return out, nil
}

func HydrateStringSafe(s string, vp ValuedParameters) string {
	var out string

	var args []any

	i := 0

	j := 0
	for j < len(s)-1 {
		// Look for opening {{
		if s[j] == '{' && s[j+1] == '{' {
			out += s[i:j] + "%s"
			j += 2
			start := j

			// Find closing }}
			for j < len(s)-1 {
				if s[j] == '}' && s[j+1] == '}' {
					content := cleanString(s[start:j])
					name := parseName(content)

					if len(name) == 0 {
						args = append(args, "")
					} else {
						if val, exists := vp.Value(name); exists {
							args = append(args, val)
						} else {
							args = append(args, "{{"+content+"}}")
						}
					}

					j += 2
					i = j

					break
				}

				j++
			}
		} else {
			j++
		}
	}

	if i < len(s) {
		out += s[i:]
	}

	return fmt.Sprintf(out, args...)
}

func parseBrackets(s string) []string {
	var results []string

	seen := make(map[string]int) // maps key to index in results

	i := 0
	for i < len(s)-1 {
		// Look for opening {{
		if s[i] == '{' && s[i+1] == '{' {
			i += 2
			start := i

			// Find closing }}
			for i < len(s)-1 {
				if s[i] == '}' && s[i+1] == '}' {
					content := cleanString(s[start:i])
					name := parseName(content)

					if len(name) > 0 {
						if idx, exists := seen[name]; exists {
							// Replace if current content is longer
							if len(content) > len(results[idx]) {
								results[idx] = content
							}
						} else {
							// Add new entry
							seen[name] = len(results)
							results = append(results, content)
						}
					}

					i += 2

					break
				}

				i++
			}
		} else {
			i++
		}
	}

	return results
}

func parseName(s string) string {
	parts := strings.SplitN(s, "|", maxParts)

	return strings.TrimSpace(parts[0])
}

func cleanString(s string) string {
	parts := strings.SplitN(s, "|", maxParts)

	if len(parts) == 1 {
		return strings.TrimSpace(parts[0])
	}

	return strings.TrimSpace(parts[0]) + "|" + strings.TrimSpace(parts[1])
}

func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}
