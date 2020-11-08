package stops

import (
	"context"
	"io"
	"math"

	"github.com/PiotrKozimor/krkstops/pb"
)

func ScoreStop(c pb.KrkStopsClient, shortName string) (float64, error) {
	stream, err := c.GetDepartures(context.Background(), &pb.Stop{ShortName: shortName})
	if err != nil {
		return 0, err
	}
	totalDepartures := 0
	for {
		_, err := stream.Recv()
		if err == nil {
			totalDepartures += 1
		}
		if err == io.EOF {
			return scoreByTotalDepartures(totalDepartures), err
		} else {
			return 0, err
		}
	}
}

func (c *Clients) InitializeScoring() {

}

func FinishScoring() {

}

func scoreByTotalDepartures(total int) float64 {
	return 1.0 + 0.125*math.Sqrt(float64(total))
}
