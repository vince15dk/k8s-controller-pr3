package main

type ServerInfo struct {
	Servers []struct {
		Image struct {
			ID    string `json:"id"`
		} `json:"image"`
		Flavor             struct {
			ID    string `json:"id"`
		} `json:"flavor"`
		ID             string `json:"id"`
		Name                             string        `json:"name"`
		TenantID                         string        `json:"tenant_id"`
	} `json:"servers"`
}
