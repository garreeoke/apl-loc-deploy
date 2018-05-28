package work

import (
	"encoding/json"
	propeller "applariat.io/propeller/types"
	propwork "applariat.io/propeller/work"
	propcommon "applariat.io/propeller/common"
	"os"
	"os/exec"
	"fmt"
	"io/ioutil"
	"applariat.io/propeller/kube"
)

const (
	image = "applariat/cluster-manager"
)

func ClusterMgr(interview *Interview) error {
	jobid := ""
	if os.Getenv("APL_LOC_JOBID") == "" {
		var aplResponse AplResponse
		err := json.Unmarshal(interview.RestData.Response, &aplResponse)
		if err != nil {
			return err
		}

		locResponseData, err := json.Marshal(aplResponse.Data)
		if err != nil {
			return err
		}
		var locResponse LocDeployResponse
		err = json.Unmarshal(locResponseData, &locResponse)
		if err != nil {
			return err
		}
		jobid = locResponse.JobId
	} else {
		jobid = os.Getenv("APL_LOC_JOBID")
	}

	// Launch on either local docker or as a k8s job?
	if os.Getenv("APL_LAUNCHER") == "k8s" {
		err := k8Job(jobid)
		if err != nil {
			return err
		}
	} else if os.Getenv("APL_LAUNCHER") == "docker" {
		err := cmdDockerPull()
		if err != nil {
			return err
		}
		err = cmdDockerRun(jobid, interview)
		if err != nil {
			return err
		}
	}

	return nil

}

func k8Job(jobid string) error {
	k8 := kube.K8{
		DeployID: "cluster-mgr",
		Name: "cluster-mgr",
		OnCluster: true,
	}
	err := k8.Auth()
	if err != nil {
		k8.Log.Println("K8 Auth error: ", err)
		return err
	}
	job := propeller.MqJob{
		JobID: jobid,
		Get: true,
		Log: propwork.Logger(fmt.Sprintf("%v(JobID)", jobid)),
		ApiToken: os.Getenv("APL_API_KEY"),
		IgnoreCertValidation: true,
	}
	// Setup the container object to use for the job
	c, err := propcommon.ClusterMgrContainer(image, &job)

	// Get the k8 job definition
	k8Job := k8.ClusterMgrJob(&job, c)
	job.Log.Println("Cluster-mgr job definition created, now submitting ...")
	namespace := os.Getenv("NAMESPACE")
	if namespace != "" {
		k8.Name = namespace
	} else {
		k8.Name = "default"
	}

	err = k8.SubmitJob(true, &k8Job)
	if err != nil {
		k8.DeleteJobPods(true, k8Job.Name)
		_ = k8.DeleteSecret([]string{k8Job.Name})
		return err
	}

	leaveJobPod := false
	if os.Getenv("APL_LEAVE_JOB_POD") == "true" {
		leaveJobPod = true
	}
	k8.DeleteJobPods(leaveJobPod, k8Job.Name)
	_ = k8.DeleteSecret([]string{k8Job.Name})
	job.Log.Println("Cluster-mgr job completed")
	return nil
}

func cmdDockerPull() error {
	fmt.Println("Running docker pull")
	params := []string{
		"pull",
		ClusterMgrImage,
	}

	cmd := exec.Command("docker", params...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}

func cmdDockerRun(jobid string, interview *Interview) error {

	fmt.Println("Running docker run")
	status := make(chan bool, 1)
	job := propeller.MqJob{
		JobID: jobid,
		Get: true,
		Log: propwork.Logger(fmt.Sprintf("%v(JobID)", jobid)),
		ApiToken: os.Getenv("APL_API_KEY"),
		IgnoreCertValidation: true,
	}
	err := propwork.AplJobPayload(&job)
	if err != nil {
		return err
	}
	var cmError error

	var confData []byte

	if interview.Conf != "" {
		confData, err =ioutil.ReadFile("./" + interview.Conf)
		if err != nil {
			return err
		}
	}
	go func(){
		params := []string{
			"run",
			"-e=APL_JOB_ID="+jobid,
			"-e=APL_API="+os.Getenv("APL_API"),
			"-e=APL_API_KEY="+os.Getenv("APL_API_KEY"),
			"-e=APL_SVC_USERNAME="+job.SvcAccount["svc_user"],
			"-e=APL_SVC_PASSWORD="+job.SvcAccount["svc_passwd"],
			"-e=APL_IGNORE_SSL=true",
		}
		if len(confData) > 0 {
			params = append(params, "-e=APL_CONF="+string(confData))
		}

		if os.Getenv("PROPELLER_APL_BUILD_IMAGE_TAG") != "" {
			params = append(params, "-e=PROPELLER_APL_BUILD_IMAGE_TAG="+os.Getenv("PROPELLER_APL_BUILD_IMAGE_TAG"))
		}
		params = append(params, ClusterMgrImage)

		cmd := exec.Command("docker", params...)
		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			cmError = err
			status <- false
		}
		status <- true
		close(status)
	}()

	for x := range status {
		if !x {
			return cmError
		}
	}
	return nil
}
