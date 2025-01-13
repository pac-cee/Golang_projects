package cloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type awsClient struct {
	ec2Client *ec2.Client
	s3Client  *s3.Client
	region    string
}

func newAWSClient(region string) (Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	return &awsClient{
		ec2Client: ec2.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
		region:    region,
	}, nil
}

// VM operations
func (c *awsClient) ListVMs() ([]VM, error) {
	result, err := c.ec2Client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	var vms []VM
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			vm := VM{
				ID:    *instance.InstanceId,
				State: string(instance.State.Name),
			}
			if instance.PublicIpAddress != nil {
				vm.PublicIP = *instance.PublicIpAddress
			}
			vms = append(vms, vm)
		}
	}

	return vms, nil
}

func (c *awsClient) CreateVM(config VMConfig) (*VM, error) {
	// Implementation for AWS EC2 instance creation
	// This is a simplified version. In a real implementation, you would:
	// 1. Validate the instance type
	// 2. Look up the AMI ID
	// 3. Configure networking
	// 4. Set up security groups
	// 5. Add tags
	return nil, fmt.Errorf("not implemented")
}

func (c *awsClient) DeleteVM(id string) error {
	_, err := c.ec2Client.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{id},
	})
	return err
}

func (c *awsClient) StartVM(id string) error {
	_, err := c.ec2Client.StartInstances(context.TODO(), &ec2.StartInstancesInput{
		InstanceIds: []string{id},
	})
	return err
}

func (c *awsClient) StopVM(id string) error {
	_, err := c.ec2Client.StopInstances(context.TODO(), &ec2.StopInstancesInput{
		InstanceIds: []string{id},
	})
	return err
}

// Storage operations
func (c *awsClient) ListBuckets() ([]Bucket, error) {
	result, err := c.s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	for _, b := range result.Buckets {
		buckets = append(buckets, Bucket{
			Name:      *b.Name,
			CreatedAt: b.CreationDate.Format(time.RFC3339),
		})
	}

	return buckets, nil
}

func (c *awsClient) CreateBucket(name string) (*Bucket, error) {
	_, err := c.s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: &name,
	})
	if err != nil {
		return nil, err
	}

	return &Bucket{
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
		Region:    c.region,
	}, nil
}

func (c *awsClient) DeleteBucket(name string) error {
	_, err := c.s3Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: &name,
	})
	return err
}

func (c *awsClient) ListObjects(bucket string) ([]Object, error) {
	result, err := c.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}

	var objects []Object
	for _, obj := range result.Contents {
		objects = append(objects, Object{
			Name:         *obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified.Format(time.RFC3339),
			ETag:         *obj.ETag,
		})
	}

	return objects, nil
}

func (c *awsClient) UploadObject(bucket, object, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = c.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &object,
		Body:   f,
	})
	return err
}

func (c *awsClient) DownloadObject(bucket, object, file string) error {
	result, err := c.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &object,
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, result.Body)
	return err
}

func (c *awsClient) DeleteObject(bucket, object string) error {
	_, err := c.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &object,
	})
	return err
}
