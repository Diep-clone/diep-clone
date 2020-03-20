package lib

func (obj *Object) Tank() {
	obj.Type = "Tank"
	obj.IsBorder = true
	obj.IsOwnCol = false
}

// Basic is god
func (obj *Object) Basic() {

}
