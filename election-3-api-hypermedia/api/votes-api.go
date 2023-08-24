package api

type Vote struct {
	VoteID    uint
	VoterID   uint
	PollID    uint
	VoteValue uint
}

// Hard-coded new vote:
func NewSampleVote() *Vote {
	return &Vote{
		VoteID:    1,
		PollID:    1,
		VoterID:   1,
		VoteValue: 1,
	}
}
