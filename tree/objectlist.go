package tree

import (
	"errors"
	"fmt"
)

type ObjectList []Object

// GetObjectIDs return a list of object id.
// The length of output slice may be shorter than length of the receiver ObjectList.
func (objects *ObjectList) GetObjectIDs() []string {
	objectIDs := make([]string, 0)
	for _, object := range *objects {
		if objID, ok := object.GetID(); ok {
			if !ok {
				continue
			}
			objectIDs = append(objectIDs, objID)
		}
	}
	return objectIDs
}

func (objects *ObjectList) Tree(root string, props ...string) (*Node, error) {
	objectIDs := objects.GetObjectIDs()
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

	groups, err := objects.GroupByPropValue(propKey)
	if err != nil {
		msg := fmt.Sprintf("GroupObjectsByProp failed, err: %s", err)
		return nil, errors.New(msg)
	}

	groupKeys := make([]string, 0)
	for k := range groups {
		groupKeys = append(groupKeys, k)
	}

	for propValue, objectList := range groups {
		node, err := objectList.Tree(propValue, propsLeft...)
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

// GroupByPropValue return a map, which groups the ObjectList by specified propKey's value.
// The returned map keys are values of the propKeys.
func (objects *ObjectList) GroupByPropValue(propKey string) (map[string]ObjectList, error) {
	// store the correspoding objects whose prop's value equals to the key of this map
	valuesMap := map[string]ObjectList{}

	updateValuesMap := func(object Object, valueKey string) {
		if _, exists := valuesMap[valueKey]; !exists {
			valuesMap[valueKey] = make(ObjectList, 0)
		}
		valuesMap[valueKey] = append(valuesMap[valueKey], object)
	}

	for _, object := range *objects {
		propValue, ok := object.GetProp(propKey)
		if !ok {
			propValue = "Unkown"
		}
		updateValuesMap(object, propValue)
	}

	return valuesMap, nil
}
