package errors

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestExtractErrorCode(t *testing.T) {
	type expected struct {
		code int
		err  error
	}
	type test struct {
		name     string
		err      string
		expected expected
	}
	var tests = []test{
		{
			"Bad gateway",
			"graphql: server returned a non-200 status code: 503",
			expected{
				503,
				nil,
			},
		},
		{
			"Not found",
			"graphql:Not Found: Not found error for rocketeer_turbot.grants ",
			expected{
				0,
				errors.Errorf("graphql:Not Found: Not found error for rocketeer_turbot.grants "),
			},
		},
		{
			"Permission denied",
			"graphql: Permission Denied: Insufficient Permissions for rocketeer_turbot.grants  ",
			expected{
				0,
				errors.Errorf("graphql: Permission Denied: Insufficient Permissions for rocketeer_turbot.grants  "),
			},
		},
		{
			"Status network authentication required",
			"graphql: server returned a non-200 status code: 511",
			expected{
				511,
				nil,
			},
		},
		{
			"gRPC error",
			"rpc error: code = Unavailable desc = transport is closing",
			expected{
				0,
				errors.Errorf("rpc error: code = Unavailable desc = transport is closing"),
			},
		},
		{
			"System error",
			"Index out of bound",
			expected{
				0,
				errors.Errorf("Index out of bound"),
			},
		},
		{
			"Bad formatting",
			"graphql: server returned a non-200 status code:       511      ",
			expected{
				511,
				errors.Errorf("Index out of bound"),
			},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		errCode, err := ExtractErrorCode(errors.Errorf(test.err))
		assert.Equal(t, test.expected.code, errCode)
		assert.ObjectsAreEqual(test.expected.err, err)
	}
}
