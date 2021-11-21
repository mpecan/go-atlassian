package v2

import (
	"context"
	"fmt"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"net/http"
	"strings"
)

type FilterShareService struct{ client *Client }

type ShareFilterScopeScheme struct {
	Scope string `json:"scope"`
}

// Scope returns the default sharing settings for new filters and dashboards for a user.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/filters/sharing#get-default-share-scope
// Official Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-filter-sharing/#api-rest-api-2-filter-defaultsharescope-get
func (f *FilterShareService) Scope(ctx context.Context) (result *ShareFilterScopeScheme, response *ResponseScheme, err error) {

	var endpoint = "rest/api/2/filter/defaultShareScope"
	request, err := f.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = f.client.call(request, &response)
	if err != nil {
		return
	}

	return
}

// SetScope sets the default sharing for new filters and dashboards for a user.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/filters/sharing#set-default-share-scope
// Official Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-filter-sharing/#api-rest-api-2-filter-defaultsharescope-put
// Valid values: GLOBAL, AUTHENTICATED, PRIVATE
func (f *FilterShareService) SetScope(ctx context.Context, scope string) (response *ResponseScheme, err error) {

	//Valid the share filter scope
	var (
		validScopeValuesAsList = []string{"GLOBAL", "AUTHENTICATED", "PRIVATE"}
		isValid                bool
	)

	for _, validScope := range validScopeValuesAsList {
		if validScope == scope {
			isValid = true
			break
		}
	}

	if !isValid {
		//Join the valid values and create the custom error
		var validScopeValuesAsString = strings.Join(validScopeValuesAsList, ",")
		return nil, fmt.Errorf("invalid scope, please provide one of the following: %v", validScopeValuesAsString)
	}

	payload := ShareFilterScopeScheme{Scope: scope}
	payloadAsReader, _ := transformStructToReader(&payload)

	var endpoint = "rest/api/2/filter/defaultShareScope"

	request, err := f.client.newRequest(ctx, http.MethodPut, endpoint, payloadAsReader)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err = f.client.call(request, nil)
	if err != nil {
		return
	}

	return
}

// Gets returns the share permissions for a filter.
// A filter can be shared with groups, projects, all logged-in users, or the public.
// Sharing with all logged-in users or the public is known as a global share permission.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/filters/sharing#get-share-permissions
// Official Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-filter-sharing/#api-rest-api-2-filter-id-permission-get
func (f *FilterShareService) Gets(ctx context.Context, filterID int) (result []*models.SharePermissionScheme, response *ResponseScheme, err error) {

	var endpoint = fmt.Sprintf("rest/api/2/filter/%v/permission", filterID)

	request, err := f.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = f.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

type PermissionFilterPayloadScheme struct {
	Type          string `json:"type,omitempty"`
	ProjectID     string `json:"projectId,omitempty"`
	GroupName     string `json:"groupname,omitempty"`
	ProjectRoleID string `json:"projectRoleId,omitempty"`
}

// Add a share permissions to a filter.
// If you add a global share permission (one for all logged-in users or the public)
// it will overwrite all share permissions for the filter.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/filters/sharing#add-share-permission
func (f *FilterShareService) Add(ctx context.Context, filterID int, payload *PermissionFilterPayloadScheme) (
	result []*models.SharePermissionScheme, response *ResponseScheme, err error) {

	payloadAsReader, err := transformStructToReader(payload)
	if err != nil {
		return nil, nil, err
	}

	var endpoint = fmt.Sprintf("rest/api/2/filter/%v/permission", filterID)

	request, err := f.client.newRequest(ctx, http.MethodPost, endpoint, payloadAsReader)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err = f.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// Get returns a share permission for a filter.
// A filter can be shared with groups, projects, all logged-in users, or the public.
// Sharing with all logged-in users or the public is known as a global share permission.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/filters/sharing#get-share-permission
// Official Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-filter-sharing/#api-rest-api-2-filter-id-permission-permissionid-get
func (f *FilterShareService) Get(ctx context.Context, filterID, permissionID int) (result *models.SharePermissionScheme,
	response *ResponseScheme, err error) {

	var endpoint = fmt.Sprintf("rest/api/2/filter/%v/permission/%v", filterID, permissionID)
	request, err := f.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	response, err = f.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// Delete deletes a share permission from a filter.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/filters/sharing#delete-share-permission
// Official Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-filter-sharing/#api-rest-api-2-filter-id-permission-permissionid-delete
func (f *FilterShareService) Delete(ctx context.Context, filterID, permissionID int) (response *ResponseScheme, err error) {

	var endpoint = fmt.Sprintf("rest/api/2/filter/%v/permission/%v", filterID, permissionID)
	request, err := f.client.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return
	}

	response, err = f.client.call(request, nil)
	if err != nil {
		return
	}

	return
}