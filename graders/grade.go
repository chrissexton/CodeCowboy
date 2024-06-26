package graders

import (
	"cso/codecowboy/classroom"
	"cso/codecowboy/graders/golang"
	"cso/codecowboy/graders/java"
	"cso/codecowboy/graders/net"
	"cso/codecowboy/store"
	"io"
)

type Grader interface {
	Grade(spec classroom.AssignmentSpec, out io.Writer) error
}

func GetGrader(language string, db *store.DB) Grader {
	switch language {
	case "go":
		return golang.NewGoGrader(db)
	case "java":
		return java.NewJavaGrader(db)
	case "net":
		return net.NewNetGrader(db)
	}
	return nil
}
