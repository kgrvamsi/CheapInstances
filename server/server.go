package server

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

func SpotInstancePriceHistory(client *ec2.EC2, params *ec2.DescribeSpotPriceHistoryInput) {

	resp, err := client.DescribeSpotPriceHistory(params)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func CreateSpotInstance(client *ec2.EC2, params *ec2.RequestSpotInstancesInput) {

	resp, err := client.RequestSpotInstances(params)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func GetSpotInstancesReq(client *ec2.EC2) {

	params := &ec2.DescribeSpotInstanceRequestsInput{}

	resp, err := client.DescribeSpotInstanceRequests(params)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func CancelSpotInstances(client *ec2.EC2, instanceid string) {
	params := &ec2.CancelSpotInstanceRequestsInput{
		SpotInstanceRequestIds: []*string{ // Required
			aws.String(instanceid), // Required
		},
	}

	resp, err := client.CancelSpotInstanceRequests(params)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(resp)
}

func GetavailableZones(client *ec2.EC2) {
}
