package work

import (
	propeller "applariat.io/propeller/types"
)

const (
	gkeInterviewFile = "gke.json"
	titleGoogleProject = "google_project"
	titleGoogleRegion = "google_region"
	titleGoogleZone = "google_zone"
)

type Gke propeller.LocDeploy

func (gke *Gke) Configure(interview *Interview) error {

	cfg, err := gke.questions(interview)
	if err != nil {
		return err
	}
	gke.Config = cfg
	return nil
}

func (gke *Gke) questions(interview *Interview) (*propeller.GkeConfig, error) {

	var cfg propeller.GkeConfig
	interviewGke := Interview{}
	var err error
	if !interview.Silent {
		interviewFile := InterviewDir + gkeInterviewFile
		interviewGke, err = ReadFile(interviewFile)
		if err != nil {
			return &cfg, err
		}
		interviewGke.LocDeploy = interview.LocDeploy
	}
	cred_value := map[string]string{}
	for i, question := range interviewGke.Questions {
		if ValidateQuestion(interview.Action, "", &question) {
			err := AskQuestion(interview.Silent, interview.Reader, &question)
			if err != nil {
				return &cfg, err
			}
			interviewGke.Questions[i].Answer = question.Answer
			interview.Questions = append(interview.Questions, question)
			switch question.Title {
			case titleGoogleProject:
				cfg.Project = question.Answer
			case titleGoogleRegion:
				cfg.Region = question.Answer
			case titleGoogleZone:
				cfg.Zone = question.Answer
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
		gke.Cred = propeller.Credential{
			Value: cred_value,
		}
	}
	return &cfg, nil
}