package logs

import (
	"bytes"
	"edgefusion-data-push/common"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"text/template"
)

type Coder interface {
	Code() string
}

func New(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Trace(err error) error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case fmt.Formatter:
		return err
	default:
		return errors.WithStack(err)
	}
}

func CodeError(code, message string) error {
	return &codeError{errors.New(message), code}
}

type codeError struct {
	e error
	c string
}

func (e *codeError) Code() string {
	return e.c
}

func (e *codeError) Error() string {
	return e.e.Error()
}

func (e *codeError) Format(s fmt.State, verb rune) {
	e.e.(fmt.Formatter).Format(s, verb)
}

// Field field
type F struct {
	k string
	v interface{}
}

// Error returns an error with code and fields
func CustomError(c common.Code, fs ...*F) error {
	m := c.String()
	if strings.Contains(m, "{{") {
		vs := map[string]interface{}{}
		for _, f := range fs {
			vs[f.k] = f.v
		}
		t, err := template.New(string(c)).Option("missingkey=zero").Parse(m)
		if err != nil {
			panic(err)
		}
		b := bytes.NewBuffer(nil)
		err = t.Execute(b, vs)
		if err != nil {
			panic(err)
		}
		m = b.String()
	}
	return CodeError(string(c), m)
}

// Field returns a field
func CustomField(k string, v interface{}) *F {
	return &F{k, v}
}
