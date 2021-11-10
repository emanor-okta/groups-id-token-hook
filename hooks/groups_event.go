package hooks

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

const (
	hookName string = "GroupCreated"
	active   string = "ACTIVE"
	inActive string = "INACTIVE"
)

func SetupEventHook(client *okta.Client) {
	fmt.Println("Setting up group creation event hook")
	configured := checkForHook(client)
	fmt.Printf("Events Hook Already Configured: %v\n", configured)
	if !configured {
		created := createHook(client)
		fmt.Printf("Events Hook Created: %v\n", created)
	}
}

func createHook(client *okta.Client) bool {
	subscriptions := okta.EventSubscriptions{
		Type:  "EVENT_TYPE",
		Items: []string{"group.lifecycle.create"},
	}
	channelConfig := okta.EventHookChannelConfig{
		Uri: Config.HOOK.BASE_URL + "/group",
	}
	channel := okta.EventHookChannel{
		Type:    "HTTP",
		Version: "1.0.0",
		Config:  &channelConfig,
	}
	hook := &okta.EventHook{
		Name:    hookName,
		Channel: &channel,
		Events:  &subscriptions,
	}

	hook, resp, err := client.EventHook.CreateEventHook(context.TODO(), *hook)
	if err != nil {
		fmt.Printf("Error creating event hook: %v\n", err)
		return false
	} else if resp.StatusCode >= 300 {
		fmt.Printf("Error creating event hook: %v\n", resp.Status)
		return false
	}

	if hook.VerificationStatus == "UNVERIFIED" {
		fmt.Println("Need to verify Hook")
		_, resp, err := client.EventHook.VerifyEventHook(context.TODO(), hook.Id)
		if err != nil {
			fmt.Printf("Error verifying event hook: %v\n", err)
			return false
		} else if resp.StatusCode >= 300 {
			fmt.Printf("Error verifying event hook: %v\n", resp.Status)
			return false
		}
		// fmt.Printf("%v\n", hook)
	}

	return true
}

func checkForHook(client *okta.Client) bool {
	hooks, resp, err := client.EventHook.ListEventHooks(context.TODO())
	if err != nil {
		fmt.Printf("Error checking for event hook: %v\n", err)
		return false
	} else if resp.StatusCode >= 300 {
		fmt.Printf("Error checking for event hook: %v\n", resp.Status)
		return false
	}

	for _, hook := range hooks {
		if hook.Name == hookName {
			if hook.Status == active {
				return true
			}
			// if not active make it active
			hook.Status = active
			_, resp, err = client.EventHook.UpdateEventHook(context.TODO(), hook.Id, *hook)
			if err != nil {
				fmt.Printf("Error activating event hook: %v\n", err)
			} else if resp.StatusCode >= 300 {
				fmt.Printf("Error activating event hook: %v\n", resp.Status)
			}
			return true
		}
	}
	return false
}
