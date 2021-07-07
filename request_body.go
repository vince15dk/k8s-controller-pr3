package main

type InstanceInfo struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec struct {
		BlockSize int    `json:"blockSize"`
		FlavorRef string `json:"flavorRef"`
		ImageRef  string `json:"imageRef"`
		InstName  string `json:"instName"`
		KeyName   string `json:"keyName"`
		Password  string `json:"password"`
		SubnetID  string `json:"subnetId"`
		TenantID  string `json:"tenantId"`
		UserName  string `json:"userName"`
	} `json:"spec"`
}
