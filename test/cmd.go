package test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func Cmd(t *testing.T, args []string, expectedLines string, cmd *cobra.Command) {
	time.Sleep(time.Millisecond)
	is := is.New(t)
	b := bytes.Buffer{}
	cmd.SetArgs(args)
	cmd.SetOut(&b)
	err := cmd.Execute()
	is.NoErr(err)
	expectedB := bytes.NewBufferString(expectedLines)
	sc := bufio.NewScanner(expectedB)
	scActual := bufio.NewScanner(&b)
	for sc.Scan() {
		expectedLine := sc.Text()
		scActual.Scan()
		actualLine := scActual.Text()
		words := strings.Split(expectedLine, " ")
		for i := range words {
			if !strings.Contains(actualLine, words[i]) {
				t.Fatalf("word '%s' not found in output \n%s", words[i], actualLine)
			}
		}
	}
}
