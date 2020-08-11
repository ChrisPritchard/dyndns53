package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var lastIPFilename = "./LATEST-IP-SET"

func main() {
	log.SetFlags(0)

	overrideIPService := flag.String("myip-service", "", "external service that when called will return the external IP address and nothing else")
	beSilent := flag.Bool("quiet", false, "controls whether any messages are printed to out")
	flag.Parse()

	logStatus := func(message string) {
		if !*beSilent {
			log.Println(message)
		}
	}

	myIPService := "https://ipinfo.io/ip"
	if overrideIPService != nil && *overrideIPService != "" {
		myIPService = *overrideIPService
	}

	lastIP, err := getLastIP()
	if err != nil {
		log.Fatal(err)
	}
	if lastIP == "" {
		logStatus("no last ip found - this must be the first run")
	} else {
		logStatus("last set ip address was " + lastIP)
	}
	logStatus("getting current ip address using " + myIPService)
	newIP, err := getCurrentIP(myIPService)
	if err != nil {
		log.Fatal(err)
	}
	logStatus("current ip address is " + newIP)
	if lastIP == newIP {
		logStatus("ip address hasn't changed, exiting")
	} else {
		ioutil.WriteFile(lastIPFilename, []byte(newIP), 0666)
	}
}

func getLastIP() (string, error) {
	if _, err := os.Stat(lastIPFilename); os.IsNotExist(err) {
		return "", nil
	}
	result, err := ioutil.ReadFile(lastIPFilename)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func getCurrentIP(myIPService string) (string, error) {
	resp, err := http.Get(myIPService)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	text := strings.Trim(string(body), "\r\n\t ")
	ip := net.ParseIP(text)
	if ip == nil {
		return "", errors.New("response from " + myIPService + " was not a valid ip address")
	}
	return text, nil
}
