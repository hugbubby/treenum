package main

import (
	"fmt"
	"github.com/Ullaakut/nmap"
	"log"
	"net"
	"strings"
)

type Port struct {
	Protocol string
	Number   int
}

//Contains all the *possibilities* for the given port.
type PortStatus struct {
	Open     bool
	Filtered bool
	Closed   bool
}

type PortEnum struct {
	Port    Port
	Status  PortStatus
	Service string
}

func scan(target string, port Port) (*PortStatus, error) {
	switch port.Protocol {
	case "tcp":
		return ConnectScan(target, port.Number), nil
	case "udp":
		return UDPScan(target, port.Number)
	default:
		log.Println("Could not scan port along protocol", port.Protocol, ", lack the functionality.")
		s := new(PortStatus)
		s.Open = true
		s.Filtered = true
		s.Closed = true
		return s, nil
	}
}

//Really basic port scan. Assumes TCP for the port number.
func ConnectScan(target string, port int) *PortStatus {
	ret := new(PortStatus)
	address := net.JoinHostPort(target, fmt.Sprint(port))
	conn, err := net.Dial("tcp", address)
	if err == nil {
		conn.Close()
		//Probably open!
		ret.Open = true
		ret.Filtered = false
		ret.Closed = false
	} else {
		//Could be anything. I'm not bothered enough to do more analysis.
		ret.Open = false
		ret.Filtered = true
		ret.Closed = true
	}
	return ret
}

func UDPScan(target string, port int) (*PortStatus, error) {
	ret := new(PortStatus)
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithUDPScan(),
		nmap.WithPorts(fmt.Sprint(port)))
	if err == nil {
	} else {
		scanResult, err := scanner.Run()
		if err == nil {
			status := scanResult.Hosts[0].Ports[0].State.State
			if strings.Contains(status, "open") {
				ret.Open = true
			}
			if strings.Contains(status, "filtered") {
				ret.Filtered = true
			}
			if strings.Contains(status, "closed") {
				ret.Closed = true
			}
			if strings.Contains(status, "unknown") {
				ret.Open = true
				ret.Filtered = true
				ret.Closed = true
			}
		}
	}
	return ret, err
}
