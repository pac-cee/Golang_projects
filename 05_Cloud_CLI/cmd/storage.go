package cmd

import (
	"fmt"

	"cloud-cli/internal/cloud"
	"cloud-cli/pkg/output"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(storageCmd)
	storageCmd.AddCommand(bucketListCmd)
	storageCmd.AddCommand(bucketCreateCmd)
	storageCmd.AddCommand(bucketDeleteCmd)
	storageCmd.AddCommand(objectListCmd)
	storageCmd.AddCommand(objectUploadCmd)
	storageCmd.AddCommand(objectDownloadCmd)
	storageCmd.AddCommand(objectDeleteCmd)

	// Bucket create flags
	bucketCreateCmd.Flags().StringP("name", "n", "", "bucket name")
	bucketCreateCmd.MarkFlagRequired("name")

	// Object upload flags
	objectUploadCmd.Flags().StringP("bucket", "b", "", "bucket name")
	objectUploadCmd.Flags().StringP("file", "f", "", "file to upload")
	objectUploadCmd.MarkFlagRequired("bucket")
	objectUploadCmd.MarkFlagRequired("file")

	// Object download flags
	objectDownloadCmd.Flags().StringP("bucket", "b", "", "bucket name")
	objectDownloadCmd.Flags().StringP("output", "o", "", "output file")
	objectDownloadCmd.MarkFlagRequired("bucket")
}

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Manage cloud storage",
	Long:  `Create and manage storage buckets and objects across different cloud providers.`,
}

var bucketListCmd = &cobra.Command{
	Use:   "list-buckets",
	Short: "List storage buckets",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		format, _ := cmd.Flags().GetString("output")

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		buckets, err := client.ListBuckets()
		if err != nil {
			fmt.Printf("Error listing buckets: %v\n", err)
			return
		}

		output.Print(buckets, format)
	},
}

var bucketCreateCmd = &cobra.Command{
	Use:   "create-bucket",
	Short: "Create a new storage bucket",
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		name, _ := cmd.Flags().GetString("name")

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		bucket, err := client.CreateBucket(name)
		if err != nil {
			fmt.Printf("Error creating bucket: %v\n", err)
			return
		}

		fmt.Printf("Created bucket: %s\n", bucket.Name)
	},
}

var bucketDeleteCmd = &cobra.Command{
	Use:   "delete-bucket [bucket-name]",
	Short: "Delete a storage bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		bucketName := args[0]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := client.DeleteBucket(bucketName); err != nil {
			fmt.Printf("Error deleting bucket: %v\n", err)
			return
		}

		fmt.Printf("Deleted bucket: %s\n", bucketName)
	},
}

var objectListCmd = &cobra.Command{
	Use:   "list-objects [bucket-name]",
	Short: "List objects in a bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		format, _ := cmd.Flags().GetString("output")
		bucketName := args[0]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		objects, err := client.ListObjects(bucketName)
		if err != nil {
			fmt.Printf("Error listing objects: %v\n", err)
			return
		}

		output.Print(objects, format)
	},
}

var objectUploadCmd = &cobra.Command{
	Use:   "upload [object-name]",
	Short: "Upload an object to a bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		bucket, _ := cmd.Flags().GetString("bucket")
		file, _ := cmd.Flags().GetString("file")
		objectName := args[0]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := client.UploadObject(bucket, objectName, file); err != nil {
			fmt.Printf("Error uploading object: %v\n", err)
			return
		}

		fmt.Printf("Uploaded %s to %s/%s\n", file, bucket, objectName)
	},
}

var objectDownloadCmd = &cobra.Command{
	Use:   "download [object-name]",
	Short: "Download an object from a bucket",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		bucket, _ := cmd.Flags().GetString("bucket")
		output, _ := cmd.Flags().GetString("output")
		objectName := args[0]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := client.DownloadObject(bucket, objectName, output); err != nil {
			fmt.Printf("Error downloading object: %v\n", err)
			return
		}

		fmt.Printf("Downloaded %s/%s to %s\n", bucket, objectName, output)
	},
}

var objectDeleteCmd = &cobra.Command{
	Use:   "delete-object [bucket-name] [object-name]",
	Short: "Delete an object from a bucket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		provider, _ := cmd.Flags().GetString("provider")
		region, _ := cmd.Flags().GetString("region")
		bucketName := args[0]
		objectName := args[1]

		client, err := cloud.NewClient(provider, region)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := client.DeleteObject(bucketName, objectName); err != nil {
			fmt.Printf("Error deleting object: %v\n", err)
			return
		}

		fmt.Printf("Deleted object: %s/%s\n", bucketName, objectName)
	},
}
