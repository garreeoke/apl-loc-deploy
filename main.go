package main

import (
	propeller "applariat.io/propeller/types"
	"fmt"
	"github.com/applariat/apl-loc-deploy/work"
	"bufio"
	"os"
	"flag"
)

var validate = []string{"APL_API_KEY","APL_API"}
var answerFile = flag.String("answer-file", "", "Specify answerFile for interview files")
var conf = flag.String("conf-file", "", "Conf file to use in current directory")
var silent = flag.Bool("silent", false, "Install silently")

func main() {

	flag.Parse()

	env := validateEnv()
	if !env {
		return
	}
	interviewFile := work.InterviewDir
	// Figure out if there is an answer file or not
	if *answerFile != "" {
		interviewFile = interviewFile + *answerFile
		if *silent {
			fmt.Println("Silent mode ... ")
			if *answerFile == "" {
				fmt.Println("No answerFile specified for file name")
				return
			}
		}
	} else {
		// If there isn't an answer file, set first file to general
		interviewFile = work.InterviewDir + work.GeneralQuestion
	}

	interview, err := work.ReadFile(interviewFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	interview.Silent = *silent
	interview.Conf = *conf
	interview.Reader = bufio.NewReader(os.Stdin)
	interview.LocDeploy= &propeller.LocDeploy {
		LocalScript: true,
	}

	// Start with general
	err = work.GeneralQuestions(&interview)
	if err != nil {
		fmt.Println(err)
		return
	}

	// If import, launch questions for kubernetes import
	if interview.Action == work.ActionImport {
		err = work.CreateK8(&interview)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = work.AplLocDeploy(&interview)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Response: ", string(interview.RestData.Response))
	if interview.RestData.StatusCode == 400 {
		return
	}

	// Get cluster mgr job and launch on local docker
	err = work.ClusterMgr(&interview)
	if err != nil {
		fmt.Println(err)
	}

}

// validateEnv ... validate environment variables
func validateEnv() bool {

	for _, v := range validate {
		if os.Getenv(v) == "" {
			fmt.Println(v + " environment variable not set")
			return false
		}
	}
	return true
}
