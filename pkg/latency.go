package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/dmonteroh/distributed-resource-collector/internal"
	"github.com/gin-gonic/gin"
)

func LatencyEndpoint(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	execMode := c.MustGet("EXEC_MODE").(string)
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	latencyTargets, err := internal.LatencyJsonToStrcut(string(jsonData))
	if err != nil {
		panic(err)
	}

	latencyResults := internal.LatencyResults{
		Hostname: latencyTargets.Hostname,
		Results:  []internal.LatencyResult{},
	}

	for _, target := range latencyTargets.Targets {
		c1 := make(chan internal.LatencyResult)
		go func(target internal.LatencyTarget) {
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			//cmd := "touch latency_" + latencyTargets.Hostname + " && echo '" + timestamp[len(timestamp)-1:] + "' > latency_" + latencyTargets.Hostname + " && cat latency_" + latencyTargets.Hostname
			cmd := "echo " + timestamp[len(timestamp)-1:]
			funcStart := time.Now()
			elapsed := int64(0)
			result, ok := sshServer(target, cmd, timestamp[len(timestamp)-1:])
			if ok {
				elapsed = time.Since(funcStart).Milliseconds()
				if execMode == "DEBUG" {
					fmt.Println(result)
				}

			} else {
				elapsed = int64(-1)
			}
			latencyResult := internal.LatencyResult{Hostname: target.Hostname, Latency: elapsed}
			if execMode == "DEBUG" {
				fmt.Println(latencyResult.String())
			}
			c1 <- latencyResult
		}(target)
		latencyResults.Results = append(latencyResults.Results, <-c1)
	}

	if execMode == "DEBUG" {
		c.JSON(200, latencyResults)
	} else {
		c.JSON(200, latencyResults)
	}
}

func sshServer(target internal.LatencyTarget, cmd string, expected string) (string, bool) {
	config := &ssh.ClientConfig{
		User: target.HostUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(target.HostPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", target.Hostname, target.Hostport), config)
	if err != nil {
		return err.Error(), false
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err.Error(), false
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err.Error(), false
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return err.Error(), false
	}
	defer stdin.Close()

	err = session.Start(cmd)
	if err != nil {
		return fmt.Sprintf("unable to execute remote command: %s", err), false
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, stdout); err != nil {
		return fmt.Sprintf("reading failed: %s", err), false
	}

	if sttyOutput := buf.String(); !strings.Contains(sttyOutput, expected) {
		return fmt.Sprintf("FALSE RESULT, expected %s and got %s", expected, sttyOutput), false
	} else {
		return buf.String(), true
	}

}
