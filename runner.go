package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	subnetRange = flag.String("range", "10.0.0.0/8", "The subnet to scan")
	ports       = flag.String("ports", "135,139,22201,9595", "The ports to scan for")

	running = 0
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	var c chan string = make(chan string)

	go calc(c)
	runNmap(c)

	// var input string
	// fmt.Scanln(&input)
}

func calc(c chan string) {
	if *subnetRange != "10.0.0.0/8" {
		c <- *subnetRange
	} else {
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				subnet := fmt.Sprintf("10.%d.%d.0/24", i, j)
				c <- subnet
			}
		}
	}
}

func runNmap(c chan string) {
	for {
		if running < 5 {
			subnet := <-c

			go func(subnet string) {
				running += 1

				var fileName string

				subnetForFilename := strings.Replace(subnet, ".", "_", -1)
				subnetForFilename = strings.Replace(subnetForFilename, "/", "-", -1)
				portsForFilename := strings.Replace(*ports, ",", "_", -1)

				fileName = fmt.Sprintf("%d-%s_p%s.txt", time.Now().Unix(), subnetForFilename, portsForFilename)

				argv := []string{
					"-sS",
					"-O",
					subnet,
					"-p",
					*ports,
					"-oG",
					fileName,
				}

				log.Printf("Starting: nmap %s", argv)
				cmd := exec.Command("nmap", argv...)
				_, err := cmd.CombinedOutput()

				if err != nil {
					log.Fatal(err)
				}

				running -= 1
			}(subnet)
		}

		time.Sleep(time.Second * 1)
	}
}
