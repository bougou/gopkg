package merge

type Map struct {
	items map[string]interface{}

	// Whether to use src field value to overrite the dst field value.
	// Default false.
	overwrite bool

	// Whether to add the src field to the dst if the field does not exist in the dst.
	// Default true.
	newField bool
}

func NewMap(m map[string]interface{}) *Map {
	return &Map{
		items: m,

		overwrite: false,
		newField:  true,
	}
}

func (m *Map) Value() map[string]interface{} {
	return m.items
}

type Option func(m *Map) *Map

func WithOverride(m *Map) *Map {
	m.overwrite = true
	return m
}

func WithNoNewField(m *Map) *Map {
	m.newField = false
	return m
}

func (m *Map) Merge(in map[string]interface{}, options ...Option) error {
	for _, option := range options {
		option(m)
	}

	for k, v := range in {
		_, exists := m.items[k]
		if !exists && m.newField {
			m.items[k] = v
			continue
		}

		srcValue, ok1 := v.(map[string]interface{})
		dstValue, ok2 := m.items[k].(map[string]interface{})
		if ok1 && ok2 {
			mm := NewMap(dstValue)
			if err := mm.Merge(srcValue, options...); err != nil {
				return err
			}
			continue
		}

		if m.overwrite {
			m.items[k] = v
		}
	}

	return nil
}
