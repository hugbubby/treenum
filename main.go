package main

import (
	"context"
	"fmt"
	"github.com/Ullaakut/nmap"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, _ := getConfig()
	if len(os.Args) > 1 {
		config.Target = os.Args[1]
	}

	go startEnum(config)

	log.Println("Autoenum is now running, CTRL-C or SIGINT to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func startEnum(config Config) {
	openPorts := make(chan Port, 100)
	go scanPorts(config.Target, openPorts)
	go enumPorts(config, openPorts)
}

func scanPorts(target string, openPorts chan<- Port) {
	protocols := [2]string{"tcp", "udp"}
	for _, k := range protocols {
		for i := 1; i < 65536; i++ {
			p := Port{k, i}
			go func(port Port) {
				couldBe, err := scan(target, port)
				if err == nil {
					if couldBe.Open && !couldBe.Filtered && !couldBe.Closed {
						openPorts <- port
					}
				} else {
					log.Println("Error scanning port ", port)
					log.Println(err)
				}
			}(p)
			time.Sleep(time.Microsecond * 10)
		}
	}
}

func enumPorts(config Config, openPorts <-chan Port) {
	baseDir := config.GetOutputDirectory() + "/" + config.Target
	for {
		port := <-openPorts
		log.Println("Detected open port", port.Protocol, "/", port.Number)
		go enumPort(config, baseDir+"/"+port.Protocol+"/"+fmt.Sprint(port.Number), port)
	}
}

func enumPort(config Config, outputDir string, port Port) {
	os.MkdirAll(outputDir, os.ModeDir|OS_USER_RWX|OS_ALL_R)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var scanner *nmap.Scanner
	var err error
	if port.Protocol == "udp" {
		scanner, err = nmap.NewScanner(nmap.WithTargets(config.Target),
			nmap.WithPorts(fmt.Sprint(port.Number)),
			nmap.WithContext(ctx),
			nmap.WithServiceInfo(),
			nmap.WithUDPScan())
	} else { //Default is TCP. We don't support other protocols for now.
		scanner, err = nmap.NewScanner(nmap.WithTargets(config.Target),
			nmap.WithPorts(fmt.Sprint(port.Number)),
			nmap.WithContext(ctx),
			nmap.WithServiceInfo())
	}
	if err == nil {
		result, err := scanner.Run()
		if err == nil {
			enumDir1 := os.Getenv("HOME") + "/.config/treenum/services/" + "all"
			host := result.Hosts[0]
			portMap := host.Ports[0]
			enumDir2 := os.Getenv("HOME") + "/.config/treenum/services/" + portMap.Service.Name
			dirs := []string{enumDir1, enumDir2}
			for _, enumDir := range dirs {
				files, err := ioutil.ReadDir(enumDir)
				if err == nil {
					for _, file := range files {
						if !file.IsDir() &&
							(file.Mode()&0111 != 0) /*File is globally executable*/ {
							cmd := exec.Command(enumDir+"/"+file.Name(), config.Target, fmt.Sprint(port.Number), fmt.Sprint(port.Protocol))
							cmd.Dir = outputDir + "/"
							go cmd.Run()
						}
					}
				}
			}
		} else {
			log.Println("Could not run nmap scan:")
			log.Println(err)
		}
	} else {
		log.Println("Could not build nmap scanner (this should not happen):")
		log.Println(err)
	}
}
