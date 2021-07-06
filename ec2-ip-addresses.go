package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func main() {

	tagPtr := flag.String("tag", "", "ec2-instance-tag-value")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
				},
			},
		},
	}

	if *tagPtr != "" {
		input = &ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name: aws.String("tag-value"),
					Values: []string{
						*tagPtr,
					},
				},
				{
					Name: aws.String("instance-state-name"),
					Values: []string{
						"running",
					},
				},
			},
		}
	}

	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return
	}

	for _, r := range result.Reservations {
		fmt.Println("Reservation ID: " + *r.ReservationId)
		fmt.Println("Instance IDs:")
		for _, i := range r.Instances {
			fmt.Println("   " + *i.InstanceId)
			if i.NetworkInterfaces[len(i.NetworkInterfaces)-1].Association != nil {
				if i.NetworkInterfaces[len(i.NetworkInterfaces)-1].Association.PublicDnsName != nil {
					fmt.Println("   " + *i.NetworkInterfaces[len(i.NetworkInterfaces)-1].Association.PublicDnsName)
				}
				if i.NetworkInterfaces[len(i.NetworkInterfaces)-1].Association.PublicIp != nil {
					fmt.Println("   " + *i.NetworkInterfaces[len(i.NetworkInterfaces)-1].Association.PublicIp)
				}
			} else {
				fmt.Println("   " + "offline")
			}
		}

		fmt.Println("")
	}
}
