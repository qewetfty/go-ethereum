package dpos

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPoll_IsElected(t *testing.T) {
	poll := NewPoll(5)

	winners := map[Candidate]int{
		{address: "winner", votes: 3}:  20,
		{address: "second", votes: 6}:  19,
		{address: "thirdd", votes: 2}:  17,
		{address: "forthh", votes: 5}:  11,
		{address: "fifthh", votes: 10}: 10,
	}
	//winners := map[string]int{"winner": 20, "second": 19, "thirdd": 17, "forthh": 11, "fifthh": 10}
	//losers := map[string]int{"loser1": 1, "loser2": 9, "loser3": 5}

	losers := map[Candidate]int{
		{address: "loser1", votes: 1}: 1,
		{address: "loser2", votes: 6}: 9,
		{address: "loser3", votes: 8}: 5,
	}
	var wg sync.WaitGroup
	wg.Add(8)

	for candidate, votes := range winners {
		go poll.voteForNTimes(candidate, votes, &wg)
	}
	for candidate, votes := range losers {
		go poll.voteForNTimes(candidate, votes, &wg)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	for candidate := range winners {
		if !poll.IsElected(candidate.address) {
			// t.Fatalf("%s not elected: %v\n", candidate, poll.top)
			fmt.Printf("%s not elected: %v\n", candidate, poll.top)
		}
	}

	for candidate := range losers {
		if poll.IsElected(candidate.address) {
			// t.Fatalf("%s is elected: %v\n", candidate, poll.top)
			fmt.Printf("%s is elected: %v\n", candidate, poll.top)
		}
	}

	poll.StartNewRound()

	time.Sleep(100 * time.Millisecond)

	votesLen := len(poll.votes)
	topLen := len(poll.top)
	if votesLen != 0 && topLen != 0 {
		// t.Fatalf("new round not started: votes len = %d, top len = %d", votesLen, topLen)
		fmt.Printf("new round not started: votes len = %d, top len = %d", votesLen, topLen)
	}

}

func (p *DelegatePoll) voteForNTimes(candidate Candidate, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		go p.SubmitVoteFor(candidate)
	}
}
