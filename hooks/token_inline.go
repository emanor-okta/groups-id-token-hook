package hooks

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

const (
	inlineType     string = "com.okta.oauth2.tokens.transform"
	inlineHookName string = "Token Inline Hook"
)

func SetupInlineHook(client *okta.Client) {
	fmt.Println("Setting up token inline hook")
	configured := checkForInlineHook(client)
	fmt.Printf("Token Hook Already Configured: %v\n", configured)
	if !configured {
		created := createInlineHook(client)
		fmt.Printf("Token Hook Created: %v\n", created)
	}
}

func createInlineHook(client *okta.Client) bool {
	channelConfig := okta.InlineHookChannelConfig{
		Uri: Config.HOOK.BASE_URL + "/token",
	}
	channel := okta.InlineHookChannel{
		Type:    "HTTP",
		Version: "1.0.0",
		Config:  &channelConfig,
	}
	hook := &okta.InlineHook{
		Name:    inlineHookName,
		Channel: &channel,
		Type:    inlineType,
		Version: "1.0.0",
	}

	_, resp, err := client.InlineHook.CreateInlineHook(context.TODO(), *hook)
	if err != nil {
		fmt.Printf("Error creating inline hook: %v\n", err)
		return false
	} else if resp.StatusCode >= 300 {
		fmt.Printf("Error creating inline hook: %v\n", resp.Status)
		return false
	}

	return true
}

func checkForInlineHook(client *okta.Client) bool {
	hooks, resp, err := client.InlineHook.ListInlineHooks(context.TODO(), &query.Params{Type: inlineType})
	if err != nil {
		fmt.Printf("Error checking for inline hook: %v\n", err)
		return false
	} else if resp.StatusCode >= 300 {
		fmt.Printf("Error checking for inline hook: %v\n", resp.Status)
		return false
	}

	for _, hook := range hooks {
		if hook.Name == inlineHookName {
			if hook.Status == active {
				return true
			}
			// if not active make it active
			hook.Status = active
			_, resp, err = client.InlineHook.ActivateInlineHook(context.TODO(), hook.Id)
			if err != nil {
				fmt.Printf("Error activating inline hook: %v\n", err)
			} else if resp.StatusCode >= 300 {
				fmt.Printf("Error activating inline hook: %v\n", resp.Status)
			}
			return true
		}
	}
	return false
}
