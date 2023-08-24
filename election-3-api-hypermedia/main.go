package main

import (
	"flag"
	"fmt"
	"voter-api/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")
	flag.Parse()
}

func main() {

	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler := api.NewVoterApi()
	apiHandler.AddVoter(1, "John", "Doe")
	apiHandler.AddPoll(1, 1)
	apiHandler.AddPoll(1, 2)
	apiHandler.AddPoll(1, 3)
	apiHandler.AddVoter(2, "Amirali", "Sajadi")
	apiHandler.AddPoll(2, 1)

	r.GET("/voters", apiHandler.GetVoterList)
	r.GET("/voters/:id", apiHandler.GetVoter)
	r.POST("/voters", apiHandler.AddVoterGoodHandler)
	r.PUT("/voters/:id", apiHandler.UpdateVoter)
	r.DELETE("/voters/:id", apiHandler.DeleteVoter)
	r.GET("/voters/:id/polls", apiHandler.GetPolls)
	r.GET("/voters/:id/polls/:pollid", apiHandler.GetPoll)
	// r.POST("/voters/:id/polls/:pollid", apiHandler.AddPollHandlerBadImplementation)
	r.POST("/voters/:id/polls", apiHandler.AddPollHandlerGoodImplementation)
	r.PUT("/voters/:id/polls/:pollid", apiHandler.UpdatePoll)
	r.DELETE("/voters/:id/polls/:pollid", apiHandler.DeletePoll)
	r.GET("/health", apiHandler.HealthCheck)
	r.GET("/crash", apiHandler.CrashSim)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)

}
