package work

import (
	"fmt"
	propeller "applariat.io/propeller/types"
	propWork "applariat.io/propeller/work"
	"os"
	"net/http"
	"encoding/json"
	"errors"
	"strconv"
	"crypto/tls"
	"applariat.io/propeller/kube"
)

const (
	k8sInterviewFile   = "k8s.json"
	titleFqdn          = "kube_fqdn"
	titlePort          = "kube_port"
	titleAuthType      = "auth_type"
	titleBasicUserName = "user"
	titleBasicPassword = "password"
	titleCaCert        = "ca_cert"
	titleClientKey     = "client_key"
	titleClientCert    = "client_cert"
	kubeCredType       = "kubernetes"
)

func CreateK8(interview *Interview) error {

	interview.LocDeploy.Cluster = propeller.KubeCluster{
		CredentialType: kubeCredType,
	}
	// Create credential and get the id
	err := questions(interview)
	if err != nil {
		return err
	}

	return nil
}

func questions(interview *Interview) error {

	// PKS
	if interview.LocDeploy.Name == "pks_auto" {
		fmt.Println("PKS_AUTO_DETECTED")
		k8 := kube.K8{
			DeployID: "cluster-mgr",
			Name: "cluster-mgr",
			OnCluster: true,
		}
		err := k8.Auth(false)
		if err != nil {
			k8.Log.Println("K8 Auth error: ", err)
			return err
		}
		interview.LocDeploy.Name = "pks-" + k8.NodeLabel("kubernetes.io/hostname") + "-" + k8.NodeLabel("bosh.id")
		fmt.Println("PKS_AUTO_NEW_NAME: ", interview.LocDeploy.Name)
		// Would love for it to be
		//interview.LocDeploy.Name = k8.NodeLabel("pks.io/cluster-name")
	}
	interview.LocDeploy.Cluster.Name = interview.LocDeploy.Name
	aplCred := propeller.AplCred{
		Name: "k8-ext-" + interview.LocDeploy.Name,
		CredType: kubeCredType,
	}
	interviewK8s := Interview{}
	var err error
	if !interview.Silent {
		interviewFile := InterviewDir + k8sInterviewFile
		interviewK8s, err = ReadFile(interviewFile)
		if err != nil {
			return err
		}
		interviewK8s.LocDeploy = interview.LocDeploy
	} else {
		interviewK8s = *interview
	}
	cred_value := map[string]string{}
	credential := propeller.Credential{}
	for i, question := range interviewK8s.Questions {
		if question.Section == K8sType {
			if ValidateQuestion(interview.Action, credential.Type, &question) {
				err := AskQuestion(interview.Silent, interview.Reader, &question)
				if err != nil {
					return err
				}
				interviewK8s.Questions[i].Answer = question.Answer
				interview.Questions = append(interview.Questions, question)
				switch question.Title {
				case titleFqdn:
					interview.LocDeploy.Cluster.Fqdn = "https://" + question.Answer
				case titlePort:
					port,_ := strconv.Atoi(question.Answer)
					interview.LocDeploy.Cluster.Port = port
				case titleAuthType:
					credential.Type = question.Answer
				case titleBasicUserName, titleBasicPassword, titleCaCert, titleClientCert, titleClientKey:
					cred_value[question.Title] = question.Answer
				}
			}
		}
	}
	credential.Value = cred_value
	aplCred.Credentials = credential

	// Post the credential and get back id
	err = postCred(&aplCred, interview)
	if err != nil {
		return err
	}

	return nil
}

func postCred(aplCred *propeller.AplCred, interview *Interview) error {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true,
		},
	}

	client := http.Client{Transport: tr}
	rc := propeller.RestClient{
		Client: &client,
	}

	payloadMap := map[string]interface{}{
		"data": aplCred,
	}
	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}

	rd := propeller.RestData{
		Verb: "POST",
		URL:     os.Getenv("APL_API") + "/credentials",
		ApiKey:  os.Getenv("APL_API_KEY"),
		Client:  &rc,
		Payload: payload,
	}

	err = propWork.AplAPI(&rd)
	if err != nil {
		return err
	}


	if rd.StatusCode == 400 {
		fmt.Println("Post Cred Error: ", string(rd.Response))
		return errors.New(string(rd.Response))
	}

	var data map[string]interface{}
	err = json.Unmarshal(rd.Response, &data)
	if err != nil {
		return err
	}


	interview.LocDeploy.Cluster.CredentialId = data["data"].(string)

	return nil
}
