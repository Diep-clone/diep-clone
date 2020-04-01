package lib

// Gun is
type Gun struct {
	Owner     *Object
	Bullet    Object
	Sx        float64
	Sy        float64
	Dir       float64
	Rdir      float64
	Bound     float64
	Reload    float64 // default delay
	WaitTime  float64 // first wait time
	DelayTime float64 // when shot
	ShotTime  float64 // click time to delay
	AutoShot  bool
}

func (g Gun) Shot() *Object {

	return nil
}

// New Gun1!!!!!!!!!!11!!!
func NewGun(value interface{}) *Gun {
	if g, ok := value.(Gun); ok {
		return &g
	} else {
		return nil
	}
}
