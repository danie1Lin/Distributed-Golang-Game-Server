package entity

type Treasure struct {
	Entity
	lifeTime int64
}

func (t *Treasure) Tick() {
	t.Disapear()
}

func (t *Treasure) Disapear() {
	t.lifeTime++
	if t.lifeTime > 1000 {
		t.Destroy()
	}
}

func (t *Treasure) BeCollected() {

}

func (t *Treasure) BeAttracted() {

}
