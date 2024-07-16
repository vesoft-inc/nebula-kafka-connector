package common

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type ResourceInfo struct {
	ResourceType        string
	OperationOnResource string
	ResourceList        []IPAndPort
	ClusterName         string
}

func genResourceRequest(resources ResourceInfo) (pc PromptContent) {
	pc.Label = fmt.Sprintf("Do you want to %s all the following %d %s listed in the config file: ", resources.OperationOnResource, len(resources.ResourceList), resources.ResourceType)
	for i, resource := range resources.ResourceList {
		if resource.Port != "" {
			pc.Label += fmt.Sprintf("%s:%s", resource.IP, resource.Port)
		} else {
			pc.Label += fmt.Sprintf("%s", resource.IP)
		}
		if i == len(resources.ResourceList)-1 {
			pc.Label += "?"
			break
		}
		if i < 2 {
			pc.Label += ", "
		} else {
			pc.Label += "...?"
			break
		}
	}
	return pc
}

type PromptContent struct {
	Label  string
	ErrMsg string
}

func PromptGetUserChoice(pc PromptContent) (answer string, err error) {
	validate := func(input string) error {
		if len(input) > 1 {
			return NgctlError("Invalid response.", "")
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Prompt:          "{{ . }} ",
		Valid:           "{{ . | green }} ",
		Invalid:         "{{ . | red }} ",
		Success:         "{{ . | bold }} ",
		ValidationError: "{{ . | red }} ",
	}
	prompt := promptui.Prompt{
		Label:     pc.Label + " (Y/N)",
		Templates: templates,
		Validate:  validate,
	}

	answer, err = prompt.Run()
	if err != nil {
		return "", NgctlError(pc.ErrMsg, "")
	}
	return answer, nil
}

func ConfirmResourceList(resrouces ResourceInfo) (confirmedResoruces ResourceInfo, err error) {
	pc := genResourceRequest(resrouces)
	answer, err := PromptGetUserChoice(pc)
	if err != nil {
		return ResourceInfo{}, NgctlError(pc.ErrMsg, "")
	}
	if strings.ToUpper(answer) == "Y" || strings.ToUpper(answer) == "" {
		// return all resources
		return resrouces, nil
	}
	if strings.ToUpper(answer) == "N" {
		confirmedResoruces.ResourceType = resrouces.ResourceType
		confirmedResoruces.OperationOnResource = resrouces.OperationOnResource
		confirmedResoruces.ClusterName = resrouces.ClusterName
		// not all are selected, select one by one
		for _, resource := range resrouces.ResourceList {
			if resource.Port != "" {
				pc.Label = "Do you want to " + resrouces.OperationOnResource + " " + resource.IP + ":" + resource.Port + "?"
			} else {
				pc.Label = "Do you want to " + resrouces.OperationOnResource + " " + resource.IP + "?"
			}
			answer, err = PromptGetUserChoice(pc)
			if err != nil {
				return ResourceInfo{}, NgctlError(pc.ErrMsg, "")
			}
			if strings.ToUpper(answer) == "Y" || strings.ToUpper(answer) == "" || strings.ToUpper(answer) == "T" {
				confirmedResoruces.ResourceList = append(confirmedResoruces.ResourceList, resource)
			}
			// else, continue
		}
		return confirmedResoruces, nil
	} else {
		// including the case with an input "Q"
		return ResourceInfo{}, nil
	}
}
