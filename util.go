package jpkg

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"slices"
)

func toSeq1[T any](items []T) iter.Seq[T] {
	return slices.Values(items)
}

func promoteSeq1[T1, T2, T3 any](s iter.Seq[T1], f func(T1) (T2, T3)) iter.Seq2[T2, T3] {
	return func(yield func(T2, T3) bool) {
		for v := range s {
			if !yield(f(v)) {
				break
			}
		}
	}
}

type item[K, V any] struct {
	k K
	v V
}

func split[K, V any](items []item[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, v := range items {
			if !yield(v.k, v.v) {
				return
			}
		}
	}
}

func collect2[T1 comparable, T2 any](iter iter.Seq2[T1, T2]) []item[T1, T2] {
	items := []item[T1, T2]{}

	for k, v := range iter {
		items = append(items, item[T1, T2]{k, v})
	}

	return items
}

func mapToSortedStream[K cmp.Ordered, V any](m map[K]V) []item[K, V] {

	items := make([]item[K, V], len(m))

	i := 0
	for k, v := range m {
		items[i] = item[K, V]{k, v}
		i++
	}

	slices.SortStableFunc(
		items,
		func(a, b item[K, V]) int {
			return cmp.Compare(a.k, b.k)
		},
	)

	return collect2(
		promoteSeq1(
			toSeq1(items),
			func(item item[K, V]) (K, V) {
				return item.k, item.v
			},
		),
	)
}

type sizeableReader[T any] interface {
	Size() T
}

func isSizeableReader(source io.Reader) (uint64, bool) {
	if sizable, isSizable := source.(sizeableReader[int]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[uint]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[int64]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[uint64]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[int32]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[uint32]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[int16]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[uint16]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[int8]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[uint8]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[byte]); isSizable {
		return uint64(sizable.Size()), true
	}

	if sizable, isSizable := source.(sizeableReader[rune]); isSizable {
		return uint64(sizable.Size()), true
	}

	return 0, false
}

func serializeMetadataToJSON(data any) (string, uint64, error) {
	if data == nil {
		data = struct{}{}
	}

	metadataJsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", 0, fmt.Errorf("error parsing metadata as json: %w", err)
	}

	str := string(metadataJsonBytes)
	len := uint64(len(str))
	return str, len, nil
}
