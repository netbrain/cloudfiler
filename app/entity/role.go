package entity

type Role struct {
	ID    int
	Name  string
	Users []User
}

func (r *Role) Equals(other interface{}) bool {
	switch o := other.(type) {
	case Role:
		return r.ID == o.ID
	case *Role:
		return r.ID == o.ID
	}
	return false
}
