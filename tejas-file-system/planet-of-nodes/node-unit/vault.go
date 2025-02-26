package nodeunit

import (
	"fmt"
	"os/exec"
)

func isExists(valutName string) bool {
	cmd := exec.Command("docker", "volume", "inspect", valutName)
	err := cmd.Run()
	return err == nil // If there is no error the volume exist
}

func CreateStore(valutName string) error {

	if isExists(valutName) {
		return nil
	}

	cmd := exec.Command("docker", "volume", "create", valutName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating volume: %w, output: %s", err, string(output))
	}
	fmt.Println("volume created", string(output))
	return nil
}
