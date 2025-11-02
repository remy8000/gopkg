package tools

import (
	"fmt"
	"os"
)

// FileExistsAndNotEmpty checks if a file exists at the given path and is not empty.
// It returns true if the file exists and has a size greater than 0, otherwise false.
// If the file does not exist or there is an error accessing it, an appropriate error is returned.
func FileExistsAndNotEmpty(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("file does not exist")
		}
		return false, err // Other errors (e.g., permission issues)
	}
	
	// Check if it's a regular file and not empty
	if !info.IsDir() && info.Size() > 0 {
		return true, nil
	}
	return false, fmt.Errorf("file is empty")

}
