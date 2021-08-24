package tree

/*

### CMDB 拓扑

obj1: propA=a1, propB=b1, propC=c1, propD=d1
obj2: propA=a1, propB=b2, propC=c1, propD=d2
obj3: propA=a1, propB=b2, propC=c1, propD=d2
obj4: propA=a2, propB=b2, propC=c2, propD=d2
obj5: propA=a2, propB=b2, propC=c1, propD=d2
obj6:           propB=b2,           propD=d1


propA = [a1, a2, a3, unkown]  propA 的值列表

// propA 对应每个值的 cmdb 对象列表
propA_objects = [
  [ obj1, obj2 ],     // a1
  [ obj3, obj4 ],     // a2
  [ obj5 ],           // a3
  [ obj6 ],           // unkown
]

// 按照 [propA, propB, propC] key 顺序展示的每个值的树结构


root:
	- a1
		- b1
			- c1
				- obj1
		- b2
			- c1
				- obj2
				- obj3
	- a2
		- b2
			- c1
				- obj5
			- c2
				- obj4
	- unkown
		- b2
			- unkown
				- obj6


root:
	objects: [ obj1, obj2, obj3, obj4, obj5, obj6 ]
	childrenKey: propA
  children:
    - value: a1
      objects: [ obj1, obj2, obj3 ]
			childrenKey: propB
      children:
        - value: b1
          objects: [ obj1 ]
					childrenKey: propC
          children:
            - value: c1
              objects: [ obj1]
              children: false

        - value: b2
          objects: [ obj2, obj3 ]
          children:
            - value: c1
              objects: [ obj2, obj3]
              children: false

    - value: a2
      objects: [ obj4, obj5 ]
			childrenKey: propB
      children:
        - value: b2
          objects: [ obj4, obj5]
					childrenKey: propC
          children:
            - value: c1
              objects: [ obj5 ]
              children: false
            - value: c2
              objects: [ obj4 ]
              children: false


    - value: unkown
      objects: [ obj6 ]
			childrenKey: propB
      children:
        - value: b2
          objects: [ obj6 ]
					childrenKey: propC
          children:
            - value: unkown
              objects: [ obj6 ]
              children: false

*/

var mockData = `
[
	{
		"id": "obj1",
		"propA": "a1",
		"propB": "b1",
		"propC": "c1",
		"propD": "d1"
	},
	{
		"id": "obj2",
		"propA": "a1",
		"propB": "b2",
		"propC": "c1",
		"propD": "d2",
		"propE": {
			"propF": "ef1"
		}
	},
	{
		"id": "obj3",
		"propA": "a1",
		"propB": "b2",
		"propC": "c1",
		"propD": "d1",
		"propE": {
			"propF": "ef2"
		}
	},
	{
		"id": "obj4",
		"propA": "a2",
		"propB": "b1",
		"propC": "c1",
		"propD": "d1",
		"propE": {
			"propF": {
				"propG": "efg3"
			}
		}
	},
	{
		"id": "obj5",
		"propA": "a2",
		"propB": "b2",
		"propC": "c1",
		"propD": "d1",
		"propE": {
			"propF": 11
		}
	},
	{
		"id": "obj6",
		"propB": "b1",
		"propD": "d1"
	},
	{
		"uuid": "obj7",
		"propB": "b2",
		"propD": "d3"
	}
]
`
