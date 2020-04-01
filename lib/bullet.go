package lib

import (
	"math"
	"time"
)

type Bullet struct {
	Object
}

func (b *Bullet) Tick() {
	b.Dx += math.Cos(b.Dir) * b.Speed
	b.Dy += math.Sin(b.Dir) * b.Speed
}

func NewBullet(value interface{}) *Bullet {
	if b, ok := value.(Bullet); ok {
		b.ID = objID
		objID++
		b.Opacity = 0.5
		b.SpawnTime = time.Now().Unix()
		b.HitTime = time.Now().Unix()
		b.DeadTime = -1
		//b.HitObject = &b
		b.IsBorder = false
		b.IsOwnCol = false
		b.IsDead = false
		b.IsCollision = false

		return &b
	} else {
		return nil
	}
}
