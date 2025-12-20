package itertools

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"testing"
)

func TestMap_DoubleIntegers(t *testing.T) {
	t.Parallel()

	input := slices.Values([]int{1, 2, 3, 4, 5})
	doubled := Map(input, func(n int) int { return n * 2 })

	want := []int{2, 4, 6, 8, 10}
	got := slices.Collect(doubled)

	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap_IntToString(t *testing.T) {
	t.Parallel()

	input := slices.Values([]int{1, 2, 3})
	stringified := Map(input, func(n int) string { return fmt.Sprintf("num-%d", n) })

	want := []string{"num-1", "num-2", "num-3"}
	got := slices.Collect(stringified)

	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap_StringToLength(t *testing.T) {
	t.Parallel()

	input := slices.Values([]string{"hello", "world", "go"})
	lengths := Map(input, func(s string) int { return len(s) })

	want := []int{5, 5, 2}
	got := slices.Collect(lengths)

	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap_EmptySequence(t *testing.T) {
	t.Parallel()

	input := slices.Values([]int{})
	doubled := Map(input, func(n int) int { return n * 2 })

	got := slices.Collect(doubled)

	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestMap_StructTransformation(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	type PersonDTO struct {
		FullName string
		Years    int
	}

	input := slices.Values([]Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	})

	dtos := Map(input, func(p Person) PersonDTO {
		return PersonDTO{FullName: p.Name, Years: p.Age}
	})

	want := []PersonDTO{
		{FullName: "Alice", Years: 30},
		{FullName: "Bob", Years: 25},
	}
	got := slices.Collect(dtos)

	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap_EarlyTermination(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5}
	doubled := Map(slices.Values(input), func(n int) int { return n * 2 })

	got := make([]int, 0, len(input))

	for v := range doubled {
		got = append(got, v)
		if len(got) == 3 {
			break
		}
	}

	want := []int{2, 4, 6}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap_ChainedMaps(t *testing.T) {
	t.Parallel()

	input := slices.Values([]int{1, 2, 3})

	// Chain multiple maps
	result := Map(
		Map(
			Map(input, func(n int) int { return n * 2 }),
			func(n int) int { return n + 1 },
		),
		func(n int) string { return fmt.Sprintf("result: %d", n) },
	)

	want := []string{"result: 3", "result: 5", "result: 7"}
	got := slices.Collect(result)

	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap_PointerTransformation(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3}

	ptrs := Map(slices.Values(input), func(n int) *int {
		val := n * 2

		return &val
	})

	got := make([]int, 0, len(input))
	for ptr := range ptrs {
		got = append(got, *ptr)
	}

	want := []int{2, 4, 6}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap_LazyEvaluation(t *testing.T) {
	t.Parallel()

	callCount := 0
	input := slices.Values([]int{1, 2, 3, 4, 5})

	mapped := Map(input, func(n int) int {
		callCount++

		return n * 2
	})

	// Map should be lazy - function not called yet
	if callCount != 0 {
		t.Errorf("expected map to be lazy, but function was called %d times", callCount)
	}

	// Consume only 2 elements
	count := 0
	for range mapped {
		count++
		if count == 2 {
			break
		}
	}

	// Should have only called the function twice
	if callCount != 2 {
		t.Errorf("expected function to be called 2 times, was called %d times", callCount)
	}
}

func TestMap_BoolTransformation(t *testing.T) {
	t.Parallel()

	input := slices.Values([]int{1, 2, 3, 4, 5})
	isEven := Map(input, func(n int) bool { return n%2 == 0 })

	want := []bool{false, true, false, true, false}
	got := slices.Collect(isEven)

	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap2_Success(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5}

	mapped := Map2(slices.Values(input), func(n int) (int, error) {
		return n * 2, nil
	})

	got := make([]int, 0, len(input))

	for val, err := range mapped {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got = append(got, val)
	}

	want := []int{2, 4, 6, 8, 10}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap2_ErrorStopsIteration(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5}

	mapped := Map2(slices.Values(input), func(n int) (int, error) {
		if n == 3 {
			return 0, fmt.Errorf("error at %d", n) // nolint:err113
		}

		return n * 2, nil
	})

	got := make([]int, 0, len(input))

	var gotErr error

	for val, err := range mapped {
		if err != nil {
			gotErr = err

			break
		}

		got = append(got, val)
	}

	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}

	if gotErr.Error() != "error at 3" {
		t.Errorf("expected 'error at 3', got '%v'", gotErr)
	}

	// Should have values before the error
	want := []int{2, 4}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap2_ErrorOnFirstElement(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3}

	mapped := Map2(slices.Values(input), func(_ int) (int, error) {
		return 0, errors.New("always fails") //nolint:err113
	})

	got := make([]int, 0, len(input))

	var gotErr error

	for val, err := range mapped {
		if err != nil {
			gotErr = err

			break
		}

		got = append(got, val)
	}

	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}

	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestMap2_EmptySequence(t *testing.T) {
	t.Parallel()

	input := []int{}

	mapped := Map2(slices.Values(input), func(n int) (int, error) {
		return n * 2, nil
	})

	got := make([]int, 0)

	for val, err := range mapped {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got = append(got, val)
	}

	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestMap2_StringConversionWithError(t *testing.T) {
	t.Parallel()

	input := []string{"123", "456", "abc", "789"}

	mapped := Map2(slices.Values(input), strconv.Atoi)

	got := make([]int, 0, len(input))

	var gotErr error

	for val, err := range mapped {
		if err != nil {
			gotErr = err

			break
		}

		got = append(got, val)
	}

	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}

	want := []int{123, 456}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap2_EarlyTermination(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5}

	mapped := Map2(slices.Values(input), func(n int) (int, error) {
		return n * 2, nil
	})

	got := make([]int, 0, len(input))

	for val, err := range mapped {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got = append(got, val)
		if len(got) == 3 {
			break
		}
	}

	want := []int{2, 4, 6}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap2_IntToString(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3}

	mapped := Map2(slices.Values(input), func(n int) (string, error) {
		return fmt.Sprintf("num-%d", n), nil
	})

	got := make([]string, 0, len(input))

	for val, err := range mapped {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got = append(got, val)
	}

	want := []string{"num-1", "num-2", "num-3"}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap2_LazyEvaluation(t *testing.T) {
	t.Parallel()

	callCount := 0
	input := slices.Values([]int{1, 2, 3, 4, 5})

	mapped := Map2(input, func(n int) (int, error) {
		callCount++

		return n * 2, nil
	})

	// Map2 should be lazy - function not called yet
	if callCount != 0 {
		t.Errorf("expected map to be lazy, but function was called %d times", callCount)
	}

	// Consume only 2 elements
	count := 0

	for _, err := range mapped {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		count++
		if count == 2 {
			break
		}
	}

	// Should have only called the function twice
	if callCount != 2 {
		t.Errorf("expected function to be called 2 times, was called %d times", callCount)
	}
}

func TestMap2_StructTransformationWithValidation(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	type PersonDTO struct {
		FullName string
		Years    int
	}

	input := []Person{
		{Name: "Alice", Age: 30},
		{Name: "", Age: 25}, // Invalid - empty name
		{Name: "Bob", Age: 35},
	}

	mapped := Map2(slices.Values(input), func(p Person) (PersonDTO, error) {
		if p.Name == "" {
			return PersonDTO{}, errors.New("name cannot be empty") //nolint:err113
		}

		return PersonDTO{FullName: p.Name, Years: p.Age}, nil
	})

	got := make([]PersonDTO, 0, len(input))

	var gotErr error

	for val, err := range mapped {
		if err != nil {
			gotErr = err

			break
		}

		got = append(got, val)
	}

	if gotErr == nil {
		t.Fatal("expected error for empty name")
	}

	// Should have one valid person before the error
	want := []PersonDTO{{FullName: "Alice", Years: 30}}
	if !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestMap2_MultipleErrorsOnlyReturnsFirst(t *testing.T) {
	t.Parallel()

	input := slices.Values([]int{1, 2, 3, 4, 5})

	mapped := Map2(input, func(n int) (int, error) {
		if n%2 == 0 {
			return 0, fmt.Errorf("error at %d", n) // nolint:err113
		}

		return n * 2, nil
	})

	var gotErr error

	for _, err := range mapped {
		if err != nil {
			gotErr = err

			break
		}
	}

	if gotErr == nil {
		t.Fatal("expected error, got nil")
	}

	// Should get error from first even number (2)
	if gotErr.Error() != "error at 2" {
		t.Errorf("expected 'error at 2', got '%v'", gotErr)
	}
}

func TestFilter_EvenNumbers(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5, 6}
	filtered := Filter(input, func(n int) bool { return n%2 == 0 })

	want := []int{2, 4, 6}

	if !slices.Equal(filtered, want) {
		t.Errorf("expected %v, got %v", want, filtered)
	}
}

func TestFilter_OddNumbers(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5, 6}
	filtered := Filter(input, func(n int) bool { return n%2 != 0 })

	want := []int{1, 3, 5}

	if !slices.Equal(filtered, want) {
		t.Errorf("expected %v, got %v", want, filtered)
	}
}

func TestFilter_StringLength(t *testing.T) {
	t.Parallel()

	input := []string{"go", "rust", "c", "python", "java"}
	filtered := Filter(input, func(s string) bool { return len(s) > 3 })

	want := []string{"rust", "python", "java"}

	if !slices.Equal(filtered, want) {
		t.Errorf("expected %v, got %v", want, filtered)
	}
}

func TestFilter_EmptySlice(t *testing.T) {
	t.Parallel()

	input := []int{}
	filtered := Filter(input, func(n int) bool { return n > 0 })

	if len(filtered) != 0 {
		t.Errorf("expected empty slice, got %v", filtered)
	}
}

func TestFilter_NoMatches(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5}
	filtered := Filter(input, func(n int) bool { return n > 10 })

	if len(filtered) != 0 {
		t.Errorf("expected empty slice, got %v", filtered)
	}
}

func TestFilter_StructFields(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	input := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 17},
		{Name: "Charlie", Age: 25},
		{Name: "Dave", Age: 16},
	}

	filtered := Filter(input, func(p Person) bool { return p.Age >= 18 })

	want := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Charlie", Age: 25},
	}

	if !slices.Equal(filtered, want) {
		t.Errorf("expected %v, got %v", want, filtered)
	}
}

func TestFilter_PointerValues(t *testing.T) {
	t.Parallel()

	val1, val2, val3 := 1, 2, 3
	input := []*int{&val1, nil, &val2, nil, &val3}

	filtered := Filter(input, func(p *int) bool { return p != nil })

	if len(filtered) != 3 {
		t.Errorf("expected 3 non-nil pointers, got %d", len(filtered))
	}

	for i, ptr := range filtered {
		if ptr == nil {
			t.Errorf("element %d should not be nil", i)
		}
	}
}

func TestFilter_NegativeNumbers(t *testing.T) {
	t.Parallel()

	input := []int{-3, -1, 0, 1, 2, 3, -5}
	filtered := Filter(input, func(n int) bool { return n < 0 })

	want := []int{-3, -1, -5}

	if !slices.Equal(filtered, want) {
		t.Errorf("expected %v, got %v", want, filtered)
	}
}

func TestFilter_PreservesOrder(t *testing.T) {
	t.Parallel()

	input := []int{5, 1, 8, 2, 9, 3, 6}
	filtered := Filter(input, func(n int) bool { return n > 4 })

	want := []int{5, 8, 9, 6}

	if !slices.Equal(filtered, want) {
		t.Errorf("expected %v, got %v", want, filtered)
	}
}

func TestFilter_ComplexPredicate(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	// Filter for numbers divisible by 2 or 3
	filtered := Filter(input, func(n int) bool {
		return n%2 == 0 || n%3 == 0
	})

	want := []int{2, 3, 4, 6, 8, 9, 10}

	if !slices.Equal(filtered, want) {
		t.Errorf("expected %v, got %v", want, filtered)
	}
}
