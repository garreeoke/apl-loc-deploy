package work

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"bufio"
)

const (
	InterviewDir      = "interviews/"
	InterviewFileType = ".json"
	ActionCreate      = "create"
	ActionImport      = "import"
	ActionDelete      = "delete"
	ActionCreateCred  = "create_credential"
	ClusterMgrImage   = "applariat/cluster-manager:develop"
	GeneralQuestion   = "general.json"
	GeneralType = "general"
	K8sType = "k8s"
)

func ReadFile(fileName string) (Interview, error) {

	fmt.Println("Reading file: ", fileName)
	var i Interview

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading interview file (%v): %v ", fileName, err)
		return i, err
	}
	err = json.Unmarshal(file, &i)
	if err != nil {
		fmt.Printf("Error reading interview file (%v): %v ", fileName, err)
		return i, err
	}
	return i, nil
}

// ValidateQuestion makes sure we need to ask the question
func ValidateQuestion(matchAction, matchType string, question *Question) bool {

	ask := false
	// If it is silent, don't ask the question
	if len(question.Actions) == 0 {
		ask = true
	} else {
		actionMatch := false
		for _, a := range question.Actions {
			if a == matchAction {
				actionMatch = true
				break
			}
		}
		// Check if there is a type, if so make sure it matches
		if actionMatch && matchType != "" {
			if question.Type == matchType {
				ask = true
			}
		} else if actionMatch {
			ask = true
		}
	}
	return ask
}
// Format the question
func AskQuestion(silent bool, reader *bufio.Reader, question *Question) error {

	if !silent {
		fmt.Println("-----------------------")
		text := question.Text

		if len(question.Accepted) > 0 {
			text = fmt.Sprintf(fmt.Sprintf("%v [%v]", text, strings.Join(question.Accepted, ",")))
		}
		fmt.Println("Q:", text)
		userAnswer := ""
		for userAnswer == "" {
			fmt.Printf("A(%v): ", question.Answer)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "" {
				if question.Answer != "" {
					input = question.Answer
				}
			}

			if validateAnswer(input, question) {
				userAnswer = input
			}

		}
		fmt.Println("Answer is:", userAnswer)
		question.Answer = userAnswer
	}
	return nil
}

func validateAnswer(input string, question *Question) bool {

	valid := false
	if len(question.Accepted) == 0 {
		valid = true
	} else {
		for _, a := range question.Accepted {
			if a == input {
				valid = true
				break
			}
		}
	}
	if !valid {
		fmt.Println("Invalid answer")
	}
	return valid
}

// SaveQuestions saves the interview for future use
func SaveQuestions(interview *Interview) error {

	fmt.Println("Saving interview file")
	savedInterview := Interview{
		Questions: interview.Questions,
	}
	save, err := json.Marshal(savedInterview)
	if err != nil {
		return err
	}
	fileName := InterviewDir + interview.LocDeploy.Name + InterviewFileType
	err = ioutil.WriteFile(fileName, save, 0644)
	if err != nil {
		return err
	}
	return nil
}
