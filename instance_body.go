package main

type ServerType struct {
	Name      string `json:"name"`
	ImageRef  string `json:"imageRef"`
	FlavorRef string `json:"flavorRef"`
	Networks  []struct {
		Subnet string `json:"subnet"`
	} `json:"networks"`
	AvailabilityZone     string `json:"availability_zone"`
	KeyName              string `json:"key_name"`
	MaxCount             int    `json:"max_count"`
	MinCount             int    `json:"min_count"`
	BlockDeviceMappingV2 []struct {
		UUID                string `json:"uuid"`
		BootIndex           int    `json:"boot_index"`
		VolumeSize          int    `json:"volume_size"`
		DeviceName          string `json:"device_name"`
		SourceType          string `json:"source_type"`
		DestinationType     string `json:"destination_type"`
		DeleteOnTermination int    `json:"delete_on_termination"`
	} `json:"block_device_mapping_v2"`
	SecurityGroups []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
}

type ServerTypeTest struct {
	Name      string `json:"name"`
	KeyName   string `json:"key_name"`
	ImageRef  string `json:"imageRef"`
	FlavorRef string `json:"flavorRef"`
	Networks  []SubnetTest `json:"networks"`
	BlockDeviceMappingV2 []BlockDevice `json:"block_device_mapping_v2"`
}

type BlockDevice struct {
	UUID                string `json:"uuid"`
	BootIndex           int    `json:"boot_index"`
	VolumeSize          int    `json:"volume_size"`
	DeviceName          string `json:"device_name"`
	SourceType          string `json:"source_type"`
	DestinationType     string `json:"destination_type"`
	DeleteOnTermination int    `json:"delete_on_termination"`
}

type SubnetTest struct {
	Subnet string `json:"subnet"`
}

type Instance struct {
	Server ServerTypeTest `json:"server"`
}
