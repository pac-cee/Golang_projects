package cloud

import "fmt"

type azureClient struct {
	region string
}

func newAzureClient(region string) (Client, error) {
	return &azureClient{region: region}, nil
}

// VM operations
func (c *azureClient) ListVMs() ([]VM, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *azureClient) CreateVM(config VMConfig) (*VM, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *azureClient) DeleteVM(id string) error {
	return fmt.Errorf("not implemented")
}

func (c *azureClient) StartVM(id string) error {
	return fmt.Errorf("not implemented")
}

func (c *azureClient) StopVM(id string) error {
	return fmt.Errorf("not implemented")
}

// Storage operations
func (c *azureClient) ListBuckets() ([]Bucket, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *azureClient) CreateBucket(name string) (*Bucket, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *azureClient) DeleteBucket(name string) error {
	return fmt.Errorf("not implemented")
}

func (c *azureClient) ListObjects(bucket string) ([]Object, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *azureClient) UploadObject(bucket, object, file string) error {
	return fmt.Errorf("not implemented")
}

func (c *azureClient) DownloadObject(bucket, object, file string) error {
	return fmt.Errorf("not implemented")
}

func (c *azureClient) DeleteObject(bucket, object string) error {
	return fmt.Errorf("not implemented")
}
