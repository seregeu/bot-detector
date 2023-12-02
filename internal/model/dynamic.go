package model

type Dynamic struct {
	ID                 int64   `json:"id"`
	UserID             int64   `json:"user_id"`
	MaxDeviceOffs      float64 `json:"max_device_offs"`
	MinDeviceOffs      float64 `json:"min_device_offs"`
	MaxDevAcceleration float64 `json:"max_dev_acceleration"`
	MinDevAcceleration float64 `json:"min_dev_acceleration"`
	MinLight           float64 `json:"min_light"`
	MaxLight           float64 `json:"max_light"`
	HitY               float64 `json:"hit_y"`
	HitX               float64 `json:"hit_x"`
}
