package main

import (
	"GoLibApp/net/ssh_server"
	"fmt"
	"strings"
)

func main() {
	sshClient := ssh_server.Init("10.226.136.196", "root", "iaas-ops!@#")

	result, err := sshClient.Run("hping3 -S 10.226.136.111 -c 1 -p 80 | grep -v Process | grep -v HPING")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)

	flag := strings.Contains(result,"1 packets received")
	fmt.Println(flag)

	/*lines := strings.Split(strings.Trim(string(result), "\n"), "\n")
	fmt.Println(lines[1])
	fmt.Println(resolveResult(lines[1]))*/


}

func resolveResult(str string) string {
	strArray := strings.Split(str, ", ")
	end := strings.Index(strArray[2], "%")
	return strArray[2][:end]
}
