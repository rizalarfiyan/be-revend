package models

type UserPoint struct {
	Identity string `json:"user_id"`
	DeviceId string `json:"device_id"`
	Success  int    `json:"success"`
	Failed   int    `json:"failed"`
}
