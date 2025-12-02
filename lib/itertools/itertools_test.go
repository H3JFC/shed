package itertools

import (
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

	input := slices.Values([]int{1, 2, 3, 4, 5})
	doubled := Map(input, func(n int) int { return n * 2 })

	// Collect only first 3 elements
	var got []int

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

	input := slices.Values([]int{1, 2, 3})

	ptrs := Map(input, func(n int) *int {
		val := n * 2

		return &val
	})

	var got []int
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

	input := slices.Values([]int{1, 2, 3, 4, 5})

	mapped := Map2(input, func(n int) (int, error) {
		return n * 2, nil
	})

	var got []int

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

	input := slices.Values([]int{1, 2, 3, 4, 5})

	mapped := Map2(input, func(n int) (int, error) {
		if n == 3 {
			return 0, fmt.Errorf("error at %d", n)
		}

		return n * 2, nil
	})

	var got []int
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

	input := slices.Values([]int{1, 2, 3})

	mapped := Map2(input, func(n int) (int, error) {
		return 0, fmt.Errorf("always fails")
	})

	var got []int
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

	input := slices.Values([]int{})

	mapped := Map2(input, func(n int) (int, error) {
		return n * 2, nil
	})

	var got []int
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

	input := slices.Values([]string{"123", "456", "abc", "789"})

	mapped := Map2(input, func(s string) (int, error) {
		return strconv.Atoi(s)
	})

	var got []int
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

	input := slices.Values([]int{1, 2, 3, 4, 5})

	mapped := Map2(input, func(n int) (int, error) {
		return n * 2, nil
	})

	var got []int

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

	input := slices.Values([]int{1, 2, 3})

	mapped := Map2(input, func(n int) (string, error) {
		return fmt.Sprintf("num-%d", n), nil
	})

	var got []string
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

	input := slices.Values([]Person{
		{Name: "Alice", Age: 30},
		{Name: "", Age: 25}, // Invalid - empty name
		{Name: "Bob", Age: 35},
	})

	mapped := Map2(input, func(p Person) (PersonDTO, error) {
		if p.Name == "" {
			return PersonDTO{}, fmt.Errorf("name cannot be empty")
		}

		return PersonDTO{FullName: p.Name, Years: p.Age}, nil
	})

	var got []PersonDTO
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
			return 0, fmt.Errorf("error at %d", n)
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
