package gtfs

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/PiotrKozimor/krkstops/pb"
)

type Grouper struct {
	grouping DeparturesDirections
}

func NewGrouper(r io.Reader) (*Grouper, error) {
	dec := gob.NewDecoder(r)
	g := Grouper{grouping: make(DeparturesDirections)}
	err := dec.Decode(&g.grouping)
	return &g, err
}

func NewDeparture(d *pb.Departure, s *pb.Stop) Departure {
	return Departure{
		StopId:      int(s.Id),
		PatternText: d.PatternText,
		Direction:   d.Direction,
	}
}

// Group will return 0 when no group was found for departure.
func (g *Grouper) Group(d Departure) []int {
	return g.grouping[d]
}

func (g *Grouper) Dump(w io.Writer) error {
	for dep := range g.grouping {
		_, err := fmt.Fprintf(w, "%s %s %d %d\n", dep.PatternText, dep.Direction, dep.StopId, g.grouping[dep])
		if err != nil {
			return err
		}
	}
	return nil
}
