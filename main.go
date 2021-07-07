package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	admissionv1 "k8s.io/api/admission/v1"
	"log"
	"net/http"
)

var (
	//router *gin.Engine
	urlCreateRepo = "https://api-identity.infrastructure.cloud.toast.com/v2.0/tokens"
	instance = &Instance{}
)

//func init(){
//	router = gin.Default()
//}
//
//func mapUrls(){
//	router.POST("")
//}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers
	client := http.Client{}
	return client.Do(request)
}

func CreateInstancePost(url string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers
	client := http.Client{}
	return client.Do(request)
}

func handleCall(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle Called")
	input := &admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(input)
	if err != nil {
		sendErr(w, fmt.Errorf("could not unmarshal review: %v", err))
		return
	}
	inst := &InstanceInfo{}
	err = json.Unmarshal(input.Request.Object.Raw, inst)
	if err != nil {
		sendErr(w, fmt.Errorf("could not unmarshal pod: %v", err))
	}
	createInstance(inst)
}

func createInstance(inst *InstanceInfo) {
	fmt.Println("getting token")
	// get tokenId
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	b := &CreateAccessRequest{Auth: Tenant{
		TenantId: inst.Spec.TenantID,
		PasswordCredentials: UserInfo{
			UserName: inst.Spec.UserName,
			Password: inst.Spec.Password,
		},
	}}
	response, _ := Post(urlCreateRepo, b, headers)
	bytes, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	var result CreateAccessResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal create repo successful response: %s", err.Error()))
	}

	fmt.Println("creating instance")
	// create instance
	newHeader := http.Header{}
	newHeader.Set("Content-Type", "application/json")
	newHeader.Set("X-Auth-Token", result.Access.Token.ID)
	urlCreateInstance := "https://kr1-api-instance.infrastructure.cloud.toast.com/v2/" + inst.Spec.TenantID + "/servers"
	instance.Server.Name = inst.Spec.InstName
	instance.Server.ImageRef = inst.Spec.ImageRef
	instance.Server.FlavorRef = inst.Spec.FlavorRef
	instance.Server.Networks = []SubnetTest{{inst.Spec.SubnetID}}
	instance.Server.KeyName = inst.Spec.KeyName
	instance.Server.BlockDeviceMappingV2 = []BlockDevice{{UUID: inst.Spec.ImageRef, BootIndex: 0,
		VolumeSize: inst.Spec.BlockSize, DeviceName: "vda", SourceType: "image", DestinationType: "volume", DeleteOnTermination: 1}}

	newResponse, err := CreateInstancePost(urlCreateInstance, instance, newHeader)
	if err != nil {
		fmt.Println(err)
	}
	newBytes, err := ioutil.ReadAll(newResponse.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer newResponse.Body.Close()
	fmt.Println(string(newBytes))
}

func sendErr(w http.ResponseWriter, err error) {
	out, err := json.Marshal(map[string]string{
		"Err": err.Error(),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(out)
}

func main() {
	fmt.Println("server started")
	mux := http.NewServeMux()
	mux.HandleFunc("/call", handleCall)
	srv := &http.Server{Addr: ":443", Handler: mux}
	log.Fatal(srv.ListenAndServeTLS("/certs/webhook.crt", "/certs/webhook-key.pem"))
}
