package apiClient

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCredentialsPrecedence(t *testing.T) {
	type expected struct {
		result bool
		Creds ClientCredentials
	}
	type test struct {
		name string
		Config   ClientConfig
		expected expected
	}
	var tests = []test{
		{
			"Config has credentials",
			ClientConfig{
				ClientCredentials{
					"xxbd857-XXXX-XXXX-XXXX-xxxxx039ff1x",
					"36xxb4f-XXXX-XXXX-XXXX-c91f44axx4f6",
					"https://example.com/",
				},
				"",
				"",
			},
			expected{
				true,
				ClientCredentials{
					"xxbd857-XXXX-XXXX-XXXX-xxxxx039ff1x",
					"36xxb4f-XXXX-XXXX-XXXX-c91f44axx4f6",
					"https://example.com/",
				},
			},
		},
		{
			"Config has profile",
			ClientConfig{
				ClientCredentials{
					"",
					"",
					"",
				},
				"",
				"test",
			},
			expected{
				true,
				ClientCredentials{
					"xxbd857-XXXX-XXXX-XXXX-xxxxx039ff1x",
					"36xxb4f-XXXX-XXXX-XXXX-c91f44axx4f6",
					"https://example.com/",
				},
			},
		},
		{
			"Empty Config",
			ClientConfig{
				ClientCredentials{
					"",
					"",
					"",
				},
				"",
				"test",
			},
			expected{
				true,
				ClientCredentials{
					os.Getenv("TURBOT_ACCESS_KEY"),
					os.Getenv("TURBOT_SECRET_KEY"),
					os.Getenv("TURBOT_WORKSPACE"),
				},
			},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		credentials, _ := GetCredentials(test.Config)
		assert.Equal(t, test.expected.result, CredentialsSet(credentials))
		assert.ObjectsAreEqual(test.expected.Creds, credentials)
	}
}
