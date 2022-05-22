package ioUtil

import (
	"errors"
	"fmt"
	"os"
)

func ExistsDir(path string, create bool) (exists bool, error error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if create {
				err := createDir(path)
				if err != nil {
					return false, fmt.Errorf("%s", err)
				}
				//log.Printf("Created directory @%q", path)
				return true, nil
			}
			//log.Printf("Directory @%q does not exist", path)
			return false, nil
		} else {
			//log.Printf("Error accessing directory @%q: %s", path, err)
			return false, fmt.Errorf("can not access directory - %s", err)
		}
	}
	if !info.IsDir() {
		//log.Printf("Error: Directory @%q is a file", path)
		return false, errors.New("directory is file")
	}
	//log.Printf("Found directory @%q", path)
	return true, nil
}

func createDir(path string) (error error) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		//log.Printf("Error: Could not create directory @%q - %s", path, err)
		return fmt.Errorf("could not create directory - %s", err)
	}
	return nil
}
