package itertools

import "iter"

func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

func Map2[T, U any](seq iter.Seq[T], f func(T) (U, error)) iter.Seq2[U, error] {
	return func(yield func(U, error) bool) {
		for v := range seq {
			u, err := f(v)
			if !yield(u, err) {
				return
			}

			if err != nil {
				return
			}
		}
	}
}

func Filter[T any](s []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}
