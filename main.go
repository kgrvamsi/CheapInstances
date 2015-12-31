package main

import (
	_ "flag"
	"fmt"
	"github.com/CheapInstances/server"
	"github.com/CheapInstances/slack"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/vaughan0/go-ini"
	"log"
	"os"
	_ "reflect"
	"strconv"
	"time"
)

var (
	choice     string
	timeNow    = time.Now()
	instanceid string
)

func main() {

	cfg, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Print(err.Error())
	}
	tokenPtr, _ := cfg.Get("slack", "token")
	channelPtr, _ := cfg.Get("slack", "channel")
	regionPtr, _ := cfg.Get("aws", "instance_region")
	instanceType, _ := cfg.Get("aws", "instance_type")
	zonePtr, _ := cfg.Get("aws", "instance_zone")
	spotPrice, _ := cfg.Get("aws", "spot_price")
	instanceImg, _ := cfg.Get("aws", "instance_image")
	instanceKey, _ := cfg.Get("aws", "instance_key")
	instanceCount, _ := cfg.Get("aws", "instance_count")

	for {
		msg :=
			` Enter any one in the following:
	1) Get the price history
	2) Create the Spot Instance
	3) Check the Spot Instance Request History
	4) Cancel the Spot Instance Request
	5) Exit`

		fmt.Println(msg)
		fmt.Scan(&choice)

		if choice == "1" {
			slack.AlertMessage(tokenPtr, channelPtr, "Getting the Price History")

			client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionPtr)})
			//	fmt.Println(reflect.TypeOf(client))

			params := &ec2.DescribeSpotPriceHistoryInput{
				InstanceTypes: []*string{
					aws.String(instanceType),
				},
				AvailabilityZone: aws.String(zonePtr),
				ProductDescriptions: []*string{
					// Linux/UNIX (Amazon VPC)
					aws.String("Linux/UNIX"),
				},
				MaxResults: aws.Int64(10),
			}

			//fmt.Println(reflect.TypeOf(params))
			//fmt.Println(params)
			server.SpotInstancePriceHistory(client, params)

		} else if choice == "2" {

			client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionPtr)})
			count, _ := strconv.ParseInt(instanceCount, 10, 64)
			param := &ec2.RequestSpotInstancesInput{
				SpotPrice:             aws.String(spotPrice),
				AvailabilityZoneGroup: aws.String(zonePtr),
				InstanceCount:         aws.Int64(count),
				LaunchSpecification: &ec2.RequestSpotLaunchSpecification{
					ImageId:      aws.String(instanceImg),
					InstanceType: aws.String(instanceType),
					KeyName:      aws.String(instanceKey),
				},
				Type: aws.String("one-time"),
				//		ValidFrom: aws.Time(time.Now()),
				//ValidUntil: aws.Time(time.Now()),
			}

			server.CreateSpotInstance(client, param)
			slack.AlertMessage(tokenPtr, channelPtr, "Requested Spot Instance Creation Initiated and Under Process for Evaluation")
		} else if choice == "3" {
			//slack.AlertMessage(tokenPtr, channelPtr, "Requested for List of Spot Instances Done by the Account")
			client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionPtr)})
			server.GetSpotInstancesReq(client)
		} else if choice == "4" {
			fmt.Println("Enter the Spot Instance Id")
			fmt.Scan(&instanceid)
			client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionPtr)})
			server.CancelSpotInstances(client, instanceid)
			//	slack.AlertMessage(tokenPtr, channelPtr, "Deleted the Spot Instance Request Id"+instanceid+")
		} else if choice == "5" {
			os.Exit(0)
		}
	}
}
