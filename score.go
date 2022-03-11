package krkstops

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/RediSearch/redisearch-go/redisearch"
	redi "github.com/gomodule/redigo/redis"
)

func (c *Cache) Score(ctx context.Context, cancel <-chan os.Signal, cli pb.KrkStopsClient, sleep time.Duration) error {
outer:
	for {
		select {
		case <-cancel:
			break outer
		default:
			stop, err := c.getStoptoScore()
			if err != nil {
				if err == redi.ErrNil {
					fmt.Printf("err: %v\n", err)
					return nil
				}
				return err
			}
			score, err := c.scoreStop(ctx, stop, cli)
			if err != nil {
				return err
			}
			err = c.saveScore(stop, score)
			if err != nil {
				return err
			}
			err = addSuggestion(c.sug, stop, score)
			if err != nil {
				return err
			}
			fmt.Printf("assigned score %f to stop %d  %s\n", score, stop.Id, stop.ShortName)
			time.Sleep(sleep)
		}

	}
	_, err := c.redis.BgSave(ctx).Result()
	return err
}

func (c *Cache) RestartScoring(ctx context.Context) error {
	stops, err := c.redis.HKeys(ctx, NAMES).Result()
	if err != nil {
		return err
	}
	_, err = c.redis.SAdd(ctx, TO_SCORE, []interface{}{stops}...).Result()
	return err
}

func (c *Cache) getStoptoScore() (*pb.Stop, error) {
	var stop pb.Stop
	id, err := redi.Int(c.conn.Do("SPOP", TO_SCORE))
	if err != nil {
		return nil, err
	}
	stop.Id = uint32(id)
	name, err := redi.String(c.conn.Do("HGET", NAMES, id))
	if err != nil {
		return nil, err
	}
	stop.Name = name
	return &stop, err

}

func (c *Cache) scoreStop(ctx context.Context, stop *pb.Stop, cli pb.KrkStopsClient) (score float64, err error) {
	deps, err := cli.GetDepartures2(ctx, stop)
	if err != nil {
		return 0, err
	} else {
		return scoreByTotalDepartures(len(deps.Departures)), nil
	}
}

func (c *Cache) saveScore(stop *pb.Stop, score float64) error {
	_, err := c.conn.Do("HSET", SCORES, stop.Id, score)
	return err
}

func addSuggestion(sug *redisearch.Autocompleter, stop *pb.Stop, score float64) error {
	splitted := strings.Split(stop.Name, " ")
	for index, value := range splitted {
		if value == "(nż)" {
			splitted = append(splitted[:index], splitted[index+1:]...)
		}
	}
	for i := 0; i < len(splitted); i++ {
		swapped := append(splitted[i:], splitted[:i]...)
		term := strings.Join(swapped, " ")
		err := sug.AddTerms(redisearch.Suggestion{Term: term, Score: score, Payload: strconv.Itoa(int(stop.Id))})
		if err != nil {
			return err
		}
	}
	return nil
}

func scoreByTotalDepartures(total int) float64 {
	return 1.0 + 0.5*math.Sqrt(float64(total))
}
