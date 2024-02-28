package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pdxfixit/awscollector"
	"github.com/pdxfixit/hostdb"
)

func main() {

	// load config
	loadConfig()

	var accounts []*awscollector.Account

	// gather general account information
	for _, ac := range config.AccountConfigs {
		err := ac.AccountInfo()
		if err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, ac.Account)
	}

	// gather details
	for _, ac := range config.AccountConfigs {

		err := ac.ListBuckets()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.ListCertificates()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeEcsClusters()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeRDSInstances()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeVirtualInterfaces()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeDynamoTables()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeLoadBalancers()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.ListHostedZones()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.ListPolicies()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.ListRoles()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.ListUsers()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeImages()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeKeyPairs()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.ListKeys()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeRepositories()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeInstances()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeSecurityGroups()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeSubnets()
		if err != nil {
			log.Fatal(err)
		}

		err = ac.DescribeVpcs()
		if err != nil {
			log.Fatal(err)
		}

	}

	// send to hostdb ... sample data output @ /sample-data
	for _, account := range accounts {
		for _, region := range account.RegionalInfo {

			var bucketsRecords,
				certificateRecords,
				clusterRecords,
				databaseRecords,
				dxRecords,
				dynamodbRecords,
				elbRecords,
				hostedZoneRecords,
				iamPolicyRecords,
				iamRoleRecords,
				iamUserRecords,
				imagesRecords,
				keypairRecords,
				keyRecords,
				repositoryRecords,
				reservationRecords,
				secgroupRecords,
				subnetRecords,
				volumeRecords,
				vpcRecords []hostdb.Record

			//
			// Buckets
			//
			for _, obj := range region.Buckets {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				bucketsRecords = append(bucketsRecords, hostdb.Record{
					Type:    "aws-s3-bucket",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err := postToHostdb(bucketsRecords, "aws-s3-bucket", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Certificates
			//
			for _, obj := range region.AcmCertificates {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				certificateRecords = append(certificateRecords, hostdb.Record{
					Type:    "aws-acm-certificate",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(certificateRecords, "aws-acm-certificate", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Clusters
			//
			for _, obj := range region.EcsClusters {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				clusterRecords = append(clusterRecords, hostdb.Record{
					Type:    "aws-ecs-cluster",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(clusterRecords, "aws-ecs-cluster", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Databases
			//
			for _, obj := range region.RdsInstances {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				databaseRecords = append(databaseRecords, hostdb.Record{
					Type:     "aws-rds-database",
					Hostname: fmt.Sprintf("%s:%d", *obj.Endpoint.Address, *obj.Endpoint.Port),
					Context:  map[string]interface{}{},
					Data:     jsonString,
				})
			}

			err = postToHostdb(databaseRecords, "aws-rds-database", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// DX
			//
			for _, obj := range region.DXVirtualInterfaces {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				dxRecords = append(dxRecords, hostdb.Record{
					Type:    "aws-directconnect",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(dxRecords, "aws-directconnect", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// DynamoDB
			//
			for _, obj := range region.DynamoDBTables {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				dynamodbRecords = append(dynamodbRecords, hostdb.Record{
					Type:    "aws-dynamodb-table",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(dynamodbRecords, "aws-dynamodb-table", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// ELBs
			//
			for _, obj := range region.Elbs {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				elbRecords = append(elbRecords, hostdb.Record{
					Type:    "aws-elb",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(elbRecords, "aws-elb", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Hosted Zones
			//
			for _, obj := range region.HostedZones {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				hostedZoneRecords = append(hostedZoneRecords, hostdb.Record{
					Type:    "aws-route53-hostedzone",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(hostedZoneRecords, "aws-route53-hostedzone", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// IAM Policies
			//
			for _, obj := range region.IamPolicies {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				iamPolicyRecords = append(iamPolicyRecords, hostdb.Record{
					Type:    "aws-iam-policy",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(iamPolicyRecords, "aws-iam-policy", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// IAM Roles
			//
			for _, obj := range region.IamRoles {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				iamRoleRecords = append(iamRoleRecords, hostdb.Record{
					Type:    "aws-iam-role",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(iamRoleRecords, "aws-iam-role", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// IAM Users
			//
			for _, obj := range region.IamUsers {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				iamUserRecords = append(iamUserRecords, hostdb.Record{
					Type:    "aws-iam-user",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(iamUserRecords, "aws-iam-user", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Images
			//
			for _, obj := range region.Images {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				imagesRecords = append(imagesRecords, hostdb.Record{
					Type:    "aws-ec2-image",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(imagesRecords, "aws-ec2-image", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Key Pairs
			//
			for _, obj := range region.KeyPairs {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				keypairRecords = append(keypairRecords, hostdb.Record{
					Type:    "aws-ec2-keypair",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(keypairRecords, "aws-ec2-keypair", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// KMS Keys
			//
			for _, obj := range region.Keys {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				keyRecords = append(keyRecords, hostdb.Record{
					Type:    "aws-kms-key",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(keyRecords, "aws-kms-key", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Repositories
			//
			for _, obj := range region.Repositories {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				repositoryRecords = append(repositoryRecords, hostdb.Record{
					Type:    "aws-ecr-repository",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(repositoryRecords, "aws-ecr-repository", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Instance Reservations
			//
			for _, obj := range region.Reservations {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				reservationRecords = append(reservationRecords, hostdb.Record{
					Type:    "aws-ec2-reservation",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(reservationRecords, "aws-ec2-reservation", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Security Groups
			//
			for _, obj := range region.SecurityGroups {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				secgroupRecords = append(secgroupRecords, hostdb.Record{
					Type:    "aws-ec2-securitygroup",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(secgroupRecords, "aws-ec2-securitygroup", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Subnets
			//
			for _, obj := range region.Subnets {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				subnetRecords = append(subnetRecords, hostdb.Record{
					Type:    "aws-ec2-subnet",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(subnetRecords, "aws-ec2-subnet", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// Volumes
			//
			for _, obj := range region.Volumes {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				volumeRecords = append(volumeRecords, hostdb.Record{
					Type:    "aws-ebs-volume",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(volumeRecords, "aws-ebs-volume", region, account)
			if err != nil {
				log.Fatal(err)
			}

			//
			// VPCs
			//
			for _, obj := range region.Vpcs {
				jsonString, err := json.Marshal(obj)
				if err != nil {
					log.Fatal(err)
				}

				vpcRecords = append(vpcRecords, hostdb.Record{
					Type:    "aws-ec2-vpc",
					Context: map[string]interface{}{},
					Data:    jsonString,
				})
			}

			err = postToHostdb(vpcRecords, "aws-ec2-vpc", region, account)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

}

func postToHostdb(records []hostdb.Record, recordType string, ri *awscollector.RegionalInfo, acct *awscollector.Account) (err error) {

	recordSet := hostdb.RecordSet{
		Type:      recordType,
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
		Context: map[string]interface{}{
			"aws-account-name":    acct.Name,
			"aws-account-id":      acct.ID,
			"aws-account-aliases": acct.Aliases,
			"aws-region":          ri.Region,
		},
		Committer: "hostdb-collector-aws",
		Records:   records,
	}

	if config.Collector.SampleData {
		if err := recordSet.Save(fmt.Sprintf("/sample-data/%s_%s_%s.json", acct.ID, ri.Region, recordType)); err != nil {
			return err
		}
	} else {
		if err := recordSet.Send(fmt.Sprintf("type=%s&aws-account-id=%s&aws-region=%s", recordType, acct.ID, ri.Region)); err != nil {
			return err
		}
	}

	return nil
}
