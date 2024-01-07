package datastructure

import "reflect"

func Map[T, V any](s []T, f func(T) V) (res []V) {
	for _, item := range s {
		res = append(res, f(item))
	}
	return res
}

func Filter[T any](s []T, f func(T) bool) (res []T) {
	for _, s := range s {
		if f(s) {
			res = append(res, s)
		}
	}
	return
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}
