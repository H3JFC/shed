package brackets

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"

	"h3jfc/shed/lib/itertools"
)

var (
	ErrMissingParameters      = errors.New("missing parameters")
	ErrNameEmpty              = errors.New("name cannot be empty")
	ErrTooLong                = errors.New("too long")
	ErrStartsWithInvalidChar  = errors.New("starts with invalid character")
	ErrContainsSpaces         = errors.New("contains spaces")
	ErrContainsInvalidSymbols = errors.New("contains invalid symbols")
	ErrParameterNotFound      = errors.New("parameter not found")
	ErrParsingValueParams     = errors.New("failed to parse value parameters")
)

var spaceRegex = regexp.MustCompile(`\s+`)

const (
	maxParts       = 2
	characterLimit = 40
	symbols        = "!@#$%^&*()-+=[]{};:'\",.<>?/\\|`~"
	bang           = '!'
)

var symbolSet map[rune]struct{}

func init() {
	symbolSet = make(map[rune]struct{})
	for _, r := range symbols {
		symbolSet[r] = struct{}{}
	}
}

type Parameter struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ValuedParameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Secret struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type (
	Parameters       []Parameter
	ValuedParameters []ValuedParameter // nolint:recvcheck
	Secrets          []Secret
)

type Brackets struct {
	Command    string
	Parameters *Parameters
	Secrets    *Secrets
}

// MarshalJSON ensures deterministic ordering by name.
func (p Parameters) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("[]"), nil
	}

	// Create a copy to avoid modifying the original
	sorted := make([]Parameter, len(p))
	copy(sorted, p)

	// Sort by name
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	// Marshal the sorted slice.
	return json.Marshal(sorted)
}

// UnmarshalJSON ensures the slice is sorted after unmarshaling.
func (p *Parameters) UnmarshalJSON(data []byte) error {
	var params []Parameter
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	// Sort by name
	sort.Slice(params, func(i, j int) bool {
		return params[i].Name < params[j].Name
	})

	*p = params

	return nil
}

func (p *Parameters) ToMap() map[string]string {
	m := make(map[string]string, len(*p))

	for i := range *p {
		m[(*p)[i].Name] = (*p)[i].Description
	}

	return m
}

func (p *Parameters) Names() []string {
	names := make([]string, 0, len(*p))

	for i := range *p {
		names = append(names, (*p)[i].Name)
	}

	return names
}

func (p *Parameters) Description(name string) (string, error) {
	m := p.ToMap()

	if desc, exists := m[name]; exists {
		return desc, nil
	}

	return "", fmt.Errorf("%w: %s", ErrParameterNotFound, name)
}

// Replace updates the description of an existing parameter or appends a new one if it doesn't exist.
func (p *Parameters) Replace(name, description string) {
	// Search for existing parameter
	for i := range *p {
		if (*p)[i].Name == name {
			(*p)[i].Description = description

			return
		}
	}

	// Parameter not found, append new one
	*p = append(*p, Parameter{
		Name:        name,
		Description: description,
	})

	// Re-sort to maintain deterministic ordering
	sort.Slice(*p, func(i, j int) bool {
		return (*p)[i].Name < (*p)[j].Name
	})
}

func (p *Parameters) MergeName(other *Parameters, name string) {
	if p == nil || other == nil {
		return
	}

	otherMap := other.ToMap()

	desc, exists := otherMap[name]
	if !exists {
		return
	}

	p.Replace(name, desc)
}

func (p *Parameters) ThreeWayMerge(before, updated *Parameters) { // nolint:cyclop
	if p == nil {
		return
	}

	beforeMap := make(map[string]string)
	if before != nil {
		beforeMap = before.ToMap()
	}

	updatedMap := make(map[string]string)
	if updated != nil {
		updatedMap = updated.ToMap()
	}

	for _, name := range p.Names() {
		priorityDesc, _ := p.Description(name)
		beforeDesc, existedBefore := beforeMap[name]
		updatedDesc, existsInUpdated := updatedMap[name]

		// New parameter: take longer of updated or priority
		if !existedBefore {
			if existsInUpdated && len(updatedDesc) > len(priorityDesc) {
				p.Replace(name, updatedDesc)
			}

			continue
		}

		// Existing parameter: check what changed
		priorityChanged := priorityDesc != beforeDesc
		updatedChanged := existsInUpdated && updatedDesc != beforeDesc

		if priorityChanged && updatedChanged {
			// Both changed: take longer
			if len(updatedDesc) > len(priorityDesc) {
				p.Replace(name, updatedDesc)
			}
		} else if updatedChanged {
			// Only updated changed: take updated
			p.Replace(name, updatedDesc)
		}
	}
}

// MarshalJSON ensures deterministic ordering by name.
func (vp ValuedParameters) MarshalJSON() ([]byte, error) {
	if vp == nil {
		return []byte("[]"), nil
	}

	// Create a copy to avoid modifying the original
	sorted := make([]ValuedParameter, len(vp))
	copy(sorted, vp)

	// Sort by name
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	// Marshal the sorted slice
	return json.Marshal(sorted)
}

// UnmarshalJSON ensures the slice is sorted after unmarshaling.
func (vp *ValuedParameters) UnmarshalJSON(data []byte) error {
	var params []ValuedParameter
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}

	// Sort by name
	sort.Slice(params, func(i, j int) bool {
		return params[i].Name < params[j].Name
	})

	*vp = params

	return nil
}

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

func ParametersFromMap(m map[string]string) Parameters {
	vp := make(Parameters, 0, len(m))

	for k, v := range m {
		vp = append(vp, Parameter{Name: k, Description: v})
	}

	return vp
}

func ParametersFromJSON(jsonStr string) (Parameters, error) {
	var m map[string]string

	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return nil, err
	}

	return ParametersFromMap(m), nil
}

func ParseParameters(input string) (Parameters, error) {
	predicate := func(p Parameter) bool {
		return rune(p.Name[0]) != bang // filter out secrets
	}

	pp := parseParamOrSecret(input, predicate)

	var err error

	pp, err = checkForInvalidParameters(pp, true)
	if err != nil {
		return nil, err
	}

	return pp, nil
}

func ParseSecrets(input string) (Secrets, error) {
	predicate := func(p Parameter) bool {
		return rune(p.Name[0]) == bang // filter out non-secrets
	}

	pp := parseParamOrSecret(input, predicate)

	pp = slices.Collect(itertools.Map(slices.Values(pp), func(p Parameter) Parameter {
		// Remove leading '!' from secret names
		return Parameter{
			Name:        p.Name[1:],
			Description: p.Description,
		}
	}))

	var err error

	pp, err = checkForInvalidParameters(pp, false)
	if err != nil {
		return nil, err
	}

	ss := itertools.Map(slices.Values(pp), func(p Parameter) Secret {
		return Secret(p)
	})

	return slices.Collect(ss), nil
}

// ParseCommand normalizes a command string by:
// - Trimming leading/trailing whitespace
// - Normalizing spacing inside {{...}} blocks
// - Normalizing spacing around | separators in parameter descriptions
// - Collapsing multiple spaces outside {{...}} blocks to single spaces.
func ParseCommand(input string) (string, error) {
	s := strings.TrimSpace(input)

	var result strings.Builder
	result.Grow(len(s))

	var outsideBrackets strings.Builder

	i := 0
	for i < len(s) {
		// Look for opening {{
		if i < len(s)-1 && s[i] == '{' && s[i+1] == '{' {
			// Flush any accumulated outside content
			if outsideBrackets.Len() > 0 {
				normalized := spaceRegex.ReplaceAllString(outsideBrackets.String(), " ")
				result.WriteString(normalized)
				outsideBrackets.Reset()
			}

			result.WriteString("{{")

			i += 2
			start := i

			// Find closing }}
			for i < len(s)-1 {
				if s[i] == '}' && s[i+1] == '}' {
					content := s[start:i]
					normalized := cleanString(content)
					result.WriteString(normalized)
					result.WriteString("}}")

					i += 2

					break
				}

				i++
			}
		} else {
			outsideBrackets.WriteByte(s[i])
			i++
		}
	}

	// Flush any remaining outside content
	if outsideBrackets.Len() > 0 {
		normalized := spaceRegex.ReplaceAllString(outsideBrackets.String(), " ")
		result.WriteString(normalized)
	}

	return result.String(), nil
}

func Parse(input string) (*Brackets, error) {
	cmd, err := ParseCommand(input)
	if err != nil {
		return nil, err
	}

	p, err := ParseParameters(input)
	if err != nil {
		return nil, err
	}

	s, err := ParseSecrets(input)
	if err != nil {
		return nil, err
	}

	return &Brackets{
		Command:    cmd,
		Parameters: &p,
		Secrets:    &s,
	}, err
}

func HydrateString(input string, vp ValuedParameters) (string, error) {
	out := HydrateStringSafe(input, vp)

	p, err := ParseParameters(input)
	if err != nil {
		return "", err
	}

	missing := vp.MissingSubset(p)
	if len(missing) > 0 {
		missingNames := itertools.Map(slices.Values(missing), func(param Parameter) string {
			return param.Name
		})

		return "", fmt.Errorf("%w: %v", ErrMissingParameters, slices.Collect(missingNames))
	}

	return out, nil
}

func HydrateStringSafe(s string, vp ValuedParameters) string {
	var out string

	var args []any

	i := 0

	j := 0

	var outSb268 strings.Builder

	for j < len(s)-1 {
		// Look for opening {{
		if s[j] == '{' && s[j+1] == '{' { //nolint:nestif
			outSb268.WriteString(s[i:j] + "%s")
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

	out += outSb268.String()

	if i < len(s) {
		out += s[i:]
	}

	return fmt.Sprintf(out, args...)
}

func HydrateStringFromJSON(cmd, jsonValueParams string) (string, error) {
	if jsonValueParams == "" {
		jsonValueParams = "{}"
	}

	vp, err := ValuedParametersFromJSON(jsonValueParams)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrParsingValueParams, err)
	}

	return HydrateStringSafe(cmd, vp), nil
}

func parseBrackets(s string) []string { //nolint:gocognit
	var results []string

	seen := make(map[string]int) // maps key to index in results

	i := 0
	for i < len(s)-1 {
		// Look for opening {{
		if s[i] == '{' && s[i+1] == '{' { //nolint:nestif
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

func parseParamOrSecret(input string, predicate func(Parameter) bool) Parameters {
	ss := parseBrackets(input)

	params := itertools.Map(slices.Values(ss), func(s string) Parameter {
		parts := strings.SplitN(s, "|", maxParts)

		if len(parts) == 1 {
			return Parameter{Name: parts[0], Description: ""}
		}

		return Parameter{Name: parts[0], Description: parts[1]}
	})

	pp := Parameters(slices.Collect(params))

	return itertools.Filter(pp, predicate)
}

func checkForInvalidParameters(pp Parameters, parameter bool) (Parameters, error) {
	errString := "parameter"
	if !parameter {
		errString = "secret"
	}

	for i := range pp {
		if len(pp[i].Name) == 0 {
			return nil, fmt.Errorf("%s %w", errString, ErrNameEmpty)
		}

		firstChar := rune(pp[i].Name[0])
		if firstChar >= '0' && firstChar <= '9' {
			return nil, fmt.Errorf("%s %w: %s", errString, ErrStartsWithInvalidChar, pp[i].Name)
		}

		for _, r := range pp[i].Name {
			if _, exists := symbolSet[r]; exists {
				return nil, fmt.Errorf("%s %w: %s", errString, ErrContainsInvalidSymbols, pp[i].Name)
			}
		}

		if len(pp[i].Name) > characterLimit {
			return nil, fmt.Errorf("%s %w: %s", errString, ErrTooLong, pp[i].Name)
		}

		if strings.Contains(pp[i].Name, " ") {
			return nil, fmt.Errorf("%s %w: %s", errString, ErrContainsSpaces, pp[i].Name)
		}
	}

	return pp, nil
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
