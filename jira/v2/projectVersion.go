package v2

import (
	"context"
	"fmt"
	models2 "github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ProjectVersionService struct{ client *Client }

// Gets returns all versions in a project.
// The response is not paginated.
// Use Get project versions paginated if you want to get the versions in a project with pagination.
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-project-projectidorkey-versions-get
func (p *ProjectVersionService) Gets(ctx context.Context, projectKeyOrID string) (result []*models2.VersionScheme, response *ResponseScheme, err error) {

	if len(projectKeyOrID) == 0 {
		return nil, nil, models2.ErrNoProjectIDError
	}

	endpoint := fmt.Sprintf("rest/api/2/project/%v/versions", projectKeyOrID)

	request, err := p.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = p.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// Search returns a paginated list of all versions in a project.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/projects/versions#get-project-versions-paginated
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-project-projectidorkey-version-get
func (p *ProjectVersionService) Search(ctx context.Context, projectKeyOrID string, options *models2.VersionGetsOptions, startAt,
	maxResults int) (result *models2.VersionPageScheme, response *ResponseScheme, err error) {

	if len(projectKeyOrID) == 0 {
		return nil, nil, models2.ErrNoProjectIDError
	}

	params := url.Values{}
	params.Add("startAt", strconv.Itoa(startAt))
	params.Add("maxResults", strconv.Itoa(maxResults))

	if options != nil {

		if len(options.Expand) != 0 {
			params.Add("expand", strings.Join(options.Expand, ","))
		}

		if len(options.Query) != 0 {
			params.Add("query", options.Query)
		}

		if len(options.Status) != 0 {
			params.Add("status", options.Status)
		}

		if len(options.OrderBy) != 0 {
			params.Add("orderBy", options.OrderBy)
		}

	}

	var endpoint = fmt.Sprintf("rest/api/2/project/%v/version?%v", projectKeyOrID, params.Encode())

	request, err := p.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = p.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// Create creates a project version.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/projects/versions#create-version
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-version-post
func (p *ProjectVersionService) Create(ctx context.Context, payload *models2.VersionPayloadScheme) (
	result *models2.VersionScheme, response *ResponseScheme, err error) {

	var endpoint = "rest/api/2/version"

	payloadAsReader, err := transformStructToReader(payload)
	if err != nil {
		return nil, nil, err
	}

	request, err := p.client.newRequest(ctx, http.MethodPost, endpoint, payloadAsReader)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err = p.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// Get returns a project version.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/projects/versions#get-version
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-version-id-get
func (p *ProjectVersionService) Get(ctx context.Context, versionID string, expand []string) (result *models2.VersionScheme,
	response *ResponseScheme, err error) {

	if len(versionID) == 0 {
		return nil, nil, models2.ErrNoVersionIDError
	}

	params := url.Values{}
	if len(expand) != 0 {
		params.Add("expand", strings.Join(expand, ","))
	}

	var endpoint strings.Builder
	endpoint.WriteString(fmt.Sprintf("rest/api/2/version/%v", versionID))

	if params.Encode() != "" {
		endpoint.WriteString(fmt.Sprintf("?%v", params.Encode()))
	}

	request, err := p.client.newRequest(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = p.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// Update updates a project version.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/projects/versions#update-version
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-version-id-put
func (p *ProjectVersionService) Update(ctx context.Context, versionID string, payload *models2.VersionPayloadScheme) (
	result *models2.VersionScheme, response *ResponseScheme, err error) {

	if len(versionID) == 0 {
		return nil, nil, models2.ErrNoVersionIDError
	}

	var endpoint = fmt.Sprintf("rest/api/2/version/%v", versionID)

	payloadAsReader, err := transformStructToReader(payload)
	if err != nil {
		return nil, nil, err
	}

	request, err := p.client.newRequest(ctx, http.MethodPut, endpoint, payloadAsReader)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err = p.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// Merge merges two project versions.
// The merge is completed by deleting the version specified in id and replacing any occurrences of its ID in fixVersion with the version ID specified in moveIssuesTo.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/projects/versions#merge-versions
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-version-id-mergeto-moveissuesto-put
func (p *ProjectVersionService) Merge(ctx context.Context, versionID, versionMoveIssuesTo string) (response *ResponseScheme,
	err error) {

	if len(versionID) == 0 {
		return nil, models2.ErrNoVersionIDError
	}

	if len(versionMoveIssuesTo) == 0 {
		return nil, models2.ErrNoVersionIDError
	}

	var endpoint = fmt.Sprintf("rest/api/2/version/%v/mergeto/%v", versionID, versionMoveIssuesTo)

	request, err := p.client.newRequest(ctx, http.MethodPut, endpoint, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = p.client.call(request, nil)
	if err != nil {
		return
	}

	return
}

// RelatedIssueCounts returns the following counts for a version:
// 1. Number of issues where the fixVersion is set to the version.
// 2. Number of issues where the affectedVersion is set to the version.
// 2. Number of issues where a version custom field is set to the version.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/projects/versions#get-versions-related-issues-count
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-version-id-relatedissuecounts-get
func (p *ProjectVersionService) RelatedIssueCounts(ctx context.Context, versionID string) (result *models2.VersionIssueCountsScheme,
	response *ResponseScheme, err error) {

	if len(versionID) == 0 {
		return nil, nil, models2.ErrNoVersionIDError
	}

	var endpoint = fmt.Sprintf("rest/api/2/version/%v/relatedIssueCounts", versionID)

	request, err := p.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = p.client.call(request, &result)
	if err != nil {
		return
	}

	return
}

// UnresolvedIssueCount returns counts of the issues and unresolved issues for the project version.
// Docs: https://docs.go-atlassian.io/jira-software-cloud/projects/versions#get-versions-unresolved-issues-count
// Atlassian Docs: https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-project-versions/#api-rest-api-2-version-id-unresolvedissuecount-get
func (p *ProjectVersionService) UnresolvedIssueCount(ctx context.Context, versionID string) (
	result *models2.VersionUnresolvedIssuesCountScheme, response *ResponseScheme, err error) {

	if len(versionID) == 0 {
		return nil, nil, models2.ErrNoVersionIDError
	}

	var endpoint = fmt.Sprintf("rest/api/2/version/%v/unresolvedIssueCount", versionID)

	request, err := p.client.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept", "application/json")

	response, err = p.client.call(request, &result)
	if err != nil {
		return
	}

	return
}