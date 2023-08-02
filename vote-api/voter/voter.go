package voter

import (
	"encoding/json"
	"time"
)

type voterPoll struct {
	PollID   uint
	VoteDate time.Time
}

// just for use in voter-api.go
type PublicVoterPoll struct {
	PollID   uint      `json:"pollID"`
	VoteDate time.Time `json:"voteDate"`
}

type Voter struct {
	VoterID     uint
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}
type VoterList struct {
	Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoter(id uint, fn, ln string) *Voter {
	return &Voter{
		VoterID:     id,
		FirstName:   fn,
		LastName:    ln,
		VoteHistory: []voterPoll{},
	}
}

func (v *Voter) AddPoll(pollID uint) {
	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: time.Now()})
}

func (v *Voter) AddPollSpecific(pollID uint, date time.Time) {
	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: date})
}

func (v *Voter) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}
