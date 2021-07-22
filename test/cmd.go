package test

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/PiotrKozimor/krkstops/mock"
	"github.com/matryer/is"
	"github.com/spf13/cobra"
)

func Cmd(t *testing.T, args []string, expectedLines string, cmd *cobra.Command) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go mock.Ttss(ctx)
	go mock.Airly(ctx)
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
		is.Equal(
			true,
			strings.Contains(output, line),
		)
	}
}
