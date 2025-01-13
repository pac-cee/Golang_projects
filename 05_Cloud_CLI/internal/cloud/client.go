package cloud

import (
	"fmt"
)

// Client represents a cloud provider client
type Client interface {
	// VM operations
	ListVMs() ([]VM, error)
	CreateVM(config VMConfig) (*VM, error)
	DeleteVM(id string) error
	StartVM(id string) error
	StopVM(id string) error

	// Storage operations
	ListBuckets() ([]Bucket, error)
	CreateBucket(name string) (*Bucket, error)
	DeleteBucket(name string) error
	ListObjects(bucket string) ([]Object, error)
	UploadObject(bucket, object, file string) error
	DownloadObject(bucket, object, file string) error
	DeleteObject(bucket, object string) error
}

// NewClient creates a new cloud provider client
func NewClient(provider, region string) (Client, error) {
	switch provider {
	case "aws":
		return newAWSClient(region)
	case "gcp":
		return newGCPClient(region)
	case "azure":
		return newAzureClient(region)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// VM represents a virtual machine instance
type VM struct {
	ID       string
	Name     string
	Type     string
	State    string
	PublicIP string
	Region   string
}

// VMConfig represents configuration for creating a VM
type VMConfig struct {
	Name  string
	Type  string
	Image string
}

// Bucket represents a storage bucket
type Bucket struct {
	Name         string
	CreatedAt    string
	Region       string
	ObjectCount  int64
	TotalSize    int64
}

// Object represents a storage object
type Object struct {
	Name         string
	Size         int64
	LastModified string
	ContentType  string
	ETag         string
}
