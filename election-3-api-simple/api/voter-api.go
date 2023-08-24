package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"voter-api/voter"

	"github.com/gin-gonic/gin"
)

type VoterApi struct {
	voterList voter.VoterList
}

func NewVoterApi() *VoterApi {
	return &VoterApi{
		voterList: voter.VoterList{
			Voters: make(map[uint]voter.Voter),
		},
	}
}

func (v *VoterApi) AddVoter(voterID uint, firstName, lastName string) {
	v.voterList.Voters[voterID] = *voter.NewVoter(voterID, firstName, lastName)
}

func (v *VoterApi) AddVoterGoodHandler(c *gin.Context) {
	var this_voter voter.Voter
	if err := c.ShouldBindJSON(&this_voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, found := v.voterList.Voters[uint(this_voter.VoterID)]
	if found {
		log.Println("Voter ID already exists: ", this_voter.VoterID)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else {
		v.voterList.Voters[this_voter.VoterID] = *voter.NewVoter(
			this_voter.VoterID,
			this_voter.FirstName,
			this_voter.LastName,
		)
	}
}

func (v *VoterApi) NewSampleVoter() {
	v.voterList.Voters[7] = *voter.NewVoter(
		7,
		"Kshitij",
		"Kayastha",
	)
}

func (v *VoterApi) AddPoll(voterID, pollID uint) {
	voter := v.voterList.Voters[voterID]
	voter.AddPoll(pollID)
	v.voterList.Voters[voterID] = voter
}

// This one does not use the payload (only uses the id and pollid in url)
func (v *VoterApi) AddPollHandlerBadImplementation(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pidS := c.Param("pollid")
	pid64, err := strconv.ParseInt(pidS, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	this_voter, found := v.voterList.Voters[uint(id64)]
	if found {
		// check if poll id already exists
		for _, value := range this_voter.VoteHistory {
			if int64(value.PollID) == pid64 {
				log.Println("The pollid for this user already exists", err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
		// add the poll
		pidUint, err := convertToUnit(pid64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		} else {
			this_voter.AddPoll(pidUint)
			v.voterList.Voters[uint(id64)] = this_voter
			log.Println("poll added to the user")
			return
		}
	} else {
		log.Println("Voter not found: ", id64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

// the version that uses the payload data
func (v *VoterApi) AddPollHandlerGoodImplementation(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// check if voter exists
	this_voter, found := v.voterList.Voters[uint(id64)]
	if found {
		var poll voter.PublicVoterPoll
		if err := c.ShouldBindJSON(&poll); err != nil {
			log.Println("Error binding JSON: ", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// check if poll id already exists
		for _, value := range this_voter.VoteHistory {
			if value.PollID == poll.PollID {
				log.Println("The pollid for this user already exists", err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
		// add the poll
		// this_voter.AddPoll(poll.PollID)
		this_voter.AddPollSpecific(poll.PollID, poll.VoteDate)
		v.voterList.Voters[uint(id64)] = this_voter
		log.Println("poll added to the user")
		return
	} else {
		log.Println("Voter not found: ", id64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

}

func (v *VoterApi) GetVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Check if the voter exists
	this_voter, found := v.voterList.Voters[uint(id64)]
	if found {
		c.JSON(http.StatusOK, this_voter)
	} else {
		log.Println("Voter not found: ", id64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func (v *VoterApi) UpdateVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Check if the voter exists
	this_voter, found := v.voterList.Voters[uint(id64)]
	if found {
		// do stuff
		var updatedVoter voter.Voter

		if err := c.ShouldBindJSON(&updatedVoter); err != nil {
			log.Println("Error binding JSON: ", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		} else {
			// URL and payload data don't match
			if id64 != uint64(updatedVoter.VoterID) {
				log.Println("The ID in URL does not match the one in payload", updatedVoter.VoterID)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			} else {
				this_voter.FirstName = updatedVoter.FirstName
				this_voter.LastName = updatedVoter.LastName
				v.voterList.Voters[uint(id64)] = this_voter
				log.Println("User Updated")
				c.AbortWithStatus(http.StatusOK)
			}
		}
	} else {
		log.Println("Voter not found: ", id64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func (v *VoterApi) DeleteVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, _ := strconv.ParseInt(idS, 10, 32)

	_, found := v.voterList.Voters[uint(id64)]
	if found {
		// detele
		delete(v.voterList.Voters, uint(id64))
		c.Status(http.StatusOK)
	} else {
		log.Println("Voter not found: ", id64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func (v *VoterApi) UpdatePoll(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	pidS := c.Param("pollid")
	pid64, err := strconv.ParseInt(pidS, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var poll voter.PublicVoterPoll
	if err := c.ShouldBindJSON(&poll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	this_voter, found := v.voterList.Voters[uint(id64)]
	if pid64 == int64(poll.PollID) && found {
		// do stuff
		for i, value := range this_voter.VoteHistory {
			if int64(value.PollID) == pid64 {
				this_voter.VoteHistory[i].PollID = poll.PollID
				this_voter.VoteHistory[i].VoteDate = poll.VoteDate
				v.voterList.Voters[uint(id64)] = this_voter
				log.Println("Poll updated")
				c.AbortWithStatus(http.StatusOK)
				return
			} else {
				log.Println("The PooID does not exist", err)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

}

func (v *VoterApi) DeletePoll(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	pidS := c.Param("pollid")
	pid64, err := strconv.ParseInt(pidS, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	this_voter := v.voterList.Voters[uint(id64)]
	delete_index := -1
	for i, value := range this_voter.VoteHistory {
		if int64(value.PollID) == pid64 {
			delete_index = i
			break
		}
	}
	this_voter.VoteHistory = append(this_voter.VoteHistory[:delete_index], this_voter.VoteHistory[delete_index+1:]...)
	v.voterList.Voters[uint(id64)] = this_voter
	log.Println("Poll deleted")
	c.AbortWithStatus(http.StatusOK)
}

func (v *VoterApi) GetVoterJson(voterID uint) string {
	voter := v.voterList.Voters[voterID]
	return voter.ToJson()
}

func (v *VoterApi) GetVoterList(c *gin.Context) {
	c.JSON(http.StatusOK, v.voterList)
}

func (v *VoterApi) GetVoterListJson() string {
	b, _ := json.Marshal(v.voterList)
	return string(b)
}

func (v *VoterApi) GetPolls(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	this_voter, found := v.voterList.Voters[uint(id64)]
	if found {
		c.JSON(http.StatusOK, this_voter.VoteHistory)
	} else {
		log.Println("Voter not found: ", id64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func (v *VoterApi) GetPoll(c *gin.Context) {
	vidS := c.Param("id")
	vid64, err := strconv.ParseInt(vidS, 10, 32)
	if err != nil {
		log.Println("Error converting voter id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	pidS := c.Param("pollid")
	pid64, err := strconv.ParseInt(pidS, 10, 32)
	if err != nil {
		log.Println("Error converting poll id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	this_voter, found := v.voterList.Voters[uint(vid64)]
	if found {
		for _, value := range this_voter.VoteHistory {
			if int64(value.PollID) == pid64 {
				c.JSON(http.StatusOK, value)
				return
			}
		}
		log.Println("Poll not found: ", pid64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else {
		log.Println("Voter not found: ", vid64)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func (td *VoterApi) CrashSim(c *gin.Context) {
	panic("Simulating an unexpected crash")
}

func (td *VoterApi) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}

func convertToUnit(x int64) (uint, error) {
	if x < 0 {
		return 0, fmt.Errorf("cannot convert negative int64 to uint")
	}
	return uint(x), nil
}
