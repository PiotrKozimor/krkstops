package score

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/PiotrKozimor/krkstops/pkg/cache"
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/go-redis/redis/v8"
	redi "github.com/gomodule/redigo/redis"
)

type Score struct {
	*cache.Cache
	redis  *redis.Client
	sugTmp *redisearch.Autocompleter
}

func NewScore(redisURI, sugKey string) (*Score, error) {
	r := redis.NewClient(&redis.Options{
		Addr: redisURI,
	})
	acTmp := redisearch.NewAutocompleter(redisURI, cache.GetTmpKey(sugKey))
	c, err := cache.NewCache(redisURI, sugKey)
	return &Score{
		redis:  r,
		sugTmp: acTmp,
		Cache:  c,
	}, err
}

func (c *Score) Score(ctx context.Context, cancel <-chan os.Signal, cli pb.KrkStopsClient, sleep time.Duration) error {
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
			err = addSuggestion(c.Sug, stop, score)
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

func (s *Score) RestartScoring(ctx context.Context) error {
	stops, err := s.redis.HKeys(ctx, cache.NAMES).Result()
	if err != nil {
		return err
	}
	_, err = s.redis.SAdd(ctx, cache.TO_SCORE, []interface{}{stops}...).Result()
	return err
}

func (s *Score) getStoptoScore() (*pb.Stop, error) {
	var stop pb.Stop
	id, err := redi.Int(s.Conn.Do("SPOP", cache.TO_SCORE))
	if err != nil {
		return nil, err
	}
	stop.Id = uint32(id)
	name, err := redi.String(s.Conn.Do("HGET", cache.NAMES, id))
	if err != nil {
		return nil, err
	}
	stop.Name = name
	return &stop, err

}

func (s *Score) scoreStop(ctx context.Context, stop *pb.Stop, cli pb.KrkStopsClient) (score float64, err error) {
	deps, err := cli.GetDepartures2(ctx, stop)
	if err != nil {
		return 0, err
	} else {
		return scoreByTotalDepartures(len(deps.Departures)), nil
	}
}

func (s *Score) saveScore(stop *pb.Stop, score float64) error {
	_, err := s.Conn.Do("HSET", cache.SCORES, stop.Id, score)
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
