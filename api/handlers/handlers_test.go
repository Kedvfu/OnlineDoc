package handlers

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestHandleLogin(t *testing.T) {
	type args struct {
		context *gin.Context
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleLogin(tt.args.context)

		})
	}
}
