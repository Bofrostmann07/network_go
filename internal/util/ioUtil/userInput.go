package ioUtil

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func UserInputYesNo(text string, yesDefault bool) bool {
	fmt.Print(text)
	userInput := input()
	userInput = strings.ToLower(userInput)

	okayResponses := []string{"y", "yes"}
	nokayResponses := []string{"n", "no"}

	if yesDefault {
		okayResponses = append(okayResponses, "")
	} else {
		nokayResponses = append(nokayResponses, "")
	}

	if containsString(okayResponses, userInput) {
		return true
	} else if containsString(nokayResponses, userInput) {
		return false
	} else {
		return UserInputYesNo(text, yesDefault)
	}
}

func input() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func ReadLine() string {
	return input()
}
