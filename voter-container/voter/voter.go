package voter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "todo:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type VoterPoll struct {
	PollID   uint
	VoteDate time.Time
}

type PublicVoterPoll struct {
	PollID   uint      `json:"pollID"`
	VoteDate time.Time `json:"voteDate"`
}

type Voter struct {
	VoterID     uint
	FirstName   string
	LastName    string
	VoteHistory []VoterPoll
}

type VoterList struct {
	// Voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
	cache
}

func NewVoter(id uint, fn, ln string) *Voter {
	return &Voter{
		VoterID:     id,
		FirstName:   fn,
		LastName:    ln,
		VoteHistory: []VoterPoll{},
	}
}

func New() (*VoterList, error) {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoterList, error) {
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	ctx := context.Background()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error() + "cache might not be available, continuing...")
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	return &VoterList{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

//------------------------------------------------------------
// REDIS HELPERS
//------------------------------------------------------------

func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

// Helper to return a ToDoItem from redis provided a key
func (list *VoterList) getItemFromRedis(key string, item *Voter) error {

	//Lets query redis for the item, note we can return parts of the
	//json structure, the second parameter "." means return the entire
	//json structure
	itemObject, err := list.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	//JSONGet returns an "any" object, or empty interface,
	//we need to convert it to a byte array, which is the
	//underlying type of the object, then we can unmarshal
	//it into our ToDoItem struct
	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

// AddItem accepts a ToDoItem and adds it to the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must not already exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if so, return an error
//
// Postconditions:
//
//	 (1) The item will be added to the DB
//		(2) The DB file will be saved with the item added
//		(3) If there is an error, it will be returned
func (t *VoterList) AddItem(item Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	redisKey := redisKeyFromId(int(item.VoterID))
	var existingItem Voter
	if err := t.getItemFromRedis(redisKey, &existingItem); err == nil {
		return errors.New("item already exists")
	}

	//Add item to database with JSON Set
	if _, err := t.jsonHelper.JSONSet(redisKey, ".", item); err != nil {
		return err
	}

	//If everything is ok, return nil for the error
	return nil
}

func (t *VoterList) GetAllItems() ([]Voter, error) {

	//Now that we have the DB loaded, lets crate a slice
	var toDoList []Voter
	var toDoItem Voter

	//Lets query redis for all of the items
	pattern := RedisKeyPrefix + "*"
	ks, _ := t.cacheClient.Keys(t.context, pattern).Result()
	for _, key := range ks {
		err := t.getItemFromRedis(key, &toDoItem)
		if err != nil {
			return nil, err
		}
		toDoList = append(toDoList, toDoItem)
	}

	return toDoList, nil
}

// func (t *VoterList) GetAllPolls() (id int, error) {
// 	// var allPolls []voterPoll
// 	// allVoters, err := t.GetAllItems()
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// for _, value := range allVoters {
// 	// 	allPolls = append(allPolls, value.VoteHistory...)
// 	// }
// 	// return allPolls, nil
// 	currentVoter = t.GetItem(id)
// 	return currentVoter.VoteHistory
// }

func (t *VoterList) GetItem(id int) (Voter, error) {

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	var item Voter
	pattern := redisKeyFromId(id)
	err := t.getItemFromRedis(pattern, &item)
	if err != nil {
		return Voter{}, err
	}

	return item, nil
}

func (t *VoterList) UpdateItem(item Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	redisKey := redisKeyFromId(int(item.VoterID))
	var existingItem Voter
	if err := t.getItemFromRedis(redisKey, &existingItem); err != nil {
		return errors.New("item does not exist")
	}

	//Add item to database with JSON Set.  Note there is no update
	//functionality, so we just overwrite the existing item
	if _, err := t.jsonHelper.JSONSet(redisKey, ".", item); err != nil {
		return err
	}

	//If everything is ok, return nil for the error
	return nil
}

func (t *VoterList) DeleteItem(id int) error {

	pattern := redisKeyFromId(id)
	numDeleted, err := t.cacheClient.Del(t.context, pattern).Result()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent item")
	}

	return nil
}

func (v *Voter) AddPoll(pollID uint) {
	v.VoteHistory = append(v.VoteHistory, VoterPoll{PollID: pollID, VoteDate: time.Now()})
}

func (v *Voter) AddPollSpecific(pollID uint, date time.Time) {
	v.VoteHistory = append(v.VoteHistory, VoterPoll{PollID: pollID, VoteDate: date})
}

func (v *Voter) ToJson() string {
	b, _ := json.Marshal(v)
	return string(b)
}
