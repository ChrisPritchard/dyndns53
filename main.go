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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var lastIPFilename = "./LATEST-IP-SET"

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	hostedZoneID := flag.String("hosted-zone-id", "", "route 53 hosted zone id (required)")
	targetDomain := flag.String("target-domain", "", "the domain to update (required)")
	overrideIPService := flag.String("ip-service", "", "external service to call to get the current ip address\nmust respond with just the ip address")
	overrideCurrentIP := flag.String("current-ip", "", "rather than querying the current ip using the ip service, just use this")
	beSilent := flag.Bool("quiet", false, "controls whether any messages are printed to out")
	flag.Parse()

	if *hostedZoneID == "" || *targetDomain == "" {
		flag.Usage()
		return
	}

	logStatus := func(message string) {
		if !*beSilent {
			log.Println(message)
		}
	}

	logStatus("targeting domain " + *targetDomain + " on hosted zone with id " + *hostedZoneID)

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
	newIP := *overrideCurrentIP
	if newIP == "" {
		logStatus("getting current ip address using " + myIPService)
		newIP, err = getCurrentIP(myIPService)
		if err != nil {
			log.Fatal(err)
		}
	} else if net.ParseIP(newIP) == nil {
		log.Fatal(errors.New("specified ip address is not valid"))
	}
	logStatus("current ip address is " + newIP)
	if lastIP == newIP {
		logStatus("ip address hasn't changed, exiting")
		return
	}
	logStatus("updating AWS...")
	err = updateAWS(*hostedZoneID, *targetDomain, newIP)
	if err != nil {
		log.Fatal(err)
	}
	logStatus("AWS updated successfully")
	ioutil.WriteFile(lastIPFilename, []byte(newIP), 0666)
	logStatus("saved new ip to status file\nfinished!")
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

func updateAWS(hostedZoneID, domain, newIPAddress string) error {
	session, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := route53.New(session)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(domain),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(newIPAddress),
							},
						},
						TTL:  aws.Int64(300),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String("Updated dynamic IP address"),
		},
		HostedZoneId: aws.String(hostedZoneID),
	}

	_, err = svc.ChangeResourceRecordSets(input)
	return err
}
