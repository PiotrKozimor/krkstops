package stops

import (
	"context"
	"errors"
	"io"
	"math"

	"github.com/PiotrKozimor/krkstops/pb"
)

const (
	STOPS_TO_SCORE = "stops.toScore"
	STOPS          = "stops"
	STOPS_NEW      = "stops.new"
	STOPS_TMP      = "stops.tmp"
)

var ScoringInitialized = errors.New("Scoring is already initialized")

func (c *Clients) ScoreStop(krk pb.KrkStopsClient, shortName string) (float64, error) {
	stream, err := krk.GetDepartures(context.Background(), &pb.Stop{ShortName: shortName})
	if err != nil {
		return 0, err
	}
	totalDepartures := 0
	for {
		_, err := stream.Recv()
		if err == nil {
			totalDepartures += 1
		} else if err == io.EOF {
			return scoreByTotalDepartures(totalDepartures), nil
		} else {
			return 0, err
		}
	}
}

func (c *Clients) GetStopToScore() (shortName string, err error) {
	shortName, err = c.Redis.SPop(STOPS_TO_SCORE).Result()
	return
}

func (c *Clients) InitializeScoring() error {
	exist, err := c.Redis.Exists(STOPS_TO_SCORE).Result()
	if err != nil {
		return err
	}
	if exist != 0 {
		return ScoringInitialized
	}
	return c.RestartScoring()
}

func (c *Clients) RestartScoring() error {
	res, err := c.Redis.SUnionStore(STOPS_TO_SCORE, STOPS).Result()
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New("stops set not created, please call 'stopctl update' command")
	}
	return nil
}

func (c *Clients) FinishScoring() error {
	res, err := c.Redis.Del(STOPS_TO_SCORE).Result()
	if err != nil {
		return err
	}
	if res != 1 {
		return errors.New("scoring already finished")
	}
	return nil
}

func scoreByTotalDepartures(total int) float64 {
	return 1.0 + 0.5*math.Sqrt(float64(total))
}

func (c *Clients) ApplyScore(score float64, shortName string) error {
	name, err := c.Redis.Get(shortName).Result()
	if err != nil {
		return err
	}
	return c.AddSuggestion(shortName, name, score)
}
