package main

import (
	lru "github.com/hashicorp/golang-lru/v2"
)

var fontBytesCache, _ = lru.New[string, []byte](1)

func Map[T, V any](ts []T, fn func(T) V) []V {
	// I've got no idea how this works, but it does
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
