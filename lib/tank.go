package lib

import (
	"time"
)

type Tank struct {
	Controller *Player
	Object
	Sight    float64
	MaxStats [8]int
}

func NewTank(value interface{}) *Tank {
	if t, ok := value.(Tank); ok {
		t.ID = objID
		objID++
		t.Mh = t.H
		t.Lh = t.H
		t.Bound = 1
		t.Stance = 1
		t.Opacity = 1
		t.SpawnTime = time.Now().Unix()
		t.HitTime = time.Now().Unix()
		t.DeadTime = -1
		//b.HitObject = &b
		t.IsBorder = true
		t.IsOwnCol = false
		t.IsDead = false
		t.IsCollision = false
		return &t
	} else {
		return nil
	}
}
