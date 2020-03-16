package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func Test_NewArray_Success(t *testing.T) {
	var a = NewArray()
	assert.NotNil(t, a)
	assert.NotNil(t, a.inner)
	assert.NotNil(t, a.RWMutex)
}

func Test_Append_Concurrently(t *testing.T) {
	c := make(chan bool)
	a := NewArray()
	go func() {
		a.Append("a") // First conflicting access.
		c <- true
	}()
	a.Append("b") // Second conflicting access.
	<-c
	for v := range a.inner {
		fmt.Println(v)
	}
}

func BenchmarkArrayAppendSync(b *testing.B) {
	var a = NewArray()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Append(rand.Int())
	}
}

func BenchmarkArrayAppend(b *testing.B) {
	var a = make([]interface{}, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = append(a, rand.Int())
	}
}
