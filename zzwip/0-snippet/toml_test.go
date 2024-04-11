package snippet

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kr/pretty"
	toml "github.com/pelletier/go-toml"
)

// https://github.com/pelletier/go-toml/blob/master/marshal.go#L219
// The following struct annotations are supported:
//   toml:"Field"      Overrides the field's name to output.
//   omitempty         When set, empty values and groups are not emitted.
//   comment:"comment" Emits a # comment on the same line. This supports new lines.
//   commented:"true"  Emits the value as commented.
// Note that pointers are automatically assigned the "omitempty" option, as TOML
// explicitly does not handle null values (saying instead the label should be
// dropped).

func t() {
	buf := new(bytes.Buffer)

	toml.NewEncoder(buf).
		QuoteMapKeys(false).
		ArraysWithOneElementPerLine(false).
		Indentation("").
		Order(toml.OrderAlphabetical).
		Order(toml.OrderPreserve).
		SetTagName("toml").
		SetTagMultiline("multiline").
		SetTagComment("comment").
		SetTagCommented("commented").
		PromoteAnonymous(false)

}

type Configuration struct {
	Inputs  Inputs  `toml:"inputs"`
	Outputs Outputs `toml:"output"`
}

type Inputs struct {
	Postgres []Postgres `toml:"postgres"`
	MySQL    []MySQL    `toml:"mysql"`
	ES       ES         `toml:"es"`
}

type Outputs struct {
	Postgres []Postgres `toml:"postgres"`
	MySQL    []MySQL    `toml:"mysql"`
}

type ES struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type Postgres struct {
	User     string `toml:"user" comment:"username"`
	Password string `toml:"pass" comment:"password" commented:"true"`
}

type MySQL struct {
	User     string `toml:"user" comment:"username"`
	Password string `toml:"pass" comment:"password" commented:"true"`
}

type Config struct {
	Postgres Postgres `toml:"postgres"`
	MySQL    MySQL    `toml:"mysql"`
}

func Test_loadToml(t *testing.T) {
	data := `
	[postgres]
	user = "pelletier"
	pass = "mypassword"`

	config, _ := toml.Load(data)

	// retrieve data directly
	user := config.Get("postgres.user").(string)
	fmt.Println("user: ", user)

	// or using an intermediate object
	postgresConfig := config.Get("postgres").(*toml.Tree)
	password := postgresConfig.Get("pass").(string)
	fmt.Println("password: ", password)

	c := Config{}
	toml.Unmarshal([]byte(data), &c)
	pretty.Print(c)
}

func Test_saveToml(t *testing.T) {
	c1 := Config{
		Postgres: Postgres{
			User:     "john",
			Password: "123456",
		},
	}
	s, _ := toml.Marshal(c1)
	fmt.Println(string(s))

	// # postgres
	// [postgres]

	// 	# password
	// 	# pass = "123456"

	// 	# username
	// 	user = "john"

}

func loadToml3() {
	data := `
	[[inputs.mysql]]
	user = "mysql_user1"
	pass = "mysql_pass1"

	[[inputs.mysql]]
	user = "mysql_user2"
	pass = "mysql_pass2"

	[[inputs.postgres]]
	user = "postgres_user1"
	pass = "postgres_pass1"

	[[inputs.postgres]]
	user = "postgres_user2"
	pass = "postgres_pass2"

	[[output.mysql]]
	user = "mysql_user1"
	pass = "mysql_pass1"

	[[output.mysql]]
	user = "mysql_user2"
	pass = "mysql_pass2"

	[[output.postgres]]
	user = "postgres_user1"
	pass = "postgres_pass1"

	[[output.postgres]]
	user = "postgres_user2"
	pass = "postgres_pass2"
	`

	c := Configuration{}
	toml.Unmarshal([]byte(data), &c)
	fmt.Println(c)
	fmt.Println(c.Inputs)
	fmt.Println(c.Inputs.MySQL)
	fmt.Println(c.Outputs)

}

func saveToml3() {
	d3 := Configuration{
		Inputs: Inputs{
			MySQL: []MySQL{
				{
					User:     "mysql_user1",
					Password: "mysql_pass1",
				},
				{
					User:     "mysql_user1",
					Password: "mysql_pass1",
				},
			},
			ES: ES{
				User:     "es_user",
				Password: "es_pass",
			},
		},
	}

	s3, _ := toml.Marshal(d3)
	fmt.Println("s3: ", string(s3))
	fmt.Println("s3:------")
}

func loadToml2() {
	data := `
		[[inputs.postgres]]

		# password
		# pass = "123456"

		# username
		user = "john"

		[[inputs.postgres]]

			# password
			# pass = "654321"

			# username
			user = "tom"
	`

	c := Configuration{}
	toml.Unmarshal([]byte(data), &c)

	fmt.Println("c: ", c)
	fmt.Println(c.Inputs)
	fmt.Println(c.Inputs.MySQL)
	fmt.Println(c.Inputs.Postgres)

}

func saveToml2() {
	c2 := Configuration{
		Inputs: Inputs{
			Postgres: []Postgres{
				{
					User:     "john",
					Password: "123456",
				},
				{
					User:     "tom",
					Password: "654321",
				},
			},
		},
	}

	s2, _ := toml.Marshal(c2)
	fmt.Println("c2: ", string(s2))
	fmt.Println("c2:------")
}
