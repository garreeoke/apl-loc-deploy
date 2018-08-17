package work

import (
	propeller "applariat.io/propeller/types"
	"bufio"
)

type Interview struct {
	Questions  []Question           `json:"questions,omitempty"`
	Action     string               `json:"action,omitempty"`
	Credential bool                 `json:"credential,omitempty"`
	Reader     *bufio.Reader        `json:"-"`
	RestData   *propeller.RestData  `json:"-"`
	LocDeploy  *propeller.LocDeploy `json:"-"`
	ObjectId   string               `json:"-"`
	Conf       string               `json:"-"`
	//Registry   string			    `json:"-"`
	//Annotations string			    `json:"-"`
	Silent     bool                 `json:"-"`
}

type Question struct {
	Actions  []string `json:"actions,omitempty"`
	Section  string   `json:"section,omitempty"`
	Title    string   `json:"title,omitempty"`
	Type     string   `json:"type,omitempty"`
	Text     string   `json:"text,omitempty"`
	Accepted []string `json:"accepted,omitempty"`
	Answer   string   `json:"answer,omitempty"`
}

type Provider interface {
	Configure(*Interview) error
}

type AplResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type LocDeployResponse struct {
	Services []string `json:"loc_deploy_services,omitempty"`
	Id       string   `json:"loc_deploy_id,omitempty"`
	JobId    string   `json:"job_id,omitempty"`
}
