package connect

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsConnect struct {
	Sess *session.Session
}

func NewAwsConnect() *AwsConnect {
	return &AwsConnect{}
}

func (a *AwsConnect) Conn(profileName string, regionName string) error {
	if len(regionName) == 0 {
		regionName = endpoints.ApSoutheast1RegionID
	}
	if len(profileName) == 0 {
		profileName = "default"
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", profileName),
		Region:      aws.String(regionName),
	})
	if err != nil {
		msg := fmt.Sprintf("Failed aws Conn session err (profileName: %s, regionName: %s): %s", profileName, regionName, err.Error())
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		msg := fmt.Sprintf("Failed aws Conn session get Credentials err (profileName: %s, regionName: %s): %s", profileName, regionName, err.Error())
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	a.Sess = sess
	return nil
}
