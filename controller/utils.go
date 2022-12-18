package main

type void struct{}

var nullptr void

func Map[TIn, TOut any](input []TIn, mapper func(TIn) TOut) []TOut {
	var output []TOut
	for _, v := range input {
		output = append(output, mapper(v))
	}
	return output
}
