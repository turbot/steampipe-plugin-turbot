package apiClient

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/machinebox/graphql"
	"github.com/mitchellh/go-homedir"
	errorsHandler "github.com/turbot/steampipe-plugin-turbot/errors"
	"github.com/turbot/steampipe-plugin-turbot/helpers"
)

// Turbot API Client
type Client struct {
	AccessKey string
	SecretKey string
	Graphql   *graphql.Client
}

func CreateClient(config ClientConfig) (*Client, error) {
	// if accessKeyId and secretAccessKey were not directly specified (either via provider parameters or environment variables)
	// look for a credentials file

	credentials, err := GetCredentials(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials, error: %s", err.Error())
	}
	return &Client{
		AccessKey: credentials.AccessKey,
		SecretKey: credentials.SecretKey,
		Graphql:   graphql.NewClient(credentials.Workspace),
	}, nil
}

func GetCredentials(config ClientConfig) (ClientCredentials, error) {
	credentials, err := getCredentialsByPrecedence(config)
	if err != nil {
		return ClientCredentials{}, err
	}
	if !CredentialsSet(credentials) {
		return ClientCredentials{}, errors.New("failed to get credentials")
	}
	// update workspace url
	credentials.Workspace, err = BuildApiUrl(credentials.Workspace)
	if err != nil {
		return ClientCredentials{}, err
	}
	return credentials, nil
}

/*
	precedence of credentials:
	- Credentials set in config
	- profile set in config
	- ENV vars {TURBOT_ACCESS_KEY, TURBOT_SECRET_KEY, TURBOT_WORKSPACE}
	- TURBOT_PROFILE env var
*/
func getCredentialsByPrecedence(config ClientConfig) (ClientCredentials, error) {
	credentials := config.Credentials
	if !CredentialsSet(credentials) {
		var err error
		credentialsPath, err := getCredentialsPath(config)
		if err != nil {
			return ClientCredentials{}, err
		}
		if len(config.Profile) != 0 {
			credentials, err = getProfileCredentials(config)
			if err != nil {
				return ClientCredentials{}, err
			}
		} else {
			var credentialsOk bool
			credentials, credentialsOk = getCredentialsFromEnv()
			// if credentials were not passed in, get from the credentials file
			if !credentialsOk {
				config.Profile = os.Getenv("TURBOT_PROFILE")
				credentials, err = loadProfile(credentialsPath, config.Profile)
				if err != nil {
					return ClientCredentials{}, err
				}
			}
		}
	}
	return credentials, nil
}

func getCredentialsFromEnv() (ClientCredentials, bool) {
	credentials := ClientCredentials{
		AccessKey: os.Getenv("TURBOT_ACCESS_KEY"),
		SecretKey: os.Getenv("TURBOT_SECRET_KEY"),
		Workspace: os.Getenv("TURBOT_WORKSPACE"),
	}
	return credentials, CredentialsSet(credentials)
}

func getProfileCredentials(config ClientConfig) (ClientCredentials, error) {
	credentialsPath, err := getCredentialsPath(config)
	if err != nil {
		return ClientCredentials{}, err
	}
	credentials, err := loadProfile(credentialsPath, config.Profile)
	if err != nil {
		return ClientCredentials{}, err
	}
	return credentials, nil
}

func getCredentialsPath(config ClientConfig) (string, error) {
	var err error
	credentialsPath := config.CredentialsPath
	if len(credentialsPath) == 0 {
		credentialsPath = os.Getenv("TURBOT_SHARED_CREDENTIALS_FILE")
	}
	// if no credentials path was specified, use ~/.turbot/credentials
	if len(credentialsPath) == 0 {
		credentialsPath = filepath.Join(userHomeDir(), ".config", "turbot", "credentials.yml")
	} else {
		credentialsPath, err = homedir.Expand(credentialsPath)
		if err != nil {
			return "", err
		}
	}
	return credentialsPath, err
}

// convert workspace into a fully formed api url
func BuildApiUrl(rawWorkspace string) (string, error) {

	// acceptable forms of workspace are:
	// bananaman-turbot.putney
	// bananaman-turbot.putney.turbot.io
	// bananaman-turbot.putney.turbot.io/
	// bananaman-turbot.putney.turbot.io/api/v5
	// bananaman-turbot.putney.turbot.io/api/v5/
	// https://bananaman-turbot.putney.turbot.io
	// https://bananaman-turbot.putney.turbot.io/api/v5

	workspace := strings.TrimSuffix(rawWorkspace, "/")

	// check for "https://"' prefix
	if !strings.HasPrefix(workspace, "https://") {
		workspace = "https://" + workspace
	}
	u, err := url.Parse(workspace)
	if err != nil {
		return "", fmt.Errorf("failed to create client - could not parse workspace url %s, error %s", rawWorkspace, err.Error())
	}
	if u.Path == "invalid" {
		return "", fmt.Errorf("failed to create client - could not parse workspace url '%s'", rawWorkspace)
	}

	if u.Path != "" {
		apiVersionRegex := regexp.MustCompile(`\/api\/v[0-9]+$|latest$`)
		if !apiVersionRegex.Match([]byte(u.Path)) {
			return "", fmt.Errorf("invalid worksapce %s", workspace)
		}
		u.Path = path.Join(u.Path, "graphql")
	} else {
		u.Path = "/api/latest/graphql"
	}

	baseUrl := u.String()
	return baseUrl, nil
}

func CredentialsSet(credentials ClientCredentials) bool {
	return len(credentials.AccessKey) != 0 && len(credentials.SecretKey) != 0 && len(credentials.Workspace) != 0
}

// Validate checks if the API workspace URL and credentials are valid.
func (client *Client) Validate() error {
	query, responseObject := validationQuery()
	err := client.doRequest(query, nil, &responseObject)
	if err == nil && !responseObject.isValid() {
		err = errors.New("authorisation failed. Verify workspace, access_key and secret_key have been set correctly")
	}
	return err
}

// UserHomeDir returns the home directory for the user the process is running under.
func userHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}

func loadProfile(credentialsPath, profile string) (ClientCredentials, error) {
	// if no profile specified, use default
	if len(profile) == 0 {
		profile = "default"
	}
	yamlFile, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return ClientCredentials{}, err
	}

	var credentialsMap = map[string]ClientCredentials{}
	err = yaml.Unmarshal(yamlFile, &credentialsMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	credentials := credentialsMap[profile]
	if !CredentialsSet(credentials) {
		return ClientCredentials{}, fmt.Errorf("failed to load all credentials for profile %s from credentials file %s", profile, credentialsPath)
	}

	return credentials, nil
}

func basicAuthHeader(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func (client *Client) BuildPropertiesFromUpdateSchema(resourceId string, properties []interface{}) ([]interface{}, error) {
	getResourceQuery := getResourceTypeIdQuery(resourceId)
	responseData := &ResourceResponse{}
	// execute api call
	if err := client.doRequest(getResourceQuery, nil, &responseData); err != nil {
		return nil, fmt.Errorf("error reading resource type id: %s", err.Error())
	}

	resourceTypeId := responseData.Resource.Turbot.ResourceTypeId

	query := readResourceQuery(resourceTypeId, properties)
	response := &ResourceSchema{}
	// execute api call
	if err := client.doRequest(query, nil, &response); err != nil {
		return nil, fmt.Errorf("error reading resource type id: %s", err.Error())
	}

	if response.Resource.UpdateSchema == nil {
		return nil, nil
	}

	var m = response.Resource.UpdateSchema.(map[string]interface{})
	var excluded []interface{}
	if value, ok := m["allOf"]; ok {
		for _, schema := range value.([]interface{}) {
			if res, ok := schema.(map[string]interface{}); ok {
				if res["type"] == "object" {
					// loop to flatten interface, so we will not get this structure - [[id1,id2],[id3,id4]]
					for _, element := range helpers.GetNullProperties(res) {
						excluded = append(excluded, element)
					}
				}
			}
		}
	}
	return excluded, nil
}

func (client *Client) doRequest(query string, vars map[string]interface{}, responseData interface{}) error {
	return client.DoRequest(query, vars, responseData)
}

// execute graphql request
func (client *Client) DoRequest(query string, vars map[string]interface{}, responseData interface{}) error {
	// make a request
	req := graphql.NewRequest(query)

	// set any variables
	for k, v := range vars {
		req.Var(k, v)
	}

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", basicAuthHeader(client.AccessKey, client.SecretKey))

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	start := time.Now()
	if err := client.Graphql.Run(ctx, req, &responseData); err != nil {
		err = errorsHandler.BuildErrorMessage(err)
		return err
	}
	log.Println("graphql.time", time.Since(start).Milliseconds())
	return nil
}

func (client *Client) handleCreateError(err error, input map[string]interface{}, resourceType string) error {
	parent := input["parent"]
	if errorsHandler.NotFoundError(err) {
		return fmt.Errorf("error creating %s: parent resource not found: %s", resourceType, parent)
	}
	return fmt.Errorf("error creating %s: %s ", resourceType, err.Error())
}

func (client *Client) handleReadError(err error, resource string, resourceType string) error {
	if errorsHandler.NotFoundError(err) {
		return fmt.Errorf("error reading %s: resource not found: %s", resourceType, resource)
	}
	return fmt.Errorf("error reading %s: %s ", resourceType, err.Error())
}

func (client *Client) handleUpdateError(err error, input map[string]interface{}, resourceType string) error {
	resource := input["id"]
	if errorsHandler.NotFoundError(err) {
		return fmt.Errorf("error updating %s: resource not found: %s", resourceType, resource)
	}
	return fmt.Errorf("error updating %s: %s ", resourceType, err.Error())
}
