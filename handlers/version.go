package handlers

import "fmt"

func Version(currentVersion string) error {
	fmt.Println("BRUCE version: " + currentVersion)
	tag, err := getLatestTag("brucehq", "bruce")
	if err != nil {
		return err
	}
	fmt.Println("Latest version: " + tag)
	return nil
}
