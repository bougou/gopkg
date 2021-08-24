package tree

// Object represents
type Object interface {
	GetID() (string, bool)
	GetProp(propKey string) (string, bool)
	GetNestedProp(nestedPropKeys ...string) (string, bool)
}

type MapObject struct {
	idKey string
	value map[string]interface{}
}

func NewMapObject(idKey string, value map[string]interface{}) *MapObject {
	return &MapObject{
		idKey: idKey,
		value: value,
	}
}

func (mo *MapObject) GetNestedProp(nestedPropKeys ...string) (string, bool) {
	tmp := mo

	for i, propKey := range nestedPropKeys {
		v, exists := tmp.value[propKey]
		if !exists {
			return "", false
		}

		if i == len(nestedPropKeys)-1 {
			s, ok := v.(string)
			if !ok {
				return "", false
			}
			return s, true
		}

		m, ok := v.(map[string]interface{})
		if !ok {
			return "", false
		}

		tmp = NewMapObject("id", m)
	}

	return "", false
}

func (mo *MapObject) GetProp(propKey string) (string, bool) {
	v, exists := mo.value[propKey]
	if !exists {
		return "", false
	}

	s, ok := v.(string)
	if !ok {
		return "", false
	}

	return s, true
}

func (mo *MapObject) GetID() (string, bool) {
	return mo.GetProp(mo.idKey)
}
