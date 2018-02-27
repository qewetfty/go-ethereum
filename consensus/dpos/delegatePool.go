package dpos

import (
	"fmt"
	"sort"
)

const (
	register = iota
	addVote
	subVote
)

type DelegatePoll struct {
	votesChan    chan CandidateWrapper
	newRoundChan chan struct{}
	votes        map[string]int64
	top          []Candidate // 当选的候选人列表
	maxElected   int         // 最多能有多少该候选人当选
}

// 候选人
type Candidate struct {
	Address string
	//Normal bool
	Vote int64 //投票数
	Nickname string // delegate name
	//publicKey string
}

type CandidateWrapper struct {
	candidate Candidate
	action    int // 0-注册,1-增加票数,2-减票
}

func NewPoll(maxElected int) *DelegatePoll {
	votes := make(map[string]int64)
	top := make([]Candidate, 0, maxElected)

	poll := &DelegatePoll{
		votesChan:    make(chan CandidateWrapper),
		newRoundChan: make(chan struct{}),
		votes:        votes,
		top:          top,
		maxElected:   maxElected,
	}

	go poll.startListening()
	return poll
}

func (p *DelegatePoll) startListening() {
	for {
		select {
		case candidateWrapper := <-p.votesChan:
			candidate := candidateWrapper.candidate
			switch candidateWrapper.action {
			case register:
				if _, ok := p.votes[candidate.Address]; !ok {
					p.votes[candidate.Address] = 0
					fmt.Printf("%s 注册代理成功\n", candidate.Address)
				} else {
					fmt.Printf("%s 代理已注册,不能重复注册\n", candidate.Address)
				}

			case addVote:
				if _, ok := p.votes[candidate.Address]; !ok {
					fmt.Printf("%s 未注册|投票失败\n", candidate.Address)
				} else {
					currentValue := p.votes[candidate.Address]
					nowVoteNumber := currentValue + candidate.Vote
					p.votes[candidate.Address] = nowVoteNumber

					p.insert(Candidate{candidate.Address, nowVoteNumber,candidate.Nickname})

					fmt.Printf("-> %s 增加 %d 票,现在票数 %d;当选列表:%v\n", candidate.Address, candidate.Vote, nowVoteNumber, p.top)
					fmt.Printf("-> 候选池列表:%v\n", p.votes)
				}
			case subVote:
				if _, ok := p.votes[candidate.Address]; !ok {
					fmt.Printf("%s 未注册|减票失败\n", candidate.Address)
				} else {
					currentValue := p.votes[candidate.Address]
					nowVoteNumber := currentValue - candidate.Vote
					if nowVoteNumber < 0 {
						nowVoteNumber = 0
					}

					p.insert(Candidate{candidate.Address, nowVoteNumber,candidate.Nickname})

					fmt.Printf("-> %s 减去 %d 票,现在票数 %d;当选列表:%v\n", candidate.Address, candidate.Vote, nowVoteNumber, p.top)
					fmt.Printf("-> 候选池列表:%v\n", p.votes)
				}

			default:
				fmt.Printf("error action %s \n", candidateWrapper.action)
			}

			//fmt.Printf("-> %s = %d ; %v\n", candidate.address, nowVoteNumber, p.top)
		case <-p.newRoundChan:
			// TODO consider clearing by range deletion to decrease GC load
			p.votes = make(map[string]int64)
			p.top = make([]Candidate, 0, p.maxElected)

		}
	}
}

// Returns minimal number of votes required to be elected int current round,i.e number
// of votes for last candidate
func (p *DelegatePoll) minVotes() int64 {
	if len(p.top) == cap(p.top) {
		return p.top[len(p.top)-1].Vote
	}
	return 0
}

func (p *DelegatePoll) insert(NewCandidate Candidate) {
	tempVotes := NewCandidate.Vote
	if len(p.top) == p.maxElected {
		minVotes := p.top[p.maxElected-1].Vote
		if tempVotes-minVotes <= 0 {
			return
		}
	}
	insertedPos := GetPosition(p.top, NewCandidate)
	if insertedPos != -1 {
		p.top[insertedPos] = NewCandidate
	} else if len(p.top) < p.maxElected {
		p.top = append(p.top, NewCandidate)
		insertedPos = len(p.top) - 1
	} else {
		insertedPos = p.maxElected - 1
		p.top[insertedPos] = NewCandidate
	}
	requiredPos := sort.Search(insertedPos, func(j int) bool {
		return p.top[j].Vote-NewCandidate.Vote < 0
	})

	if requiredPos != insertedPos {
		temp := p.top[requiredPos]
		p.top[requiredPos] = NewCandidate
		p.top[insertedPos] = temp
	}

}

func GetPosition(top []Candidate, candidate Candidate) int {
	position := -1
	for i := 0; i < len(top); i++ {
		if top[i].Address == candidate.Address {
			position = i
			break
		}
	}
	return position
}

func (p *DelegatePoll) IsElected(candidate string) (result bool) {
	if len(p.top) == 0 {
		return
	}
	votesN := p.votes[candidate]
	votes := p.minVotes()
	return votesN-votes >= 0
}

func (p *DelegatePoll) SubmitVoteFor(candidate CandidateWrapper) (err error) {
	// todo
	// if no active round err = ...;return
	p.votesChan <- candidate
	return
}

func (p *DelegatePoll) StartNewRound() {
	p.newRoundChan <- struct{}{}
}
