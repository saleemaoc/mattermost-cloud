package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	log "github.com/sirupsen/logrus"
)

// AWS interface for use by other packages.
type AWS interface {
	GetCertificateSummaryByTag(key, value string) (*acm.CertificateSummary, error)

	GetAndClaimVpcResources(clusterID, owner string, logger log.FieldLogger) (ClusterResources, error)
	ReleaseVpc(clusterID string, logger log.FieldLogger) error

	GetPrivateZoneDomainName(logger log.FieldLogger) (string, error)
	CreatePrivateCNAME(dnsName string, dnsEndpoints []string, logger log.FieldLogger) error
	CreatePublicCNAME(dnsName string, dnsEndpoints []string, logger log.FieldLogger) error

	DeletePrivateCNAME(dnsName string, logger log.FieldLogger) error
	DeletePublicCNAME(dnsName string, logger log.FieldLogger) error

	TagResource(resourceID, key, value string, logger log.FieldLogger) error
	UntagResource(resourceID, key, value string, logger log.FieldLogger) error
	IsValidAMI(AMIImage string) (bool, error)
}

// Client is a client for interacting with AWS resources.
type Client struct {
	ACM            acmiface.ACMAPI
	EC2            ec2iface.EC2API
	IAM            iamiface.IAMAPI
	RDS            rdsiface.RDSAPI
	S3             s3iface.S3API
	Route53        route53iface.Route53API
	SecretsManager secretsmanageriface.SecretsManagerAPI

	AvailabilityZones []*string
}

// NewClient returns a new AWS client.
func NewClient(sess *session.Session, availabilityZones ...string) *Client {
	return &Client{
		ACM:            acm.New(sess),
		EC2:            ec2.New(sess),
		IAM:            iam.New(sess),
		RDS:            rds.New(sess),
		S3:             s3.New(sess),
		Route53:        route53.New(sess),
		SecretsManager: secretsmanager.New(sess),

		AvailabilityZones: PointerToStringArray(availabilityZones, aws.String),
	}
}
