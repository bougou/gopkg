package regpat

// Note, to avoid escaping issues
// use [.]     instead of   \.
// use [0-9]   instead of   \d

const (
	// Would:
	//  - 123     (√)
	//  - 123.456 (x)
	//  - .456    (x)
	//  - 123.    (x)
	NumberInteger string = `[+-]?[0-9]+`

	// Would:
	//  - 123     (x)
	//  - 123.456 (√)
	//  - .456    (x)
	//  - 123.    (x)
	NumberFloat string = `[+-]?[0-9]+[.][0-9]+`

	// Would:
	//  - 123     (√)
	//  - 123.456 (√)
	//  - .456    (x)
	//  - 123.    (x)
	Number string = `[+-]?[0-9]+([.][0-9]+)?`

	// Would:
	//  - 123     (√)
	//  - 123.456 (√)
	//  - .456    (√)
	//  - 123.    (x)
	NumberAllowNoInteger string = `[+-]?([0-9]*[.])?[0-9]+`

	// Would:
	//  - 123     (√)
	//  - 123.456 (√)
	//  - .456    (√)
	//  - 123.    (√)
	NumberAllowNoDecimal string = `[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+)`
)
