package pkg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/dmonteroh/distributed-resource-collector/internal"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/shomali11/parallelizer"
)

func LatencyEndpoint(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	execMode := c.MustGet("EXEC_MODE").(string)
	targetsApp := c.MustGet("TARGETS_APP").(string)
	latencyTargets := latencyTargetsHandler(targetsApp)

	latencyHandler(execMode, latencyTargets)
}

func ManualLatencyEndpoint(c *gin.Context) {
	defer internal.RecoverEndpoint(c)
	execMode := c.MustGet("EXEC_MODE").(string)
	jsonData, _ := ioutil.ReadAll(c.Request.Body)
	latencyTargets, err := internal.LatencyTargetsJsonToStruct(string(jsonData))
	if err != nil {
		panic(err)
	}
	latencyResults := latencyHandler(execMode, latencyTargets)
	c.JSON(200, latencyResults)
}

func latencyTargetsHandler(url string) internal.LatencyTargets {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		panic(err)
	}
	defer res.Body.Close()
	jsonData, _ := ioutil.ReadAll(res.Body)
	latencyTargets, err := internal.LatencyTargetsJsonToStruct(string(jsonData))
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		panic(err)
	}
	return latencyTargets
}

func latencyHandler(execMode string, latencyTargets internal.LatencyTargets) internal.LatencyResults {
	latencyGroup := parallelizer.NewGroup(parallelizer.WithPoolSize(12), parallelizer.WithJobQueueSize(3))
	defer latencyGroup.Close()
	latencyResults := internal.LatencyResults{
		Source:  latencyTargets.Source,
		Results: []internal.LatencyResult{},
	}

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(len(latencyTargets.Targets))

	c1 := make(chan internal.LatencyResult)

	go func() {
		for c := range c1 {
			latencyResults.Results = append(latencyResults.Results, c)
		}
	}()

	for _, target := range latencyTargets.Targets {
		go func(target internal.LatencyTarget) {
			handleLatencyTarget(target, execMode, c1, waitGroup)
		}(target)
	}
	// Timestamp after operations
	waitGroup.Wait()
	//Without this sleep, the last result is skipped, don't know why
	time.Sleep(time.Microsecond * 15)
	tmpTime := time.Now()
	latencyResults.Timestamp = internal.LatencyTimestamp{
		TimeLocal:   tmpTime,
		TimeSeconds: tmpTime.Unix(),
		TimeNano:    tmpTime.UnixNano(),
	}

	return latencyResults
}

func handleLatencyTarget(target internal.LatencyTarget, execMode string, c1 chan internal.LatencyResult, waitGroup *sync.WaitGroup) {
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
	defer waitGroup.Done()
}

func sshServer(target internal.LatencyTarget, cmd string, expected string) (string, bool) {
	config := &ssh.ClientConfig{
		User: target.HostUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(target.HostPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoRSA,
			ssh.KeyAlgoDSA,
			ssh.KeyAlgoECDSA256,
			ssh.KeyAlgoECDSA384,
			ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		// optional tcp connect timeout
		Timeout: 5 * time.Second,
	}

	client, err := ssh.Dial("tcp", target.Hostname+":"+target.Hostport, config)
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

	err = session.Run(cmd)
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

func sendLatency(targetUrl string, latencyUrl string, execMode string) {
	defer recoverHeartbeat()
	latencyTargets := latencyTargetsHandler(targetUrl)
	latencyResults := latencyHandler(execMode, latencyTargets)
	if execMode == "DEBUG" {
		fmt.Println("DEUBG MODE - POST")
		fmt.Println(latencyResults.String())
	}
	res, err := http.Post(latencyUrl, "application/json", bytes.NewBuffer([]byte(latencyResults.String())))
	if err != nil {
		fmt.Println(res)
		fmt.Println(err)
		panic(err)
	}

	defer res.Body.Close()
	fmt.Println(res.Body)
}

func LatencyCron(cron *gocron.Scheduler, seconds int, targetUrl string, latencyUrl string, execMode string) {

	defer recoverCron()
	cronRes, cronErr := cron.Every(seconds).Seconds().Do(sendLatency, targetUrl, latencyUrl, execMode)
	if cronErr != nil {
		panic(cronErr)
	}
	go func() {
		if execMode == "DEBUG" {
			_, cronErrDebug := cron.Every(seconds).Seconds().Do(debugCron, cronRes)
			if cronErrDebug != nil {
				panic(cronErr)
			}
		}
	}()
}