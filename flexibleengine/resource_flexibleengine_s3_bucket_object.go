package flexibleengine

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func resourceS3BucketObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceS3BucketObjectPut,
		Read:   resourceS3BucketObjectRead,
		Update: resourceS3BucketObjectPut,
		Delete: resourceS3BucketObjectDelete,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"acl": {
				Type:         schema.TypeString,
				Default:      "private",
				Optional:     true,
				ValidateFunc: validateS3BucketObjectAclType,
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_language": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
			},

			"server_side_encryption": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateS3BucketObjectServerSideEncryption,
				Computed:     true,
			},

			"etag": {
				Type: schema.TypeString,
				// This will conflict with SSE-C and SSE-KMS encryption and multi-part upload
				// if/when it's actually implemented. The Etag then won't match raw-file MD5.
				// See http://docs.aws.amazon.com/AmazonS3/latest/API/RESTCommonResponseHeaders.html
				Optional: true,
				Computed: true,
			},

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"website_redirect": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceS3BucketObjectPut(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	s3conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine s3 client: %s", err)
	}

	var body io.ReadSeeker

	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error opening S3 bucket object source (%s): %s", source, err)
		}

		body = file
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	} else {
		return fmt.Errorf("Must specify \"source\" or \"content\" field")
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    aws.String(d.Get("acl").(string)),
		Body:   body,
	}

	if v, ok := d.GetOk("cache_control"); ok {
		putInput.CacheControl = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_type"); ok {
		putInput.ContentType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		putInput.ContentEncoding = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_language"); ok {
		putInput.ContentLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		putInput.ContentDisposition = aws.String(v.(string))
	}

	if v, ok := d.GetOk("server_side_encryption"); ok {
		putInput.ServerSideEncryption = aws.String(v.(string))
	}

	if v, ok := d.GetOk("website_redirect"); ok {
		putInput.WebsiteRedirectLocation = aws.String(v.(string))
	}

	resp, err := s3conn.PutObject(putInput)
	if err != nil {
		return fmt.Errorf("Error putting object in S3 bucket (%s): %s", bucket, err)
	}

	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	d.Set("etag", strings.Trim(*resp.ETag, `"`))

	d.Set("version_id", resp.VersionId)
	d.SetId(key)
	return resourceS3BucketObjectRead(d, meta)
}

func resourceS3BucketObjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	s3conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine s3 client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	resp, err := s3conn.HeadObject(
		&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		// If S3 returns a 404 Request Failure, mark the object as destroyed
		if awsErr, ok := err.(awserr.RequestFailure); ok && awsErr.StatusCode() == 404 {
			d.SetId("")
			log.Printf("[WARN] Error Reading Object (%s), object not found (HTTP status 404)", key)
			return nil
		}
		return err
	}
	log.Printf("[DEBUG] Reading S3 Bucket Object meta: %s", resp)

	d.Set("cache_control", resp.CacheControl)
	d.Set("content_disposition", resp.ContentDisposition)
	d.Set("content_encoding", resp.ContentEncoding)
	d.Set("content_language", resp.ContentLanguage)
	d.Set("content_type", resp.ContentType)
	d.Set("version_id", resp.VersionId)
	d.Set("server_side_encryption", resp.ServerSideEncryption)
	d.Set("website_redirect", resp.WebsiteRedirectLocation)

	d.Set("etag", strings.Trim(*resp.ETag, `"`))

	return nil
}

func resourceS3BucketObjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	s3conn, err := config.computeS3conn(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating FlexibleEngine s3 client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	if _, ok := d.GetOk("version_id"); ok {
		// Bucket is versioned, we need to delete all versions
		vInput := s3.ListObjectVersionsInput{
			Bucket: aws.String(bucket),
			Prefix: aws.String(key),
		}
		out, err := s3conn.ListObjectVersions(&vInput)
		if err != nil {
			return fmt.Errorf("Failed listing S3 object versions: %s", err)
		}

		for _, v := range out.Versions {
			input := s3.DeleteObjectInput{
				Bucket:    aws.String(bucket),
				Key:       aws.String(key),
				VersionId: v.VersionId,
			}
			_, err := s3conn.DeleteObject(&input)
			if err != nil {
				return fmt.Errorf("Error deleting S3 object version of %s:\n %s:\n %s",
					key, v, err)
			}
		}
	} else {
		// Just delete the object
		input := s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}
		_, err := s3conn.DeleteObject(&input)
		if err != nil {
			return fmt.Errorf("Error deleting S3 bucket object: %s  Bucket: %q Object: %q", err, bucket, key)
		}
	}

	return nil
}

func validateS3BucketObjectAclType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	cannedAcls := map[string]bool{
		s3.ObjectCannedACLPrivate:                true,
		s3.ObjectCannedACLPublicRead:             true,
		s3.ObjectCannedACLPublicReadWrite:        true,
		s3.ObjectCannedACLAuthenticatedRead:      true,
		s3.ObjectCannedACLAwsExecRead:            true,
		s3.ObjectCannedACLBucketOwnerRead:        true,
		s3.ObjectCannedACLBucketOwnerFullControl: true,
	}

	sentenceJoin := func(m map[string]bool) string {
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, fmt.Sprintf("%q", k))
		}
		sort.Strings(keys)

		length := len(keys)
		words := make([]string, length)
		copy(words, keys)

		words[length-1] = fmt.Sprintf("or %s", words[length-1])
		return strings.Join(words, ", ")
	}

	if _, ok := cannedAcls[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid canned ACL type %q. Valid types are either %s",
			k, value, sentenceJoin(cannedAcls)))
	}
	return
}

func validateS3BucketObjectStorageClassType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	storageClass := map[string]bool{
		s3.StorageClassStandard:          true,
		s3.StorageClassReducedRedundancy: true,
		s3.StorageClassStandardIa:        true,
	}

	if _, ok := storageClass[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Storage Class type %q. Valid types are either %q, %q, or %q",
			k, value, s3.StorageClassStandard, s3.StorageClassReducedRedundancy,
			s3.StorageClassStandardIa))
	}
	return
}

func validateS3BucketObjectServerSideEncryption(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	serverSideEncryption := map[string]bool{
		s3.ServerSideEncryptionAes256: true,
		s3.ServerSideEncryptionAwsKms: true,
	}

	if _, ok := serverSideEncryption[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Server Side Encryption value %q. Valid values are %q and %q",
			k, value, s3.ServerSideEncryptionAes256, s3.ServerSideEncryptionAwsKms))
	}
	return
}
