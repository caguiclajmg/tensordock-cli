package api

type Server struct {
	Cost struct {
		Charged    float32 `json:"charged"`
		HourOff    float32 `json:"hour_off"`
		HourOn     float32 `json:"hour_on"`
		MinutesOff float32 `json:"minutes_off"`
		MinutesOn  float32 `json:"minutes_on"`
	} `json:"cost"`
	CPUModel     string                       `json:"cpu_model"`
	GPUCount     int                          `json:"gpu_count"`
	GPUModel     string                       `json:"gpu_model"`
	Id           string                       `json:"id"`
	Ip           string                       `json:"ip"`
	Links        map[string]map[string]string `json:"links"`
	Location     string                       `json:"location"`
	Name         string                       `json:"name"`
	Ram          int                          `json:"ram"`
	Status       string                       `json:"status"`
	Storage      int                          `json:"storage"`
	StorageClass string                       `json:"storage_class"`
	Type         string                       `json:"type"`
	VCPUs        int                          `json:"vcpus"`
}

type BillingDetails struct {
	Balance            float32 `json:"balance"`
	HourlySpendingRate float32 `json:"hourly_spending_rate"`
}
