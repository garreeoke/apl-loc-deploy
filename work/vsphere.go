package work

import (
	propeller "applariat.io/propeller/types"
)

const (
	vsphereInterviewFile = "vsphere.json"
	titleVsphereDatastore = "datastore"
	titleVsphereDisktype = "diskformat"
)

type Vsphere propeller.LocDeploy

func (vsphere *Vsphere) Configure(interview *Interview) error {

	cfg, err := vsphere.questions(interview)
	if err != nil {
		return err
	}
	vsphere.Config = cfg
	return nil
}

func (vsphere *Vsphere) questions(interview *Interview) (*propeller.VsphereConfig, error) {

	var cfg propeller.VsphereConfig
	interviewVsphere := Interview{}
	var err error
	if !interview.Silent {
		interviewFile := InterviewDir + vsphereInterviewFile
		interviewVsphere, err = ReadFile(interviewFile)
		if err != nil {
			return &cfg, err
		}
		interviewVsphere.LocDeploy = interview.LocDeploy
	} else {
		interviewVsphere = *interview
	}
	for i, question := range interviewVsphere.Questions {
		if question.Section == propeller.ProviderVsphere {
			if ValidateQuestion(interview.Action, "", &question) {
				err := AskQuestion(interview.Silent, interview.Reader, &question)
				if err != nil {
					return &cfg, err
				}
				interviewVsphere.Questions[i].Answer = question.Answer
				interview.Questions = append(interview.Questions, question)
				switch question.Title {
				case titleVsphereDatastore:
					cfg.Datastore = question.Answer
				case titleVsphereDisktype:
					cfg.DiskFormat = question.Answer
				}
			}
		}
	}
	cfg.Size = "S"
	return &cfg, nil
}
