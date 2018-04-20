package main

import (
"os/exec"
"fmt"
)

func main() {
	cmd := exec.Command("/bin/bash", "./exec.sh")
	bytes, err := cmd.Output()
	if err != nil {
		fmt.Printf("cmd.Output: %+v", err)
		return
	}
	fmt.Println(string(bytes))
}
