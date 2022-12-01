package connect

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/viper"
)

type AwsConnect struct {
	Sess *session.Session
}

func NewAwsConnect() *AwsConnect {
	return &AwsConnect{}
}

func (a *AwsConnect) Conn() error {
	costProfileName := viper.GetString("aws.credentials.cost-bill-name")
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", costProfileName),
		Region:      aws.String(endpoints.ApSoutheast1RegionID),
	})
	if err != nil {
		msg := fmt.Sprintf("Failed aws Conn session err: %s", err.Error())
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	_, err = sess.Config.Credentials.Get()
	if err != nil {
		msg := fmt.Sprintf("Failed aws Conn session get err: %s", err.Error())
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	a.Sess = sess
	return nil
}
