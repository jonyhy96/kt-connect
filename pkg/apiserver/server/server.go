package server

import "github.com/jonyhy96/kt-connect/pkg/apiserver/common"

// Init ...
func Init(context common.Context) error {
	r := NewRouter(context)
	return r.Run(":8000")
}
