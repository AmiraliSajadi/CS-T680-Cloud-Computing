package api

type pollOption struct {
	PollOptionID    uint
	PollOptionValue string
}

type Poll struct {
	PollID       uint
	PollTitle    string
	PollQuestion string
	PollOptions  []pollOption
}

// Hard-coded new poll:
func NewSamplePoll() *Poll {
	return &Poll{
		PollID:       1,
		PollTitle:    "Favorite Pet",
		PollQuestion: "What type of pet do you like best?",
		PollOptions: []pollOption{
			{PollOptionID: 1, PollOptionValue: "Dog"},
			{PollOptionID: 2, PollOptionValue: "Cat"},
			{PollOptionID: 3, PollOptionValue: "Fish"},
			{PollOptionID: 4, PollOptionValue: "Bird"},
			{PollOptionID: 5, PollOptionValue: "NONE"},
		},
	}
}
