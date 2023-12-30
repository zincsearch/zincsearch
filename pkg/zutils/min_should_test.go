package zutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMin(t *testing.T) {
	cases := []struct {
		subCount int
		value    interface{}
		want     int
	}{
		// Simple Integer
		{subCount: 5, value: 3, want: 3},
		{subCount: 10, value: "2", want: 2},
		{subCount: 8, value: int64(7), want: 7},
		{subCount: 9, value: 3.0, want: 3},
		{subCount: 5, value: 5.7, want: 5},

		{subCount: 5, value: -10, want: 1},
		{subCount: 3, value: 5, want: 3},

		// Negative Integer
		{subCount: 10, value: -2, want: 8},
		{subCount: 8, value: "-5", want: 3},
		{subCount: 15, value: -3.0, want: 12},
		{subCount: 9, value: -3.5, want: 5},

		// percent
		{subCount: 10, value: "80%", want: 8},
		{subCount: 10, value: "-20%", want: 8},
		{subCount: 5, value: "75%", want: 3},
		{subCount: 5, value: "-25%", want: 4},

		// combination
		{subCount: 4, value: "5<90%", want: 4},
		{subCount: 5, value: "5<90%", want: 5},
		{subCount: 7, value: "5<3", want: 3},
		{subCount: 5, value: "2<-25%", want: 4},

		// multi combinations
		{subCount: 2, value: "2<-25% 9<-3", want: 2},
		{subCount: 5, value: "4<-25% 9<-3", want: 4},
		{subCount: 10, value: "4<-40% 9<-3", want: 7},
	}
	for _, c := range cases {
		v, err := CalculateMin(c.subCount, c.value)
		assert.Nil(t, err)
		assert.Equal(t, c.want, v)
	}
}
