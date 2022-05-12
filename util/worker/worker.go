package worker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"log"
	"master/util"
	"os"
	"strings"
	"sync"
)

var lbd *lambda.Lambda
var Urls []string
var useS3 = false
var s3uploader *s3manager.Uploader
var s3c *s3.S3

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	config := &aws.Config{Region: aws.String("ap-northeast-1")}
	sess2 := session.Must(session.NewSession(config))
	lbd = lambda.New(sess, config)
	//s3c = s3.Add(sess, config)
	s3uploader = s3manager.NewUploader(sess2)
	s3c = s3.New(sess2)
}

func GetAllNames() ([]string, error) {
	input := &lambda.ListFunctionsInput{}
	functions, err := lbd.ListFunctions(input)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, v := range functions.Functions {
		names = append(names, *v.FunctionName)
	}
	return names, nil
}

func GetALlUrls() ([]string, error) {
	names, err := GetAllNames()
	if err != nil {
		return nil, err
	}
	var urls []string
	ch := make(chan string)
	for _, v := range names {
		go func(v string) {
			configs, err := lbd.ListFunctionUrlConfigs(&lambda.ListFunctionUrlConfigsInput{FunctionName: &v})
			if err != nil {
				log.Println(err)
				ch <- ""
				return
			}
			ch <- *configs.FunctionUrlConfigs[0].FunctionUrl
			//fmt.Printf("%+v\n", configs.FunctionUrlConfigs[0])
		}(v)
	}
	for range names {
		if s := <-ch; s != "" {
			urls = append(urls, s)
		}
	}
	return urls, nil
}

func CreateWorker(name string) error {
	Code := &lambda.FunctionCode{
		S3Bucket: aws.String("wangyu-code"),
		S3Key:    aws.String("function.zip"),
	}
	input := &lambda.CreateFunctionInput{
		Architectures: []*string{aws.String("x86_64")},
		Code:          Code,
		FunctionName:  &name,
		Handler:       aws.String("worker"),
		Role:          aws.String("arn:aws:iam::372623093426:role/service-role/worker-role-5er0369v"),
		Runtime:       aws.String("go1.x"),
		MemorySize:    aws.Int64(128),
	}
	_, err := lbd.CreateFunction(input)
	if err != nil {
		return err
	}
	NONE := "NONE"
	input2 := &lambda.CreateFunctionUrlConfigInput{
		AuthType:     &NONE,
		FunctionName: &name,
		Cors:         &lambda.Cors{},
	}
	_, err = lbd.CreateFunctionUrlConfig(input2)
	if err != nil {
		return err
	}
	Action := "lambda:InvokeFunctionUrl"
	StatementId := "FunctionURLAllowPublicAccess"
	star := "*"
	input3 := &lambda.AddPermissionInput{
		FunctionName:        &name,
		Action:              &Action,
		StatementId:         &StatementId,
		FunctionUrlAuthType: &NONE,
		//SourceAccount:       &star,
		Principal: &star,
	}
	_, err = lbd.AddPermission(input3)
	if err != nil {
		return err
	}
	return nil
}

func RemoveWorker(name string) error {
	urlInput := &lambda.DeleteFunctionUrlConfigInput{FunctionName: &name}
	_, err := lbd.DeleteFunctionUrlConfig(urlInput)
	if err != nil {
		return err
	}
	delInput := &lambda.DeleteFunctionInput{
		FunctionName: &name,
	}
	_, err = lbd.DeleteFunction(delInput)
	if err != nil {
		return err
	}
	return nil
}

func RemoveAllWorkers() error {
	names, err := GetAllNames()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	lm := util.NewGoLimit(5)
	for _, v := range names {
		wg.Add(1)
		go func(v string) {
			lm.Add()
			defer lm.Done()
			defer wg.Done()
			err := RemoveWorker(v)
			if err != nil {
				log.Println(err)
			}
		}(v)
	}
	wg.Wait()
	return nil
}

func GetFunctionUrl(v string) string {
	config, err := lbd.GetFunctionUrlConfig(&lambda.GetFunctionUrlConfigInput{FunctionName: aws.String(v)})
	if err != nil {
		log.Fatalln(err)
	}
	return *config.FunctionUrl
}

func PutCode(s, k string) {
	/*_, err := s3c.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String("wangyu-code"),
		Key:    aws.String(k),
	})
	if err != nil {
		log.Fatalln(err)
	}*/
	file, err := os.Open(s)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = s3uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("wangyu-code"),
		Key:    aws.String(k),
		Body:   file,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func Init() {
	names, err := GetAllNames()
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(names); i++ {
		if !strings.Contains(names[i], "worker") {
			names = append(names[:i], names[i+1:]...)
			i--
		}
	}
	if len(names) < 25 {
		PutCode("/tmp/function.zip", "function.zip")
		var wg sync.WaitGroup
		lm := util.NewGoLimit(5)
		for i := 0; i < 50-len(names); i++ {
			wg.Add(1)
			go func() {
				lm.Add()
				defer lm.Done()
				defer wg.Done()
				newUUID, err := uuid.NewUUID()
				if err != nil {
					log.Fatalln(err)
				}
				err = CreateWorker("worker-" + newUUID.String())
				if err != nil {
					log.Fatalln(err)
				}
			}()
		}
		wg.Wait()
	}
	Urls, err = GetALlUrls()
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range Urls {
		log.Println(v)
	}
}
