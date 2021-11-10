package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	dataTypes "github.com/emanor-okta/group_ids_hook/server/types"
)

var grpMap map[string]string

func SetGroups(groups map[string]string) {
	grpMap = groups
}

func HandleVerify(res http.ResponseWriter, req *http.Request) {
	verification := req.Header.Get("x-okta-verification-challenge")
	res.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["verification"] = verification
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("Error Marshalling /verify Json Response: %v\n", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write(nil)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(jsonResp)
}

func GroupHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		HandleVerify(res, req)
		return
	}

	var data dataTypes.GroupCreated
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding Event: %v\n", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// fmt.Printf("%+v\n", data)
	for _, event := range data.Data.Events {
		if event.EventType != "group.lifecycle.create" || event.Outcome.Result != "SUCCESS" {
			continue
		}

		for _, target := range event.Target {
			if target.Type != "UserGroup" {
				continue
			}
			grpMap[target.DisplayName] = target.Id
		}
	}

	res.WriteHeader(http.StatusOK)
	res.Write(nil)
}

func TokenHandler(res http.ResponseWriter, req *http.Request) {
	var data dataTypes.TokenRequest
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding Token Request: %v\n", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// fmt.Printf("%+v\n", data)
	var v []string
	for _, grp := range data.Data.Claims.Groups {
		v = append(v, grpMap[grp])
	}
	cmdVal := dataTypes.Value{
		Op:    "add",
		Path:  "/claims/groupIds",
		Value: v,
	}
	values := dataTypes.Values{
		Type:  "com.okta.identity.patch",
		Value: []*dataTypes.Value{&cmdVal},
	}
	reply := dataTypes.TokenResponse{
		Commands: []*dataTypes.Values{&values},
	}

	jsonResp, err := json.Marshal(reply)
	if err != nil {
		fmt.Printf("Error marhalling token response: %v\n", err)
	}
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonResp)
}
