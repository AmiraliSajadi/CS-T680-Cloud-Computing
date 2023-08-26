package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"voter-container/voter"

	"github.com/gin-gonic/gin"
)

type VoterApi struct {
	list *voter.VoterList
}

func New() (*VoterApi, error) {
	dbHandler, err := voter.New()
	if err != nil {
		return nil, err
	}
	return &VoterApi{list: dbHandler}, nil
}

//	func (v *VoterApi) AddVoter(voterID uint, firstName, lastName string) {
//		v.voterList.Voters[voterID] = *voter.NewVoter(voterID, firstName, lastName)
//	}

// func (v *VoterApi) AddVoterGoodHandler(c *gin.Context) {
// 	var this_voter voter.Voter
// 	if err := c.ShouldBindJSON(&this_voter); err != nil {
// 		log.Println("Error binding JSON: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	_, found := v.voterList.Voters[uint(this_voter.VoterID)]
// 	if found {
// 		log.Println("Voter ID already exists: ", this_voter.VoterID)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	} else {
// 		v.voterList.Voters[this_voter.VoterID] = *voter.NewVoter(
// 			this_voter.VoterID,
// 			this_voter.FirstName,
// 			this_voter.LastName,
// 		)
// 	}
// }

func (v *VoterApi) AddVoterGoodHandler(c *gin.Context) {
	var currentVoter voter.Voter

	if err := c.ShouldBindJSON(&currentVoter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.list.AddItem(currentVoter); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	c.JSON(http.StatusOK, currentVoter)
}

// func (v *VoterApi) AddPoll(voterID, pollID uint) {
// 	voter := v.voterList.Voters[voterID]
// 	voter.AddPoll(pollID)
// 	v.voterList.Voters[voterID] = voter
// }

// // This one does not use the payload (only uses the id and pollid in url)
// func (v *VoterApi) AddPollHandlerBadImplementation(c *gin.Context) {
// 	idS := c.Param("id")
// 	id64, err := strconv.ParseUint(idS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting id to uint64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	pidS := c.Param("pollid")
// 	pid64, err := strconv.ParseInt(pidS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting poll id to int64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	this_voter, found := v.voterList.Voters[uint(id64)]
// 	if found {
// 		// check if poll id already exists
// 		for _, value := range this_voter.VoteHistory {
// 			if int64(value.PollID) == pid64 {
// 				log.Println("The pollid for this user already exists", err)
// 				c.AbortWithStatus(http.StatusBadRequest)
// 				return
// 			}
// 		}
// 		// add the poll
// 		pidUint, err := convertToUnit(pid64)
// 		if err != nil {
// 			c.AbortWithStatus(http.StatusBadRequest)
// 			return
// 		} else {
// 			this_voter.AddPoll(pidUint)
// 			v.voterList.Voters[uint(id64)] = this_voter
// 			log.Println("poll added to the user")
// 			return
// 		}
// 	} else {
// 		log.Println("Voter not found: ", id64)
// 		c.AbortWithStatus(http.StatusNotFound)
// 		return
// 	}
// }

// the version that uses the payload data
func (v *VoterApi) AddPollHandlerGoodImplementation(c *gin.Context) {
	idS := c.Param("pid")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	currentVoter, err := v.list.GetItem(int(id64))
	if err != nil {
		log.Println("Error getting voter: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// add the poll
	// currentVoter.AddPoll(poll.PollID)
	var newPoll voter.VoterPoll

	if err := c.BindJSON(&newPoll); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	currentVoter.VoteHistory = append(currentVoter.VoteHistory, newPoll)
	err = v.list.UpdateItem(currentVoter)
	if err != nil {
		log.Println("Error updating voter: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println("Poll added to the user")
	c.JSON(http.StatusOK, currentVoter)
}

func (v *VoterApi) GetVoter(c *gin.Context) {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	todoItem, err := v.list.GetItem(int(id64))
	if err != nil {
		log.Println("Item not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	c.JSON(http.StatusOK, todoItem)
}

// func (v *VoterApi) UpdateVoter(c *gin.Context) {
// 	idS := c.Param("id")
// 	id64, err := strconv.ParseUint(idS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting id to uint64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	// Check if the voter exists
// 	this_voter, found := v.voterList.Voters[uint(id64)]
// 	if found {
// 		// do stuff
// 		var updatedVoter voter.Voter

// 		if err := c.ShouldBindJSON(&updatedVoter); err != nil {
// 			log.Println("Error binding JSON: ", err)
// 			c.AbortWithStatus(http.StatusBadRequest)
// 			return
// 		} else {
// 			// URL and payload data don't match
// 			if id64 != uint64(updatedVoter.VoterID) {
// 				log.Println("The ID in URL does not match the one in payload", updatedVoter.VoterID)
// 				c.AbortWithStatus(http.StatusBadRequest)
// 				return
// 			} else {
// 				this_voter.FirstName = updatedVoter.FirstName
// 				this_voter.LastName = updatedVoter.LastName
// 				v.voterList.Voters[uint(id64)] = this_voter
// 				log.Println("User Updated")
// 				c.AbortWithStatus(http.StatusOK)
// 			}
// 		}
// 	} else {
// 		log.Println("Voter not found: ", id64)
// 		c.AbortWithStatus(http.StatusNotFound)
// 		return
// 	}
// }

func (v *VoterApi) UpdateVoter(c *gin.Context) {
	var todoItem voter.Voter
	if err := c.ShouldBindJSON(&todoItem); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.list.UpdateItem(todoItem); err != nil {
		log.Println("Error updating item: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, todoItem)
}

func (v *VoterApi) DeleteVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, _ := strconv.ParseInt(idS, 10, 32)

	if err := v.list.DeleteItem(int(id64)); err != nil {
		log.Println("Error deleting item: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

// func (v *VoterApi) UpdatePoll(c *gin.Context) {
// 	idS := c.Param("id")
// 	id64, err := strconv.ParseUint(idS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting id to uint64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	pidS := c.Param("pollid")
// 	pid64, err := strconv.ParseInt(pidS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting poll id to int64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	var poll voter.PublicVoterPoll
// 	if err := c.ShouldBindJSON(&poll); err != nil {
// 		log.Println("Error binding JSON: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	this_voter, found := v.voterList.Voters[uint(id64)]
// 	if pid64 == int64(poll.PollID) && found {
// 		// do stuff
// 		for i, value := range this_voter.VoteHistory {
// 			if int64(value.PollID) == pid64 {
// 				this_voter.VoteHistory[i].PollID = poll.PollID
// 				this_voter.VoteHistory[i].VoteDate = poll.VoteDate
// 				v.voterList.Voters[uint(id64)] = this_voter
// 				log.Println("Poll updated")
// 				c.AbortWithStatus(http.StatusOK)
// 				return
// 			} else {
// 				log.Println("The PooID does not exist", err)
// 				c.AbortWithStatus(http.StatusBadRequest)
// 				return
// 			}
// 		}
// 	} else {
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// }

// func (v *VoterApi) DeletePoll(c *gin.Context) {
// 	idS := c.Param("id")
// 	id64, err := strconv.ParseUint(idS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting id to uint64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}
// 	pidS := c.Param("pollid")
// 	pid64, err := strconv.ParseInt(pidS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting poll id to int64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	this_voter := v.voterList.Voters[uint(id64)]
// 	delete_index := -1
// 	for i, value := range this_voter.VoteHistory {
// 		if int64(value.PollID) == pid64 {
// 			delete_index = i
// 			break
// 		}
// 	}
// 	this_voter.VoteHistory = append(this_voter.VoteHistory[:delete_index], this_voter.VoteHistory[delete_index+1:]...)
// 	v.voterList.Voters[uint(id64)] = this_voter
// 	log.Println("Poll deleted")
// 	c.AbortWithStatus(http.StatusOK)
// }

// func (v *VoterApi) GetVoterJson(voterID uint) string {
// 	voter := v.voterList.Voters[voterID]
// 	return voter.ToJson()
// }

func (v *VoterApi) GetVoterList(c *gin.Context) {
	todoList, err := v.list.GetAllItems()
	if err != nil {
		log.Println("Error Getting All Items: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if todoList == nil {
		// todoList = make([]voter.Voter, 0)
		todoList = make([]voter.Voter, 0)
	}

	c.JSON(http.StatusOK, todoList)
}

// func (v *VoterApi) GetVoterListJson() string {
// 	b, _ := json.Marshal(v.voterList)
// 	return string(b)
// }

func (v *VoterApi) GetPolls(c *gin.Context) {
	idS := c.Param("id")
	id64, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to uint64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	a, err := v.list.GetItem(int(id64))
	if err != nil {
		fmt.Println("Error getting item:", err)
		c.JSON(http.StatusInternalServerError, nil) // Return an appropriate HTTP status code
		return
	}
	c.JSON(http.StatusOK, a.VoteHistory)
}

func (v *VoterApi) GetPoll(c *gin.Context) {
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

	a, err := v.list.GetItem(int(id64))
	if err != nil {
		fmt.Println("Error getting item:", err)
		c.JSON(http.StatusInternalServerError, nil) // Return an appropriate HTTP status code
		return
	}
	c.JSON(http.StatusOK, a.VoteHistory[pid64])
}

// func (v *VoterApi) GetPoll(c *gin.Context) {
// 	vidS := c.Param("id")
// 	vid64, err := strconv.ParseInt(vidS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting voter id to int64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	pidS := c.Param("pollid")
// 	pid64, err := strconv.ParseInt(pidS, 10, 32)
// 	if err != nil {
// 		log.Println("Error converting poll id to int64: ", err)
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		return
// 	}

// 	this_voter, found := v.voterList.Voters[uint(vid64)]
// 	if found {
// 		for _, value := range this_voter.VoteHistory {
// 			if int64(value.PollID) == pid64 {
// 				c.JSON(http.StatusOK, value)
// 				return
// 			}
// 		}
// 		log.Println("Poll not found: ", pid64)
// 		c.AbortWithStatus(http.StatusNotFound)
// 		return
// 	} else {
// 		log.Println("Voter not found: ", vid64)
// 		c.AbortWithStatus(http.StatusNotFound)
// 		return
// 	}
// }

// func (td *VoterApi) CrashSim(c *gin.Context) {
// 	panic("Simulating an unexpected crash")
// }

// func (td *VoterApi) HealthCheck(c *gin.Context) {
// 	c.JSON(http.StatusOK,
// 		gin.H{
// 			"status":             "ok",
// 			"version":            "1.0.0",
// 			"uptime":             100,
// 			"users_processed":    1000,
// 			"errors_encountered": 10,
// 		})
// }

func convertToUnit(x int64) (uint, error) {
	if x < 0 {
		return 0, fmt.Errorf("cannot convert negative int64 to uint")
	}
	return uint(x), nil
}
