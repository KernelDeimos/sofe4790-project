package strparse

import (
	"strconv"

	"github.com/KernelDeimos/sofe4790/estate"
)

func ParseF(state *estate.ErrorState, in string) float64 {
	out, err := strconv.ParseFloat(in, 64)
	if err != nil {
		state.Add(err)
		return 0.0
	}
	return out
}

func ParseI(state *estate.ErrorState, value string) int {
	valueInt, err := strconv.Atoi(value)
	state.Add(err)
	return valueInt
}
