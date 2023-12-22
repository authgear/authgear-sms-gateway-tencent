package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type CustomSMSGatewayPayload struct {
	Body string `json:"body,omitempty"`
	To   string `json:"to,omitempty"`
}

var codeRegexp = regexp.MustCompile(`\d{6}`)

func extractCode(body string) string {
	code := codeRegexp.FindString(body)
	if code == "" {
		return body
	}
	return code
}

func main() {
	godotenv.Load()

	TENCENTCLOUD_SECRET_ID := os.Getenv("TENCENTCLOUD_SECRET_ID")
	TENCENTCLOUD_SECRET_KEY := os.Getenv("TENCENTCLOUD_SECRET_KEY")
	TENCENTCLOUD_REGION := os.Getenv("TENCENTCLOUD_REGION")
	TENCENTCLOUD_SMS_SDK_APP_ID := os.Getenv("TENCENTCLOUD_SMS_SDK_APP_ID")
	TENCENTCLOUD_SMS_TEMPLATE_ID := os.Getenv("TENCENTCLOUD_SMS_TEMPLATE_ID")

	credential := common.NewCredential(
		TENCENTCLOUD_SECRET_ID,
		TENCENTCLOUD_SECRET_KEY,
	)

	cpf := profile.NewClientProfile()
	client, _ := sms.NewClient(credential, TENCENTCLOUD_REGION, cpf)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload CustomSMSGatewayPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Printf("failed to decode hook payload: %v\n", err)
			return
		}

		request := sms.NewSendSmsRequest()
		request.SmsSdkAppId = common.StringPtr(TENCENTCLOUD_SMS_SDK_APP_ID)
		request.TemplateId = common.StringPtr(TENCENTCLOUD_SMS_TEMPLATE_ID)
		request.PhoneNumberSet = common.StringPtrs([]string{payload.To})
		code := extractCode(payload.Body)
		request.TemplateParamSet = common.StringPtrs([]string{code})

		response, err := client.SendSms(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			log.Printf("tencent cloud sdk error: %v\n", err)
			return
		}
		if err != nil {
			log.Printf("other error: %v\n", err)
			return
		}

		b, err := json.Marshal(response.Response)
		if err != nil {
			log.Printf("failed to marshal response: %v\n", err)
			return
		}

		log.Printf("%v\n", string(b))
	})

	log.Printf("listening on :7000\n")
	http.ListenAndServe(":7000", nil)
}
