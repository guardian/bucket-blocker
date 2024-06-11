package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

func blockPublicAccess(s3Client *s3.S3, name string) (*s3.PutPublicAccessBlockOutput, error) {
	resp, err := s3Client.PutPublicAccessBlock(&s3.PutPublicAccessBlockInput{
		Bucket: aws.String(name),
		PublicAccessBlockConfiguration: &s3.PublicAccessBlockConfiguration{
			BlockPublicAcls:       aws.Bool(true),
			IgnorePublicAcls:      aws.Bool(true),
			BlockPublicPolicy:     aws.Bool(true),
			RestrictPublicBuckets: aws.Bool(true),
		},
	})
	if err != nil {
		return resp, err
	}
	fmt.Println("Public access blocked for bucket: " + name)
	return resp, nil
}

func validateCredentials(stsClient *sts.STS, profile string) (*sts.GetCallerIdentityOutput, error) {
	resp, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return resp, errors.New("Could not find valid credentials for profile: " + profile)
	}
	return resp, nil
}

func bucketBlocksPublicAccess(s3Client *s3.S3, bucketName string) bool {
	s3BucketPolicy, err := s3Client.GetPublicAccessBlock(&s3.GetPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Println("No public access ACL found for bucket: " + bucketName)
		return false
	}
	return *s3BucketPolicy.PublicAccessBlockConfiguration.BlockPublicAcls
}

func bucketInRegion(s3Client *s3.S3, bucketName string, region string) bool {
	location, err := s3Client.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		fmt.Println("Error getting bucket location for bucket: " + bucketName + " Error: " + err.Error())
		return false
	}
	if location.LocationConstraint == nil {
		return region == "us-east-1"
	}
	return *location.LocationConstraint == region
}

func main() {
	excludedBuckets := flag.String("buckets", "", "A comma separated list of bucket names to exclude from public access blocking")
	profile := flag.String("profile", "", "The name of the profile to use")
	region := flag.String("region", "", "The region of the bucket")
	flag.Parse()

	if *profile == "" {
		fmt.Println("Please provide a profile name")
		return
	}

	if *region == "" {
		fmt.Println("Please provide a region")
		return
	}

	parsedBuckets := splitAndTrim(*excludedBuckets)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           *profile,
		Config: aws.Config{
			Region: aws.String(*region),
		},
	}))

	stsClient := sts.New(sess)
	_, err := validateCredentials(stsClient, *profile)
	if err != nil {
		fmt.Println(err)
		return
	}

	s3Client := s3.New(sess)

	buckets, err := s3Client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Error listing buckets: " + err.Error())
		return
	}
	fmt.Println("Found " + fmt.Sprintf("%d", len(buckets.Buckets)) + " buckets in " + *profile)

	var openRegionalBuckets s3.ListBucketsOutput
	for _, bucket := range buckets.Buckets {
		if bucketInRegion(s3Client, *bucket.Name, *region) && !bucketBlocksPublicAccess(s3Client, *bucket.Name) {
			openRegionalBuckets.Buckets = append(openRegionalBuckets.Buckets, bucket)
		}
	}
	fmt.Println("Found " + fmt.Sprintf("%d", len(openRegionalBuckets.Buckets)) + " buckets in " + *region + " that do not have a public access ACL")

	for _, bucket := range openRegionalBuckets.Buckets[:10] {
		if !contains(parsedBuckets, *bucket.Name) {
			fmt.Println("Blocking public access for bucket: " + *bucket.Name)
			// _, err = blockPublicAccess(s3Client, *bucket.Name)
			// if err != nil {
			// 	fmt.Println("Error blocking public access: " + err.Error())
			// 	return
			// }
		} else {
			fmt.Println("Skipping bucket: " + *bucket.Name)

		}
	}
}
