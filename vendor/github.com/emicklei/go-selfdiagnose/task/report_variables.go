package task

import (
	"fmt"
	"sort"

	selfdiagnose "github.com/emicklei/go-selfdiagnose"
)

// Reports the variables and values the system is using
type ReportVariables struct {
	Description string
	VariableMap map[string]interface{}
}

func (r ReportVariables) Run(ctx *selfdiagnose.Context, result *selfdiagnose.Result) {
	result.Passed = true
	// creates a summary of all variables and values ordered alphabetically
	summary := ""
	sortedNameList := make([]string, 0, len(r.VariableMap))
	for name := range r.VariableMap {
		sortedNameList = append(sortedNameList, name)
	}
	sort.Strings(sortedNameList)
	for _, name := range sortedNameList {
		summary += fmt.Sprintf("%s = %v<br/>", name, r.VariableMap[name])
	}
	result.Reason = summary
}

func (r ReportVariables) Comment() string {
	if len(r.Description) != 0 {
		return r.Description
	}
	return "Configuration"
}
