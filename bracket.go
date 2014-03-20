package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Matchup struct {
	Seeds      [2]uint8
	Names      [2]string
	Difficulty float32
	Winner     uint8
}

func (matchup *Matchup) Simulate() {
	var high, low uint8
	if matchup.Seeds[0] < matchup.Seeds[1] {
		high = 0
		low = 1
	} else {
		high = 1
		low = 0
	}
	matchup.Difficulty = 1 / float32(matchup.Seeds[low]-matchup.Seeds[high]+1)

	if rand.Float32() < matchup.Difficulty {
		matchup.Winner = low
	} else {
		matchup.Winner = high
	}
}

func (matchup *Matchup) String() string {
	return fmt.Sprintf("%s vs. %s (D: %5.3f) -> %s", matchup.Names[0], matchup.Names[1], matchup.Difficulty, matchup.Names[matchup.Winner])
}

func NewMatchup(seed1, seed2 uint8) *Matchup {
	return &Matchup{
		Seeds: [2]uint8{seed1, seed2},
		Names: [2]string{fmt.Sprintf("%02d", seed1), fmt.Sprintf("%02d", seed2)},
	}
}

type Round struct {
	Name     string
	Matchups []*Matchup
}

func (round *Round) String() string {
	var results []string
	results = append(results, fmt.Sprintf("===== %s", round.Name))
	for _, matchup := range round.Matchups {
		results = append(results, matchup.String())
	}
	return strings.Join(results, "\n")
}

func NextRound(name string, matchups []*Matchup) (round *Round) {
	round = &Round{Name: name, Matchups: make([]*Matchup, len(matchups)/2)}
	lower := 0
	upper := len(matchups) - 1
	for i := range round.Matchups {
		first := matchups[lower]
		second := matchups[upper]
		matchup := NewMatchup(first.Seeds[first.Winner], second.Seeds[second.Winner])
		matchup.Simulate()

		round.Matchups[i] = matchup
		lower++
		upper--
	}
	return
}

type Region struct {
	Name         string
	RoundTwo     *Round
	RoundThree   *Round
	SweetSixteen *Round
	EliteEight   *Round
}

func (region *Region) String() string {
	var results []string
	results = append(results, fmt.Sprintf("========== %s", region.Name))
	if region.RoundTwo != nil {
		results = append(results, region.RoundTwo.String())
		if region.RoundThree != nil {
			results = append(results, region.RoundThree.String())
			if region.SweetSixteen != nil {
				results = append(results, region.SweetSixteen.String())
				if region.EliteEight != nil {
					results = append(results, region.EliteEight.String())
				}
			}
		}
	}
	return strings.Join(results, "\n")
}

func NewRegion(name string) *Region {
	return &Region{
		Name:     name,
		RoundTwo: &Round{Name: "2nd Round", Matchups: make([]*Matchup, 8)},
	}
}

func main() {
	var (
		seed int64
		err  error
	)

	if len(os.Args) > 1 {
		var intSeed int
		intSeed, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		seed = int64(intSeed)
	} else {
		seed = time.Now().Unix()
	}
	fmt.Printf("Seed: %d\n", seed)
	rand.Seed(seed)

	regions := []*Region{
		NewRegion("South"),
		NewRegion("East"),
		NewRegion("West"),
		NewRegion("Midwest"),
	}
	for _, region := range regions {
		high := uint8(1)
		low := uint8(16)
		for i := range region.RoundTwo.Matchups {
			matchup := NewMatchup(high, low)
			matchup.Simulate()
			region.RoundTwo.Matchups[i] = matchup

			high++
			low--
		}

		region.RoundThree = NextRound("3rd Round", region.RoundTwo.Matchups)
		region.SweetSixteen = NextRound("Sweet Sixteen", region.RoundThree.Matchups)
		region.EliteEight = NextRound("Elite Eight", region.SweetSixteen.Matchups)

		fmt.Println(region)
	}
}
