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

	// assert
	assert.True(t, g.isLoop("c", "a"))
	assert.True(t, g.isLoop("c", "b"))
	assert.True(t, g.isLoop("b", "a"))
}

func Test_isLoop_noloop(t *testing.T) {
	// arrange
	// a <- b <- c
	// d <- e <- f
	g := NewGraph()
	g.AddDependency("a", "b")
	g.AddDependency("b", "c")
	g.AddDependency("d", "e")
	g.AddDependency("e", "f")

	// action

	// assert
	for _, s1 := range []string{"a", "b", "c"} {
		for _, s2 := range []string{"d", "e", "f"} {
			assert.False(t, g.isLoop(s1, s2))
			assert.False(t, g.isLoop(s2, s1))
		}
	}
}

func Test_AddDependency_loop(t *testing.T) {
	// arrange
	g := NewGraph()
	g.AddDependency("a", "b")

	// action
	b := g.AddDependency("b", "a")

	// assert
	assert.False(t, b)
}

func Test_AddDependency_parentAndChildNotExist(t *testing.T) {
	// arrange
	g := NewGraph()

	// action
	b := g.AddDependency("a", "b")

	// assert
	assert.True(t, b)
}

func Test_AddDependency_parentAndChildExist(t *testing.T) {
	// arrange
	g := NewGraph()
	g.AddDependency("a", "b")

	// action
	b := g.AddDependency("a", "b")

	// assert
	assert.True(t, b)
}
