package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	_ "reflect"
	"strconv"
	"time"

	"github.com/CheapInstances/server"
	"github.com/CheapInstances/slack"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
	"github.com/vaughan0/go-ini"
)

var (
	choice     string
	timeNow    = time.Now()
	instanceid string
	reqData    = [][]string{}
)

func main() {

	cfg, err := ini.LoadFile("config.ini")
	if err != nil {
		log.Print(err.Error())
	}

	flag.Usage = func() {
		help := `Usage: cheapinstances -count="3" -instance="m3.medium" -price="0.10" -region="us-east-1" -zone="us-east-1c"`
		fmt.Println(help)
		flag.PrintDefaults()
	}

	// Config File Paramters Parsed Variables
	tokenPtr, _ := cfg.Get("slack", "token")
	channelPtr, _ := cfg.Get("slack", "channel")
	regionType, _ := cfg.Get("aws", "instance_region")
	instanceType, _ := cfg.Get("aws", "instance_type")
	zoneType, _ := cfg.Get("aws", "instance_zone")
	spotPrice, _ := cfg.Get("aws", "spot_price")
	instanceImg, _ := cfg.Get("aws", "instance_image")
	instanceKey, _ := cfg.Get("aws", "instance_key")
	instanceCount, _ := cfg.Get("aws", "instance_count")

	// Flag Parsed Variables

	instancePtr := flag.String("instance", "", "This Represents the Instance Type")
	regionPtr := flag.String("region", "", "This Represents the Region to Use")
	zonePtr := flag.String("zone", "", "This Represents the zone to use for the respective Region")
	pricePtr := flag.String("price", "", "This Represents the Price to use for the Requested instance")
	instCountPtr := flag.String("count", "", "This is used for the number of instance requests you need")

	flag.Parse()

	if *instancePtr == "" && *regionPtr == "" && *zonePtr == "" && *pricePtr == "" && *instCountPtr == "" {
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
				// slack.AlertMessage(tokenPtr, channelPtr, "Getting the Price History")

				client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionType)})
				//	fmt.Println(reflect.TypeOf(client))

				// Why are we doing this?
				// server.GetAvailableZones(client)

				// server.GetTheLeastZone(instanceType, client)
				params := &ec2.DescribeSpotPriceHistoryInput{
					InstanceTypes: []*string{
						aws.String(instanceType),
					},
					AvailabilityZone: aws.String(zoneType),
					ProductDescriptions: []*string{
						// Linux/UNIX (Amazon VPC)
						aws.String("Linux/UNIX"),
					},
					MaxResults: aws.Int64(10),
				}

				_, data := server.SpotInstancePriceHistory(client, params)

				fmt.Println(data.SpotPriceHistory)

			} else if choice == "2" {

				reader := bufio.NewReader(os.Stdin)
				fmt.Println("Enter the UserData:")
				userData, _ := reader.ReadString('\n')
				fmt.Println(userData + "Is added to the Instance")
				client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionType)})
				count, _ := strconv.ParseInt(instanceCount, 10, 64)
				param := &ec2.RequestSpotInstancesInput{
					SpotPrice:             aws.String(spotPrice),
					AvailabilityZoneGroup: aws.String(zoneType),
					InstanceCount:         aws.Int64(count),
					LaunchSpecification: &ec2.RequestSpotLaunchSpecification{
						ImageId:      aws.String(instanceImg),
						InstanceType: aws.String(instanceType),
						KeyName:      aws.String(instanceKey),
						UserData:     aws.String(userData),
					},
					Type: aws.String("one-time"),
					//		ValidFrom: aws.Time(time.Now()),
					//ValidUntil: aws.Time(time.Now()),
				}

				server.CreateSpotInstance(client, param)
				slack.AlertMessage(tokenPtr, channelPtr, "Requested Spot Instance Creation Initiated and Under Process for Evaluation")
			} else if choice == "3" {
				//slack.AlertMessage(tokenPtr, channelPtr, "Requested for List of Spot Instances Done by the Account")

				client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionType)})
				_, reqInst := server.GetSpotInstancesReq(client)
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"SpotId", "State", "Status", "Type"})
				for i, _ := range reqInst.SpotInstanceRequests {
					data := []string{*(reqInst.SpotInstanceRequests[i].SpotInstanceRequestId), *(reqInst.SpotInstanceRequests[i].State), *(reqInst.SpotInstanceRequests[i].Status.Code), *(reqInst.SpotInstanceRequests[i].Type)}
					table.Append(data)
				}
				table.Render()

			} else if choice == "4" {
				fmt.Println("Enter the Spot Instance Id")
				fmt.Scan(&instanceid)
				client := ec2.New(session.New(), &aws.Config{Region: aws.String(regionType)})
				server.CancelSpotInstances(client, instanceid)
				//	slack.AlertMessage(tokenPtr, channelPtr, "Deleted the Spot Instance Request Id"+instanceid+")
			} else if choice == "5" {
				os.Exit(0)
			}
		}
	} else {

		fmt.Println(*instancePtr, *zonePtr, *regionPtr, *pricePtr)
	}
}
