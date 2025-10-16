package os

type (
	OS int
)

const (
	Darwin OS = iota
	Linux
	Windows
	Unknown
)

const (
	linux   string = "linux"
	darwin  string = "darwin"
	windows string = "windows"
	unknown string = "unknown"
)

var toString = map[OS]string{
	Darwin:  darwin,
	Linux:   linux,
	Windows: windows,
	Unknown: unknown,
}

var toMap = map[string]OS{
	darwin:  Darwin,
	linux:   Linux,
	windows: Windows,
}

func ToString(os OS) string {
	if val, ok := toString[os]; ok {
		return val
	}

	return toString[Unknown]
}

func ToOS(os string) OS {
	if val, ok := toMap[os]; ok {
		return val
	}

	return Unknown
}
