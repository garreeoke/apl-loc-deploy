package work

import (
	propeller "applariat.io/propeller/types"
	propWork "applariat.io/propeller/work"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"crypto/tls"
)

func ProviderQuestions(interview *Interview) error {

	fmt.Println("Provider Questions")
	l, err := json.Marshal(interview.LocDeploy)
	if err != nil {
		return err
	}
	switch interview.LocDeploy.Type {
	case propeller.ProviderAws:
		var aws Aws
		json.Unmarshal(l, &aws)
		err = action(&aws, interview)
		if err != nil {
			return err
		}
		l, err = json.Marshal(aws)
		if err != nil {
			return err
		}
	case propeller.ProviderGke:
		var gke Gke
		json.Unmarshal(l, &gke)
		err = action(&gke, interview)
		if err != nil {
			return err
		}
		l, err = json.Marshal(gke)
		if err != nil {
			return err
		}
	case propeller.ProviderAzure:
		var acs Acs
		json.Unmarshal(l, &acs)
		err = action(&acs, interview)
		if err != nil {
			return err
		}
		l, err = json.Marshal(acs)
		if err != nil {
			return err
		}
	case propeller.ProviderVsphere:
		var vsphere Vsphere
		json.Unmarshal(l, &vsphere)
		err = action(&vsphere, interview)
		if err != nil {
			return err
		}
		l, err = json.Marshal(vsphere)
		if err != nil {
			return err
		}
	case propeller.ProviderMetal :
		fmt.Println("Loc_deploy type metal, no more questsions")
	default:
		return errors.New("Unknown provider type: " + interview.LocDeploy.Type)
	}

	err = json.Unmarshal(l, &interview.LocDeploy)
	if err != nil {
		return err
	}

	return nil
}

func action(p Provider, interview *Interview) error {

	var err error
	switch interview.Action {
	case ActionCreate, ActionImport:
		err = p.Configure(interview)
	case ActionDelete:
		// Go straight to delete
	}
	if err != nil {
		return err
	}

	return nil
}

//AplLocDeploy ... submit payload to apl
func AplLocDeploy(interview *Interview) error {

	rd := propeller.RestData{
		URL:    os.Getenv("APL_API") + "/loc_deploys",
		ApiKey: os.Getenv("APL_API_KEY"),
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: tr}

	rc := propeller.RestClient{
		Client: &client,
	}
	rd.Client = &rc

	var err error
	switch interview.Action {
	case ActionCreate, ActionImport:
		// Provider specific interviews
		err = ProviderQuestions(interview)
		if err != nil {
			fmt.Println(err)
			return err
		}
		rd.Verb = "POST"
		// Send to apl
		payloadMap := map[string]interface{}{
			"data": interview.LocDeploy,
		}
		payload, err := json.Marshal(payloadMap)
		if err != nil {
			return err
		}
		fmt.Println("PAYLOAD: ", string(payload))
		rd.Payload = payload
	case ActionDelete:
		rd.Verb = "DELETE"
		rd.URL = rd.URL + "/" + interview.ObjectId
	}
	// Save the interview questions to file
	if interview.Action == ActionCreate || interview.Action == ActionImport {
		err := SaveQuestions(interview)
		if err != nil {
			fmt.Println(err)
		}
	}
	interview.RestData = &rd
	err = propWork.AplAPI(interview.RestData)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully posted loc_deploy to apl_api_url: %v \n", interview.RestData.URL)
	return nil

}
