package tree

type Node struct {
	Val string `json:"val" yaml:"val"`

	// Objects IDs
	Objects []string `json:"objects" yaml:"objects"`

	ChildrenKey string  `json:"chidrenKey" yaml:"chidrenKey"`
	Children    []*Node `json:"children" yaml:"children"`
}
