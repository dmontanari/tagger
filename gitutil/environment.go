package gitutil

import (
	"os"
)

type Env struct {
	PAT      string
	Username string
}

func NewEnv() Env {

	e := Env{PAT: "", Username: ""}

	if os.Getenv("HTTP_PAT") != "" {
		e.PAT = os.Getenv("HTTP_PAT")
	}
	if os.Getenv("HTTP_USERNAME") != "" {
		e.Username = os.Getenv("HTTP_USERNAME")
	}

	return e

}
