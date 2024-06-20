//go:build !solution

package genericsum

import (
	"cmp"
	"math/cmplx"
	"slices"

	"golang.org/x/exp/constraints"
)

func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func SortSlice[S ~[]E, E cmp.Ordered](a S) {
	slices.Sort(a)
}

func MapsEqual[M ~map[K]V, K comparable, V comparable](a, b M) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bValue, ok := b[k]; !ok || bValue != v {
			return false
		}
	}

	return true
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, elem := range s {
		if elem == v {
			return true
		}
	}
	return false
}

func MergeChans[C <-chan E, E any](chs ...C) C {
	result := make(chan E)

	go func() {
		openedChans := len(chs)

		for openedChans > 0 {
			for _, ch := range chs {
				select {
				case value, isOpened := <-ch:
					if isOpened {
						result <- value
						continue
					} else {
						openedChans--
					}
				default:
					continue
				}
			}
		}

		close(result)

	}()

	return result
}

func IsHermitianMatrix[T constraints.Integer | constraints.Complex | constraints.Float](m [][]T) bool {
	height := len(m)
	if height < 1 {
		return true
	}
	width := len(m[0])
	if width < 1 {
		return true
	}
	if height != width {
		return false
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			el := m[x][y]
			switch any(el).(type) {
			case complex64:
				oposite, ok := any(m[y][x]).(complex64)
				if !ok {
					return false
				}
				cur, ok := any(m[x][y]).(complex64)
				if !ok {
					return false
				}

				if real(oposite) != real(cur) {
					return false
				}
				if -imag(oposite) != imag(cur) {
					return false
				}

			case complex128:
				oposite, ok := any(m[y][x]).(complex128)
				if !ok {
					return false
				}
				cur, ok := any(m[x][y]).(complex128)
				if !ok {
					return false
				}

				if cmplx.Conj(cur) != oposite {
					return false
				}
			default:
				if m[x][y] != m[y][x] {
					return false
				}
			}
		}
	}

	return true
}
