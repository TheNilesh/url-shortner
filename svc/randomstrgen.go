package svc

import (
	"math/rand"
	"time"
)

type RandomStrGen interface {
	Generate() string
}

type randomStrGen struct {
	minLength int
	maxLength int
	charset   string
	rand      *rand.Rand
}

func NewRandomStrGen(minLength int, maxLength int, charset string) RandomStrGen {
	return &randomStrGen{
		minLength: minLength,
		maxLength: maxLength,
		charset:   charset,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *randomStrGen) Generate() string {
	length := r.minLength
	if r.maxLength > r.minLength {
		length = r.minLength + r.rand.Intn(r.maxLength-r.minLength+1)
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = r.charset[r.rand.Intn(len(r.charset))]
	}
	return string(b)
}
