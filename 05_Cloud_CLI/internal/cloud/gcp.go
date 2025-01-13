package cloud

import "fmt"

type gcpClient struct {
	region string
}

func newGCPClient(region string) (Client, error) {
	return &gcpClient{region: region}, nil
}

// VM operations
func (c *gcpClient) ListVMs() ([]VM, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *gcpClient) CreateVM(config VMConfig) (*VM, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *gcpClient) DeleteVM(id string) error {
	return fmt.Errorf("not implemented")
}

func (c *gcpClient) StartVM(id string) error {
	return fmt.Errorf("not implemented")
}

func (c *gcpClient) StopVM(id string) error {
	return fmt.Errorf("not implemented")
}

// Storage operations
func (c *gcpClient) ListBuckets() ([]Bucket, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *gcpClient) CreateBucket(name string) (*Bucket, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *gcpClient) DeleteBucket(name string) error {
	return fmt.Errorf("not implemented")
}

func (c *gcpClient) ListObjects(bucket string) ([]Object, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *gcpClient) UploadObject(bucket, object, file string) error {
	return fmt.Errorf("not implemented")
}

func (c *gcpClient) DownloadObject(bucket, object, file string) error {
	return fmt.Errorf("not implemented")
}

func (c *gcpClient) DeleteObject(bucket, object string) error {
	return fmt.Errorf("not implemented")
}
