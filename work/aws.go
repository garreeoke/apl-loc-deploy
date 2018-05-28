package work

import (
	propeller "applariat.io/propeller/types"
)

const (
	awsInterviewFile = "aws.json"
	titleAccessKeyId = "aws_access_key_id"
	titleAccessKeySecret = "aws_secret_access_key"
	titleTld = "top_level_domain"
	titleRegion = "region"
	titleAzone = "avalability_zone"
	titleClusterSize = "cluster_size"
)

type Aws propeller.LocDeploy

func (aws *Aws) Configure(interview *Interview) error {

	cfg, err := aws.questions(interview)
	if err != nil {
		return err
	}
	aws.Config = cfg

	return nil
}

// interviews reads the interviews file and returns aws config
func (aws *Aws) questions(interview *Interview) (*propeller.AwsConfig, error) {

	var cfg propeller.AwsConfig
	interviewAws := Interview{}
	var err error
	if !interview.Silent {
		interviewFile := InterviewDir + awsInterviewFile
		interviewAws, err = ReadFile(interviewFile)
		if err != nil {
			return &cfg, err
		}
		interviewAws.LocDeploy = interview.LocDeploy
	} else {
		interviewAws = *interview
	}

	cred_value := map[string]string{}
	for i, question := range interviewAws.Questions {
		if ValidateQuestion(interview.Action, "", &question) {
			err := AskQuestion(interview.Silent, interview.Reader, &question)
			if err != nil {
				return &cfg, err
			}
			interviewAws.Questions[i].Answer = question.Answer
			interview.Questions = append(interview.Questions, question)
			switch question.Title {
			case titleAccessKeyId, titleAccessKeySecret:
				cred_value[question.Title] = question.Answer
			case titleTld:
				cfg.TopLevelDomain = question.Answer
			case titleRegion:
				cfg.Region = question.Answer
			case titleAzone:
				cfg.Azone = question.Answer
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
		aws.Cred = propeller.Credential{
			Value: cred_value,
		}
	}
	return &cfg, nil
}