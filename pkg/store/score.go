package store

import (
	"bytes"
	"encoding/gob"
	"os"
)

const (
	scoreFile = "score.gob"
)

func WriteScore(scores map[uint]uint) error {
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(scores)
	if err != nil {
		return err
	}

	err = os.WriteFile(scoreFile, buf.Bytes(), 0644)
	return err
}
