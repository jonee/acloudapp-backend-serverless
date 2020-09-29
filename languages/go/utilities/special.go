// Created by Jonee Ryan Ty
// Copyright ACloudApp

/**
 * Special Utility functions
 */

package utilities

import (
	"log"
	"strings"
	"time"

	"encoding/json"
)

func GetBearerToken(authHeader string) (token string) {
	if strings.HasPrefix(strings.ToUpper(authHeader), "BEARER") {
		token = authHeader[6:]
		token = strings.TrimSpace(token)
	}

	return token
}

func GetJsonStringBodyReturn(success bool, data map[string]interface{}) string {
	data["success"] = 0
	if success {
		data["success"] = 1
	}

	tmp, _ := json.Marshal(data)
	ret := string(tmp)

	log.Println(ret)

	return ret
}

func GetFormErrorsJsonStringBodyReturn(success bool, messageKey string, formErrors map[string]interface{}) string {
	d := make(map[string]interface{})
	d["success"] = 0
	if success {
		d["success"] = 1
	}

	if messageKey != "" {
		d["message_key"] = messageKey
	}

	d["form_errors"] = formErrors

	tmp, _ := json.Marshal(d)
	ret := string(tmp)

	log.Println(ret)

	return ret
}

func GetSimpleJsonStringBodyReturn(success bool, messageKey string, messageParameters map[string]string) string {
	d := make(map[string]interface{})
	d["success"] = 0
	if success {
		d["success"] = 1
	}

	if messageKey != "" {
		d["message_key"] = messageKey
	}

	if messageParameters != nil {
		d["message_parameters"] = messageParameters
	}

	tmp, _ := json.Marshal(d)
	ret := string(tmp)

	log.Println(ret)

	return ret
}

func ProcessTemplate(messageParameters map[string]string, emailText string, emailHTML string, subject string) (string, string, string) {
	for k, v := range messageParameters {
		emailText = strings.Replace(emailText, "%%"+k+"%%", v, -1)
		emailHTML = strings.Replace(emailHTML, "%%"+k+"%%", v, -1)
		subject = strings.Replace(subject, "%%"+k+"%%", v, -1)
	}

	emailText = strings.Replace(emailText, "%%subject%%", subject, -1)
	emailHTML = strings.Replace(emailHTML, "%%subject%%", subject, -1)

	return emailText, emailHTML, subject
}

func PrintMilestone(mapStore map[string]interface{}, s string) {
	t2 := time.Now()
	t1 := mapStore["last_milestone_time"].(time.Time)
	log.Println(s, (t2.UnixNano()-t1.UnixNano())/int64(time.Millisecond))

	mapStore["last_milestone_time"] = t2

	// PrintMilestone
}

/*
func GetS3File(s3File string) []byte {
	s3Creds := credentials.NewStaticCredentials(acaConfiguration.AWS_ACCESS_KEY_ID, acaConfiguration.AWS_SECRET_ACCESS_KEY, "")
	_, err := s3Creds.Get()
	if err != nil {
		log.Println("bad aws credentials")
	}

	s3Cfg := aws.NewConfig().WithRegion(acaConfiguration.AWS_REGION).WithCredentials(s3Creds)
	s3Downloader := s3manager.NewDownloader(session.New(s3Cfg))

	// extract key from s3 link
	tmpB := "s3://" + acaConfiguration.AWS_S3_BUCKET + "/"
	s3Key := s3File[len(tmpB):]
	log.Println("s3Key: ", s3Key)

	buff := &aws.WriteAtBuffer{}

	_, err = s3Downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(acaConfiguration.AWS_S3_BUCKET),
		Key:    aws.String(s3Key),
	})

	if err != nil {
		log.Println("Unable to download object %q from bucket %q, %v", s3Key, acaConfiguration.AWS_S3_BUCKET, err)
	}

	return buff.Bytes()

	// GetS3File
}
*/
