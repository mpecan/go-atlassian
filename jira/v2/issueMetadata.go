package v2

import (
	"context"
	"fmt"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/google/go-querystring/query"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
	"strings"
)

type IssueMetadataService struct{ client *Client }

// Get edit issue metadata returns the edit screen fields for an issue that are visible to and editable by the user.
// Use the information to populate the requests in Edit issue.
// Atlassian URL: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issues/#api-rest-api-2-issue-issueidorkey-editmeta-get
func (i *IssueMetadataService) Get(ctx context.Context, issueKeyOrID string, overrideScreenSecurity, overrideEditableFlag bool) (
	result gjson.Result, response *ResponseScheme, err error) {

	if len(issueKeyOrID) == 0 {
		return gjson.Result{}, nil, models.ErrNoIssueKeyOrIDError
	}

	params := url.Values{}

	if overrideEditableFlag {
		params.Add("overrideEditableFlag", "true")
	}

	if overrideScreenSecurity {
		params.Add("overrideScreenSecurity", "true")
	}

	var endpoint strings.Builder
	endpoint.WriteString(fmt.Sprintf("rest/api/2/issue/%v/editmeta", issueKeyOrID))

	if params.Encode() != "" {
		endpoint.WriteString(fmt.Sprintf("?%v", params.Encode()))
	}

	request, err := i.client.newRequest(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = i.client.call(request, nil)
	if err != nil {
		return
	}

	return gjson.ParseBytes(response.Bytes.Bytes()), response, nil
}

type IssueMetadataCreateOptions struct {
	ProjectIDs     []string `url:"projectIds,omitempty"`
	ProjectKeys    []string `url:"projectKeys,omitempty"`
	IssueTypeIDs   []string `url:"issuetypeIds,omitempty"`
	IssueTypeNames []string `url:"issuetypeNames,omitempty"`
	Expand         string   `url:"expand,omitempty"`
}

// Create returns details of projects, issue types within projects, and, when requested, the create screen fields for each issue type for the user.
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issues/#api-rest-api-2-issue-createmeta-get
func (i *IssueMetadataService) Create(ctx context.Context, opts *IssueMetadataCreateOptions) (result gjson.Result,
	response *ResponseScheme, err error) {

	params, _ := query.Values(opts)
	var endpoint strings.Builder
	endpoint.WriteString("rest/api/2/issue/createmeta")
	endpoint.WriteString(fmt.Sprintf("?%v", params.Encode()))

	request, err := i.client.newRequest(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = i.client.call(request, nil)
	if err != nil {
		return
	}

	return gjson.ParseBytes(response.Bytes.Bytes()), response, nil
}