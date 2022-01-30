package utils

import (
	"errors"
	"fmt"
	"strings"
)

type Testing struct {
	testCount          int
	successCount       int
	logPrefix          string
	failureMessageList []string
	failureTestList    []string
}

func CreateTesting(logPrefix string) Testing {
	return Testing{
		testCount:          0,
		successCount:       0,
		logPrefix:          "[" + logPrefix + "] ",
		failureMessageList: make([]string, 0),
		failureTestList:    make([]string, 0),
	}
}

func (t *Testing) TestIntEqual(expected int, actual int, successMessage string, failMessage string, fatal bool) error {
	if expected == actual {
		t.resultSuccess(successMessage)
		return nil
	} else {
		t.failureTestList = append(t.failureTestList, " Expected: "+fmt.Sprint(expected)+" Actual: "+fmt.Sprint(actual))
		return t.resultFailure(failMessage, fatal)
	}
}

func (t *Testing) TestStringEqual(expected string, actual string, successMessage string, failMessage string, fatal bool) error {
	if expected == actual {
		t.resultSuccess(successMessage)
		return nil
	} else {
		t.failureTestList = append(t.failureTestList, " Expected: "+fmt.Sprint(strings.ReplaceAll(expected, "\n", " "))+"\n Actual  : "+fmt.Sprint(strings.ReplaceAll(actual, "\n", " ")))
		return t.resultFailure(failMessage, fatal)
	}
}

func (t *Testing) TestBool(expected bool, actual bool, successMessage string, failMessage string, fatal bool) error {
	if expected == actual {
		t.resultSuccess(successMessage)
		return nil
	} else {
		t.failureTestList = append(t.failureTestList, " Expected: "+fmt.Sprint(expected)+" Actual: "+fmt.Sprint(actual))
		return t.resultFailure(failMessage, fatal)
	}
}

func (t *Testing) TestError(err error, successMessage string, failMessage string, fatal bool) error {
	if err == nil {
		t.resultSuccess(successMessage)
		return nil
	} else {
		t.failureTestList = append(t.failureTestList, " Not nil error : "+err.Error())
		return t.resultFailure(failMessage, fatal)
	}
}

func (t *Testing) resultSuccess(successMessage string) {
	t.testCount++
	t.successCount++
	LogSuccess(t.logPrefix + successMessage)
}

func (t *Testing) resultFailure(failMessage string, fatal bool) error {
	t.testCount++
	err := errors.New(t.logPrefix + failMessage + " :")
	if fatal {
		HandlePanicError(err)
	} else {
		t.failureMessageList = append(t.failureMessageList, failMessage)
		LogError(err.Error())
		tmp := strings.Split(t.failureTestList[len(t.failureTestList)-1], "\n")
		for _, tx := range tmp {
			LogError(t.logPrefix + "\t->" + tx)
		}
	}
	return nil
}

func (t *Testing) Conclusion() error {
	if t.successCount == t.testCount {
		LogSuccess(t.logPrefix + "[CONCLUSION] All tests succeeded (" + t.successRatio() + ") All seems good !")
		return nil
	} else {
		LogError(t.logPrefix + "[CONCLUSION] " + t.successRatio() + " successful tests. List of incorrect results :")
		for i, failMsg := range t.failureMessageList {
			prefix := t.logPrefix + "Error " + fmt.Sprint(i+1) + "/" + fmt.Sprint(len(t.failureMessageList)) + " : "
			LogError(prefix + failMsg)
			//LogError(t.logPrefix + "\t->" + t.failureTestList[i])
		}
		return errors.New(t.failureRatio())
	}
}

func (t *Testing) successRatio() string {
	return fmt.Sprint(t.successCount) + "/" + fmt.Sprint(t.testCount)
}

func (t *Testing) failureRatio() string {
	return fmt.Sprint(t.testCount-t.successCount) + "/" + fmt.Sprint(t.testCount)
}
