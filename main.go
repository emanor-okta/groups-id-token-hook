package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"

	"github.com/emanor-okta/group_ids_hook/server/handlers"
	"github.com/emanor-okta/group_ids_hook/server/hooks"
)

var grpMap map[string]string
var client *okta.Client

func main() {
	ch := make(chan int)
	grpMap = make(map[string]string)

	fmt.Println("Server Starting...")
	go startServer()

	fmt.Println("sleeping 5 seconds")
	time.Sleep(time.Second * 5)
	fmt.Println("continue")

	loadGroups()
	handlers.SetGroups(grpMap)

	fmt.Println("Loaded Groups:")
	for name, id := range grpMap {
		fmt.Printf("id: %s, %s\n", id, name)
	}

	hooks.SetupEventHook(client)
	hooks.SetupInlineHook(client)

	wait := <-ch
	fmt.Printf("Should not see %v\n", wait)
}

func loadGroups() {
	var ctx context.Context
	var err error
	ctx, client, err = okta.NewClient(context.TODO(), okta.WithOrgUrl(hooks.Config.CLIENT.ORG_URL), okta.WithToken(hooks.Config.CLIENT.TOKEN))
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
	}
	fmt.Println(*client)
	grps, _, err := client.Group.ListGroups(ctx, nil)
	if err != nil {
		fmt.Printf("Error loading Groups: %v\n", err)
	}

	for _, grp := range grps {
		// fmt.Printf("Name: %s, ID: %s\n", grp.Profile.Name, grp.Id)
		grpMap[grp.Profile.Name] = grp.Id
	}
}

func startServer() {
	http.HandleFunc("/group", handlers.GroupHandler)
	http.HandleFunc("/token", handlers.TokenHandler)

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Server startup failed: %s\n", err)
	}
}
