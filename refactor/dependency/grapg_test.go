package refactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isLoop_sameName(t *testing.T) {
	// arrange
	g := NewGraph()

	// action
	b := g.isLoop("a", "a")

	// assert
	assert.True(t, b)
}

func Test_isLoop_parentNotExist(t *testing.T) {
	// arrange
	g := NewGraph()
	g.AddDependency("a", "b")

	// action
	b := g.isLoop("c", "b")

	// assert
	assert.False(t, b)
}

func Test_isLoop_childNotExist(t *testing.T) {
	// arrange
	g := NewGraph()
	g.AddDependency("a", "b")

	// action
	b := g.isLoop("a", "c")

	// assert
	assert.False(t, b)
}

func Test_isLoop_loop(t *testing.T) {
	// arrange
	// a <- b <- c
	g := NewGraph()
	g.AddDependency("a", "b")
	g.AddDependency("b", "c")

	// action
	b := g.isLoop("c", "a")

	// assert
	assert.True(t, b)
}
