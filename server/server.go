package server

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	_ "reflect"
	"strconv"
)

var (
	region string
)

type RegionOutput struct {
	Regions []struct {
		Endpoint   string `json:"Endpoint"`
		RegionName string `json:"RegionName"`
	} `json:"Regions"`
}

type ZonesOutput struct {
	AvailabilityZones []struct {
		RegionName string `json:"RegionName"`
		State      string `json:"State"`
		ZoneName   string `json:"ZoneName"`
	} `json:"AvailabilityZones"`
}

type Data struct {
	Region string   `json:"region,omitempty"`
	Zones  []string `json:"zones"`
}

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

//func GetTheLeastZone(instance string, client *ec2.EC2) {
func GetTheLeastZone(instance string, client *ec2.EC2) {
	_, k := GetAvailableZones(client)
	fmt.Println(len(k))

	for i := range k {
		var val []float64
		region := k[i].Region
		fmt.Println(region)
		svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
		fmt.Println(len(k[i].Zones))
		for j := range k[i].Zones {

			params := &ec2.DescribeSpotPriceHistoryInput{
				InstanceTypes: []*string{
					aws.String(instance),
				},
				AvailabilityZone:    aws.String(k[i].Zones[j]),
				ProductDescriptions: []*string{aws.String("Linux/UNIX")},
				MaxResults:          aws.Int64(1),
			}
			resp, _ := svc.DescribeSpotPriceHistory(params)
			if len(resp.SpotPriceHistory) == 0 {
				fmt.Println("No Resources Available")
			} else {
				value := *(resp.SpotPriceHistory[0].SpotPrice)
				converted, _ := strconv.ParseFloat(value, 32)
				val = append(val, converted)
			}
		}
	}
}

/*
* This Function will get all the Zones available for the User
* From Different Regions Available for his account
 */
func GetAvailableZones(client *ec2.EC2) (err error, datas []Data) {

	params := &ec2.DescribeRegionsInput{}

	resp, err := client.DescribeRegions(params)
	if err != nil {
		log.Println(err.Error())
	}

	out, err := json.MarshalIndent(&resp, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(reflect.TypeOf(string(out)))
	var r *RegionOutput

	er := json.Unmarshal([]byte(out), &r)
	if er != nil {
		fmt.Println(er)
	}
	datas = make([]Data, len(r.Regions))

	for i := 0; i < len(r.Regions); i++ {
		region := r.Regions[i].RegionName
		svc := ec2.New(session.New(), &aws.Config{Region: aws.String(region)})
		param := &ec2.DescribeAvailabilityZonesInput{}
		k, err := svc.DescribeAvailabilityZones(param)
		if err != nil {
			fmt.Print(err.Error())
		}
		znes, err := json.MarshalIndent(&k, "", " ")
		if err != nil {
			fmt.Println(err)
		}
		var z *ZonesOutput
		error := json.Unmarshal([]byte(znes), &z)
		if error != nil {
			fmt.Println(error)
		}
		q := []string{}
		for _, availZones := range z.AvailabilityZones {
			q = append(q, availZones.ZoneName)
		}
		datas[i] = Data{region, q}
		//datas = append(datas, Data{region, q})
	}
	return
}
