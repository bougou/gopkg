package tree

import (
	"errors"
	"fmt"
)

// Treeify generates a tree for a list of objects by their attributes.
// The first argument is a list of objects.
// The second argument is the object key name which represents the uniqueness of the object.
// The left arguments are tree layers.
func Treeify(root string, objects []map[string]interface{}, idKey string, props ...string) (*Node, error) {
	objectsByID := map[string]map[string]interface{}{}
	objectIDs := []string{}

	for _, object := range objects {
		if objID, ok := object[idKey]; ok {
			id, ok := objID.(string)
			if !ok {
				// the value of the key of the object is not string, ignore this object
				continue
			}
			objectsByID[id] = object
			objectIDs = append(objectIDs, id)
		} else {
			// the object does not have key, ignore this object
			continue
		}
	}

	out := &Node{
		Val:         root,
		Objects:     objectIDs,
		ChildrenKey: "",
		Children:    make([]*Node, 0),
	}

	if len(props) == 0 {
		return out, nil
	}

	var propKey string
	var propsLeft []string

	propKey = props[0]
	if len(props) == 1 {
		propsLeft = []string{}
	} else {
		propsLeft = props[1:]
	}
	out.ChildrenKey = propKey

	groups, err := GroupObjectsByPropValue(objects, idKey, propKey)
	if err != nil {
		msg := fmt.Sprintf("GroupObjectsByProp failed, err: %s", err)
		return nil, errors.New(msg)
	}

	for propValue, childrenObjects := range groups {
		node, err := Treeify(propValue, childrenObjects, idKey, propsLeft...)
		if err != nil {
			msg := fmt.Sprintf("Treeify failed, err: %s", err)
			return nil, errors.New(msg)
		}
		if node == nil {
			continue
		}
		out.Children = append(out.Children, node)
	}

	return out, nil
}

// GroupObjectsByProp groups `objects` by specified `prop` key's value, the `idKey` is the unique identifier key of the objects
func GroupObjectsByPropValue(objects []map[string]interface{}, idKey string, propKey string) (map[string][]map[string]interface{}, error) {
	// store the correspoding objects whose prop's value equals to the key of this map
	valuesMap := make(map[string][]map[string]interface{})

	updateValuesMap := func(object map[string]interface{}, valueKey string) {
		if _, exists := valuesMap[valueKey]; !exists {
			valuesMap[valueKey] = make([]map[string]interface{}, 0)
		}
		valuesMap[valueKey] = append(valuesMap[valueKey], object)
	}

	for _, object := range objects {
		propValue, ok := object[propKey]
		if !ok {
			propValue = "Unknown"
		}

		_propValue, ok := propValue.(string)
		if ok {
			propValue = "Unknown"
		}

		updateValuesMap(object, _propValue)
	}

	return valuesMap, nil
}
