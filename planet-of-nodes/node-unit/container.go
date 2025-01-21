package nodeunit

import (
	"fmt"
	"os/exec"
	"strconv"
)

func SpinUpContainer(nodeId string, vaultname string, capacity int64, port string) error {

	if err := CreateStore(vaultname); err != nil {
		return fmt.Errorf("error of creating data valult: %w", err)
	}

	cmd := exec.Command("docker", "run", "-d",
		"--name", nodeId,
		"--env", fmt.Sprintf("NODE_ID=%s", nodeId),
		"--env", fmt.Sprintf("STORAGE_CAPACITY=%s", strconv.FormatInt(capacity, 10)),
		"-v", fmt.Sprintf("%s:/data", vaultname),
		"-p", fmt.Sprintf("%s:%s", port, "8080"),
		"node-image")

	fmt.Println("Executing command:", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Output:", string(output))
		return fmt.Errorf("failed to run docker: %v", err)
	}
	fmt.Println("Container started:", string(output))
	return nil
}
