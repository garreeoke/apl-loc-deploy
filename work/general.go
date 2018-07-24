package work

import (
	propeller "applariat.io/propeller/types"
	"fmt"
	"errors"
)

func GeneralQuestions(interview *Interview) error {

	fmt.Println("General Questions")
	dns := propeller.LocDNS{}
	// General interview interview
	for i, question := range interview.Questions {
		// Format the question
		if question.Section == GeneralType {
			if ValidateQuestion(interview.Action, "", &question) {
				err := AskQuestion(interview.Silent, interview.Reader, &question)
				if err != nil {
					return err}
				interview.Questions[i].Answer = question.Answer
				switch question.Title {
				case "loc_deploy_name":
					interview.LocDeploy.Name = question.Answer
				case "registry_url":
					interview.LocDeploy.Registry = propeller.Registry{
						External: true,
						Url: question.Answer,
					}
				case "action":
					interview.Action = question.Answer
					if interview.Action == ActionCreate {
						interview.LocDeploy.AplManaged = true
					} else if interview.Action == ActionImport {
						interview.LocDeploy.AplManaged = false
					}
				case "loc_deploy_id":
					interview.ObjectId = question.Answer
				case "confirm_delete":
					if question.Answer == "no" {
						return errors.New("Goodbye ... ")
					}
				case "provider":
					interview.LocDeploy.Type = question.Answer
				case "provider_credential_id":
					// Get the credential id implement later ...
					if question.Answer != ActionCreate {
						// TODO lookup cred id
						interview.LocDeploy.CredentialId = question.Answer
					} else {
						// TODO Create new credential
					}
				case "dns_enabled":
					if question.Answer == "true" {
						dns.Enabled = true
					}
				case "hosted_zone":
					dns.HostedZone = question.Answer
				case "allowed_wordloads":
					// TODO map workloads
					// Do nothing right now ... if defaults to all that would be sweet :-)
				}
			}
		}
	}

	if dns.Enabled {
		interview.LocDeploy.DNS = dns
	}
	return nil
}
