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
	output := b.String()
	expectedB := bytes.NewBufferString(expectedLines)
	sc := bufio.NewScanner(expectedB)
	for sc.Scan() {
		line := sc.Text()
		if !strings.Contains(output, line) {
			t.Fatalf("line '%s' not found in output \n%s", line, output)
		}
	}
}
