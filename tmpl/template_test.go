package tmpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrintableCommand(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("cat",
		getPrintableCommand("cat", nil))
	assert.Equal("cat a.txt",
		getPrintableCommand("cat", []string{"a.txt"}))
	assert.Equal(`cat "some file.txt"`,
		getPrintableCommand("cat", []string{"some file.txt"}))
	assert.Equal(`cat "some file.txt" another.txt`,
		getPrintableCommand("cat", []string{"some file.txt", "another.txt"}))
	assert.Equal("go test -v",
		getPrintableCommand("go", []string{"test", "-v"}))
	assert.Equal(`echo "He asked \"why not?\""`,
		getPrintableCommand("echo", []string{`He asked "why not?"`}))
}
