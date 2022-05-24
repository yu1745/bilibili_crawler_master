package worker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/yu1745/bilibili_crawler_master/util"
	"log"
	"os"
	"strings"
	"sync"
)

var lbd *lambda.Lambda
var Workers []Worker

type Worker struct {
	Name string
	Url  string
}

//var useS3 = false
var s3uploader *s3manager.Uploader

//var s3c *s3.S3

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	config := &aws.Config{Region: aws.String("ap-northeast-1")}
	sess2 := session.Must(session.NewSession(config))
	lbd = lambda.New(sess, config)
	//s3c = s3.Add(sess, config)
	s3uploader = s3manager.NewUploader(sess2)
	//s3c = s3.New(sess2)
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

func GetALlUrls() ([]Worker, error) {
	names, err := GetAllNames()
	if err != nil {
		return nil, err
	}
	var urls []Worker
	for _, v := range names {
		configs, err := lbd.ListFunctionUrlConfigs(&lambda.ListFunctionUrlConfigsInput{FunctionName: &v})
		if err != nil {
			log.Println(err)
		}
		urls = append(urls, Worker{
			Name: v,
			Url:  *configs.FunctionUrlConfigs[0].FunctionUrl,
		})
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
	/*NONE := "NONE"
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
	}*/
	return nil
}

func RemoveWorker(name string) error {
	/*urlInput := &lambda.DeleteFunctionUrlConfigInput{FunctionName: &name}
	_, err := lbd.DeleteFunctionUrlConfig(urlInput)
	if err != nil {
		return err
	}*/
	delInput := &lambda.DeleteFunctionInput{
		FunctionName: &name,
	}
	_, err := lbd.DeleteFunction(delInput)
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

func Invoke(name string, payload []byte) ([]byte, error) {
	output, err := lbd.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(name),
		Payload:      payload,
	})
	/*if output.FunctionError != nil {
		println(*output.  FunctionError)
	}*/
	return output.Payload, err
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

func Init(num int) {
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
	if len(names) < num {
		PutCode("/tmp/function.zip", "function.zip")
		var wg sync.WaitGroup
		lm := util.NewGoLimit(5)
		for i := 0; i < num-len(names); i++ {
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
	/*bytes, err := os.ReadFile("/tmp/names")
	if err != nil {
		bytes = make([]byte, 0)
	}
	_ = json.Unmarshal(bytes, &Workers)*/
	//if len(Workers) == 0 {
	/*Workers, err = GetALlUrls()
	if err != nil {
		log.Fatalln(err)
	}*/
	allNames, err := GetAllNames()
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range allNames {
		Workers = append(Workers, Worker{Name: v})
	}
	//marshal, _ := json.Marshal(&Workers)
	//_ = os.WriteFile("/tmp/names", marshal, 0644)
	//}
	for _, v := range Workers {
		log.Println(v.Name)
	}
}
