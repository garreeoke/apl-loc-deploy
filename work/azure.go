package work

import (
	propeller "applariat.io/propeller/types"
)

const (
	azureInterviewFile = "azure.json"
	titleAzureLocation = "azure_location"
	titleAzureResourceGroup = "resource_group"
)

type Acs propeller.LocDeploy

func (acs *Acs) Configure(interview *Interview) error {

	cfg, err := acs.questions(interview)
	if err != nil {
		return err
	}
	acs.Config = cfg
	return nil
}

func (acs *Acs) questions(interview *Interview) (*propeller.AcsConfig, error) {

	var cfg propeller.AcsConfig
	interviewAcs := Interview{}
	var err error
	if !interview.Silent {
		interviewFile := InterviewDir + azureInterviewFile
		interviewAcs, err = ReadFile(interviewFile)
		if err != nil {
			return &cfg, err
		}
		interviewAcs.LocDeploy = interview.LocDeploy
	} else {
		interviewAcs = *interview
	}
	cred_value := map[string]string{}
	for i, question := range interviewAcs.Questions {
		if ValidateQuestion(interview.Action, "", &question) {
			err := AskQuestion(interview.Silent, interview.Reader, &question)
			if err != nil {
				return &cfg, err
			}
			interviewAcs.Questions[i].Answer = question.Answer
			interview.Questions = append(interview.Questions, question)
			switch question.Title {
			case titleAzureLocation:
				cfg.Location = question.Answer
			case titleAzureResourceGroup:
				cfg.ResourceGroup = question.Answer
			case titleClusterSize:
				cfg.Size = question.Answer
			}
		}
	}
	if interview.Action == ActionImport {
		cfg.Size = "S"
	}
	if len(cred_value) > 0 {
		// Add cred value to credential
		acs.Cred = propeller.Credential{
			Value: cred_value,
		}
	}
	return &cfg, nil
}