package krkstops

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PiotrKozimor/krkstops/pb"
	"github.com/RediSearch/redisearch-go/redisearch"
	redi "github.com/gomodule/redigo/redis"
	"golang.org/x/net/context"
)

func (db *Cache) Score(ctx context.Context, c <-chan os.Signal, cli pb.KrkStopsClient, sleep time.Duration) error {
outer:
	for {
		select {
		case <-c:
			break outer
		default:
			stop, err := db.getStoptoScore()
			if err != nil {
				if err == redi.ErrNil {
					return nil
				}
				return err
			}
			score, err := db.scoreStop(ctx, stop, cli)
			if err != nil {
				return err
			}
			err = db.saveScore(stop, score)
			if err != nil {
				return err
			}
			err = db.addSuggestion(stop, score)
			if err != nil {
				return err
			}
			fmt.Printf("assigned score %f to stop %d  %s\n", score, stop.Id, stop.ShortName)
			time.Sleep(sleep)
		}

	}
	return nil
}

func (db *Cache) getStoptoScore() (*pb.Stop, error) {
	var stop pb.Stop
	id, err := redi.Int(db.conn.Do("SPOP", TO_SCORE))
	if err != nil {
		return nil, err
	}
	stop.Id = uint32(id)
	name, err := redi.String(db.conn.Do("HGET", NAMES, id))
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

func (db *Cache) saveScore(stop *pb.Stop, score float64) error {
	_, err := db.conn.Do("HSET", SCORES, stop.Id, score)
	return err
}

func (c *Cache) addSuggestion(stop *pb.Stop, score float64) error {
	splitted := strings.Split(stop.Name, " ")
	for index, value := range splitted {
		if value == "(nż)" {
			splitted = append(splitted[:index], splitted[index+1:]...)
		}
	}
	for i := 0; i < len(splitted); i++ {
		swapped := append(splitted[i:], splitted[:i]...)
		term := strings.Join(swapped, " ")
		err := c.sugTmp.AddTerms(redisearch.Suggestion{Term: term, Score: score, Payload: strconv.Itoa(int(stop.Id))})
		c.sugTmp.DeleteTerms()
		if err != nil {
			return err
		}
	}
	return nil
}

func scoreByTotalDepartures(total int) float64 {
	return 1.0 + 0.5*math.Sqrt(float64(total))
}
