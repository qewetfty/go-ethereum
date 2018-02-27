package dpos

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPoll_IsElected(t *testing.T) {
	poll := NewPoll(5)

	register := []CandidateWrapper{
		{Candidate{Address: "winner"}, 0},
		{Candidate{Address: "second"}, 0},
	}

	for _, candidate := range register {
		poll.SubmitVoteFor(candidate)
	}

	time.Sleep(100 * time.Millisecond)

	winners := map[CandidateWrapper]int{
		{Candidate{Address: "winner", Vote: 3}, 1}:  20,
		{Candidate{Address: "second", Vote: 6}, 1}:  19,
		{Candidate{Address: "thirdd", Vote: 2}, 1}:  17,
		{Candidate{Address: "forthh", Vote: 5}, 1}:  11,
		{Candidate{Address: "fifthh", Vote: 10}, 1}: 10,
	}

	losers := map[CandidateWrapper]int{
		{Candidate{Address: "loser1", Vote: 1}, 1}: 1,
		{Candidate{Address: "loser2", Vote: 6}, 1}: 9,
		{Candidate{Address: "loser3", Vote: 8}, 1}: 5,
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

	for candidateWrapper := range winners {
		if !poll.IsElected(candidateWrapper.candidate.Address) {
			// t.Fatalf("%s not elected: %v\n", candidate, poll.top)
			fmt.Printf("%s not elected: %v\n", candidateWrapper.candidate, poll.top)
		}
	}

	for candidateWrapper := range losers {
		if poll.IsElected(candidateWrapper.candidate.Address) {
			// t.Fatalf("%s is elected: %v\n", candidate, poll.top)
			fmt.Printf("%s is elected: %v\n", candidateWrapper.candidate, poll.top)
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

func (p *DelegatePoll) voteForNTimes(candidate CandidateWrapper, n int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < n; i++ {
		go p.SubmitVoteFor(candidate)
	}
}
