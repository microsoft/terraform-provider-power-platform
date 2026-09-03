// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func newApplicationClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *client) AddApplicationUser(ctx context.Context, environmentId string, applicationId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s/addAppUser", environmentId),
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.ADMIN_MANAGEMENT_APP_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	// Create the request body
	requestBody := map[string]string{
		"servicePrincipalAppId": applicationId,
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, requestBody, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *client) GetRootBusinessUnitId(ctx context.Context, environmentId string) (string, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return "", err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/businessunits", constants.DATAVERSE_API_VERSION),
	}
	values := url.Values{}
	values.Add("$select", "businessunitid,name")
	values.Add("$filter", "parentbusinessunitid eq null")
	apiUrl.RawQuery = values.Encode()

	var response applicationBusinessUnitArrayDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return "", err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return "", err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return "", err
	}
	if len(response.Value) == 0 {
		return "", customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("root business unit not found in environment '%s'", environmentId))
	}
	if len(response.Value) > 1 {
		return "", fmt.Errorf("expected exactly one root business unit in environment '%s', got %d", environmentId, len(response.Value))
	}

	return response.Value[0].BusinessUnitId, nil
}

func (client *client) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	if env.Properties.LinkedEnvironmentMetadata.InstanceURL == "" {
		return "", fmt.Errorf("environment %s does not have Dataverse", environmentId)
	}

	// Parse the instance URL to get the host
	instanceURL := env.Properties.LinkedEnvironmentMetadata.InstanceURL
	instanceURLParsed, err := url.Parse(instanceURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse instance URL %s: %v", instanceURL, err)
	}

	return instanceURLParsed.Host, nil
}

func (client *client) ApplicationUserExists(ctx context.Context, environmentId string, applicationId string) (bool, error) {
	// Reuse GetApplicationUserSystemId to check if application user exists
	_, err := client.GetApplicationUserSystemId(ctx, environmentId, applicationId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return false, nil
		}
		// For other errors (like forbidden access), propagate the error
		return false, err
	}

	// If no error, the application user exists
	return true, nil
}

func (client *client) GetApplicationUser(ctx context.Context, environmentId string, applicationId string) (*applicationUserDto, error) {
	users, err := client.getApplicationUsersByApplicationId(ctx, environmentId, applicationId)
	if err != nil {
		return nil, err
	}

	active := make([]applicationUserDto, 0, len(users))
	for _, user := range users {
		if user.DeletedState == 0 {
			active = append(active, user)
		}
	}

	switch len(active) {
	case 0:
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("application user '%s' not found in environment '%s'", applicationId, environmentId))
	case 1:
		sort.Slice(active[0].SecurityRoles, func(i, j int) bool {
			return active[0].SecurityRoles[i].RoleId < active[0].SecurityRoles[j].RoleId
		})
		return &active[0], nil
	default:
		ids := make([]string, 0, len(active))
		for _, user := range active {
			ids = append(ids, user.SystemUserId)
		}
		sort.Strings(ids)
		return nil, fmt.Errorf(
			"multiple active application users found for application '%s' in environment '%s': %s",
			applicationId,
			environmentId,
			strings.Join(ids, ", "))
	}
}

func (client *client) getApplicationUsersByApplicationId(ctx context.Context, environmentId string, applicationId string) ([]applicationUserDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/systemusers", constants.DATAVERSE_API_VERSION),
	}
	values := url.Values{}
	values.Add("$select", "applicationid,systemuserid,fullname,isdisabled,deletedstate,_businessunitid_value")
	values.Add("$expand", "systemuserroles_association($select=roleid,name,_businessunitid_value)")
	values.Add("$filter", fmt.Sprintf("applicationid eq %s", applicationId))
	apiUrl.RawQuery = values.Encode()

	var response applicationUsersResponseDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound || len(response.Value) == 0 {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("application user '%s' not found in environment '%s'", applicationId, environmentId))
	}

	return response.Value, nil
}

func (client *client) GetApplicationUserSystemId(ctx context.Context, environmentId string, applicationId string) (string, error) {
	user, err := client.GetApplicationUser(ctx, environmentId, applicationId)
	if err != nil {
		return "", err
	}
	return user.SystemUserId, nil
}

func (client *client) CreateScopedApplicationUser(ctx context.Context, environmentId, applicationId, businessUnitId string) (*applicationUserDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/systemusers", constants.DATAVERSE_API_VERSION),
	}

	// "businessunitid" is the ReferencingEntityNavigationPropertyName of the
	// systemuser -> businessunit relationship, verified against
	// EntityDefinitions(LogicalName='systemuser')/ManyToOneRelationships. Not
	// every relationship names its navigation property after the lookup column,
	// so check the metadata before adding another @odata.bind here.
	requestBody := map[string]any{
		"applicationid":             applicationId,
		"accessmode":                4,
		"isdisabled":                false,
		"businessunitid@odata.bind": fmt.Sprintf("/businessunits(%s)", businessUnitId),
	}

	response, err := client.createSystemUser(ctx, apiUrl.String(), requestBody)
	if err != nil {
		if purgeErr := client.purgeDeletedApplicationUsersByApplicationId(ctx, environmentId, applicationId, err); purgeErr != nil {
			return nil, purgeErr
		}

		response, err = client.createSystemUser(ctx, apiUrl.String(), requestBody)
		if err != nil {
			return nil, err
		}
	}
	if err := client.Api.HandleForbiddenResponse(response); err != nil {
		return nil, err
	}

	entityID := response.GetHeader(constants.HEADER_ODATA_ENTITY_ID)
	if entityID == "" {
		return nil, errors.New("no entity record id returned from the odata-entityid header")
	}

	re := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
	match := re.FindAllStringSubmatch(entityID, -1)
	if len(match) == 0 {
		return nil, errors.New("no entity record id returned from the odata-entityid header")
	}

	return client.GetApplicationUserBySystemUserId(ctx, environmentId, match[len(match)-1][0])
}

// createSystemUser posts the application user, replaying the request while Dataverse cannot resolve
// the application id in Entra. A service principal created moments earlier has not necessarily
// replicated yet, and that rejection is a definitive 400, so nothing was committed by the attempt.
func (client *client) createSystemUser(ctx context.Context, apiUrl string, requestBody map[string]any) (*api.Response, error) {
	maxRetries := int(constants.ENTRA_APPLICATION_PROPAGATION_POLL_TIMEOUT / constants.ENTRA_APPLICATION_PROPAGATION_POLL_INTERVAL)

	for retry := 0; ; retry++ {
		response, err := client.Api.Execute(ctx, nil, "POST", apiUrl, nil, requestBody, []int{http.StatusNoContent, http.StatusCreated}, nil)
		if err == nil || !isApplicationNotFoundInEntra(err) || retry >= maxRetries {
			return response, err
		}

		tflog.Debug(ctx, "Dataverse cannot see the application in Entra yet, retrying the application user create")
		if sleepErr := client.Api.SleepWithContext(ctx, constants.ENTRA_APPLICATION_PROPAGATION_POLL_INTERVAL); sleepErr != nil {
			return response, sleepErr
		}
	}
}

// isApplicationNotFoundInEntra reports the definitive 400 Dataverse returns when it cannot resolve
// an application user's application id in Entra, either because it has not replicated yet or
// because the registration is gone.
func isApplicationNotFoundInEntra(err error) bool {
	var httpErr customerrors.UnexpectedHttpStatusCodeError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusBadRequest {
		return false
	}
	return strings.Contains(string(httpErr.Body), constants.DATAVERSE_APPLICATION_NOT_IN_ENTRA_ERROR_CODE)
}

func (client *client) purgeDeletedApplicationUsersByApplicationId(ctx context.Context, environmentId, applicationId string, createErr error) error {
	users, err := client.getApplicationUsersByApplicationId(ctx, environmentId, applicationId)
	if err != nil {
		if errors.Is(err, customerrors.ErrObjectNotFound) {
			return createErr
		}

		return fmt.Errorf("%w\n\nAn additional lookup for deleted application users failed: %s", createErr, err.Error())
	}

	purgedAny := false
	for _, user := range users {
		if user.DeletedState == 0 {
			continue
		}

		if err := client.PermanentlyDeleteSystemUser(ctx, environmentId, user.SystemUserId); err != nil {
			return fmt.Errorf("%w\n\nA conflicting deleted application user '%s' was found but could not be permanently deleted: %s", createErr, user.SystemUserId, err.Error())
		}

		purgedAny = true
	}

	if !purgedAny {
		return createErr
	}

	return nil
}

func (client *client) SetApplicationUserDisabledState(ctx context.Context, environmentId, systemUserId, applicationId string, disabled bool) (*applicationUserDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/systemusers(%s)", constants.DATAVERSE_API_VERSION, systemUserId),
	}

	requestBody := map[string]any{
		"isdisabled":    disabled,
		"applicationid": applicationId,
	}

	resp, err := client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, requestBody, []int{http.StatusNoContent, http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}

	return client.GetApplicationUserBySystemUserId(ctx, environmentId, systemUserId)
}

func (client *client) GetApplicationUserBySystemUserId(ctx context.Context, environmentId, systemUserId string) (*applicationUserDto, error) {
	return client.GetPrincipalBySystemUserId(ctx, environmentId, systemUserId)
}

func (client *client) GetPrincipalBySystemUserId(ctx context.Context, environmentId, systemUserId string) (*applicationUserDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/systemusers(%s)", constants.DATAVERSE_API_VERSION, systemUserId),
	}
	values := url.Values{}
	values.Add("$select", "applicationid,systemuserid,fullname,isdisabled,deletedstate,_businessunitid_value")
	values.Add("$expand", "systemuserroles_association($select=roleid,name,_businessunitid_value)")
	apiUrl.RawQuery = values.Encode()

	var response applicationUserDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("principal not found for system user ID %s", systemUserId))
	}
	if response.DeletedState != 0 {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("principal not found for system user ID %s", systemUserId))
	}

	sort.Slice(response.SecurityRoles, func(i, j int) bool {
		return response.SecurityRoles[i].RoleId < response.SecurityRoles[j].RoleId
	})

	return &response, nil
}

// GetRoleHolder reads a principal that holds security roles, resolving its business unit and its
// currently assigned roles. Users and teams live in different tables, so the holder picks the shape.
func (client *client) GetRoleHolder(ctx context.Context, environmentId string, holder roleHolder) (*roleHolderDto, error) {
	if !holder.isTeam() {
		user, err := client.GetPrincipalBySystemUserId(ctx, environmentId, holder.id())
		if err != nil {
			return nil, err
		}
		return &roleHolderDto{Id: user.SystemUserId, BusinessUnitId: user.BusinessUnitId, SecurityRoles: user.SecurityRoles}, nil
	}

	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   holder.path(constants.DATAVERSE_API_VERSION),
	}
	values := url.Values{}
	values.Add("$select", holder.selectFields())
	values.Add("$expand", fmt.Sprintf("%s($select=roleid,name,_businessunitid_value)", holder.association()))
	apiUrl.RawQuery = values.Encode()

	var response teamDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("principal not found for %s", holder))
	}

	slices.SortFunc(response.SecurityRoles, func(a, b applicationSecurityRoleDto) int {
		return strings.Compare(a.RoleId, b.RoleId)
	})

	// Dataverse only lets owner teams and Microsoft Entra group teams hold security roles; an
	// access team (teamtype 1) cannot, so failing here beats a confusing association error later.
	if response.TeamType == TEAM_TYPE_ACCESS {
		return nil, fmt.Errorf("team %s (%s) is an access team, and access teams cannot hold security roles; use an owner team or a Microsoft Entra group team", response.TeamId, response.Name)
	}

	return &roleHolderDto{Id: response.TeamId, BusinessUnitId: response.BusinessUnitId, SecurityRoles: response.SecurityRoles}, nil
}

func (client *client) GetDataverseSecurityRoles(ctx context.Context, environmentId, businessUnitId string) ([]applicationSecurityRoleDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/roles", constants.DATAVERSE_API_VERSION),
	}
	values := url.Values{}
	values.Add("$select", "roleid,name,_businessunitid_value")
	if businessUnitId != "" {
		values.Add("$filter", fmt.Sprintf("_businessunitid_value eq %s", businessUnitId))
	}
	apiUrl.RawQuery = values.Encode()

	var response applicationSecurityRoleArrayDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return response.Value, nil
}

// GetSecurityRoleById fetches one security role row by its immutable id. A missing role returns
// ErrObjectNotFound with the environment named, so a caller can distinguish a bad id from a
// service fault.
func (client *client) GetSecurityRoleById(ctx context.Context, environmentId, roleId string) (*applicationSecurityRoleDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/roles(%s)", constants.DATAVERSE_API_VERSION, roleId),
	}
	values := url.Values{}
	values.Add("$select", "roleid,name,_businessunitid_value")
	apiUrl.RawQuery = values.Encode()

	var role applicationSecurityRoleDto
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &role)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("security role '%s' not found in environment '%s'", roleId, environmentId))
	}

	return &role, nil
}

func (client *client) ResolveSecurityRoleNames(ctx context.Context, environmentId, businessUnitId string, roleNames []string) ([]applicationSecurityRoleDto, error) {
	allRoles, err := client.GetDataverseSecurityRoles(ctx, environmentId, businessUnitId)
	if err != nil {
		return nil, err
	}

	rolesByName := make(map[string][]applicationSecurityRoleDto)
	for _, role := range allRoles {
		rolesByName[role.Name] = append(rolesByName[role.Name], role)
	}

	resolved := make([]applicationSecurityRoleDto, 0, len(roleNames))
	for _, roleName := range roleNames {
		matches := rolesByName[roleName]
		if len(matches) == 0 {
			return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("security role '%s' not found in business unit '%s'", roleName, businessUnitId))
		}
		if len(matches) > 1 {
			return nil, fmt.Errorf("security role '%s' is ambiguous in business unit '%s'", roleName, businessUnitId)
		}
		resolved = append(resolved, matches[0])
	}

	sort.Slice(resolved, func(i, j int) bool {
		return resolved[i].RoleId < resolved[j].RoleId
	})

	return resolved, nil
}

// associationPostError marks an error raised by the association POST itself or its response, as
// opposed to the environment host resolution that precedes it. Only a marked error can mean the
// association was attempted, so only marked errors may classify as ambiguous outcomes.
type associationPostError struct {
	err error
}

func (e associationPostError) Error() string {
	return e.err.Error()
}

func (e associationPostError) Unwrap() error {
	return e.err
}

func (client *client) AddPrincipalSecurityRoles(ctx context.Context, environmentId string, holder roleHolder, roleIds []string) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		// The association was never attempted, so this must not read as an ambiguous outcome.
		return fmt.Errorf("could not resolve the environment host before associating: %w", err)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   holder.associationPath(constants.DATAVERSE_API_VERSION),
	}

	for _, roleId := range roleIds {
		roleToAssociate := map[string]any{
			"@odata.id": fmt.Sprintf("https://%s/api/data/%s/roles(%s)", environmentHost, constants.DATAVERSE_API_VERSION, roleId),
		}
		// Exactly one POST: replaying an association whose outcome is unknown could mask a
		// concurrent caller's grant, so the caller classifies the failure instead of retrying.
		resp, err := client.Api.ExecuteWithoutRetry(ctx, nil, "POST", apiUrl.String(), nil, roleToAssociate, []int{http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
		if err != nil {
			return associationPostError{err: err}
		}
		if err := client.Api.HandleForbiddenResponse(resp); err != nil {
			return associationPostError{err: err}
		}
		if err := client.Api.HandleNotFoundResponse(resp); err != nil {
			return associationPostError{err: err}
		}
	}

	return nil
}

func (client *client) RemovePrincipalSecurityRoles(ctx context.Context, environmentId string, holder roleHolder, roleIds []string) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	for _, roleId := range roleIds {
		apiUrl := &url.URL{
			Scheme: constants.HTTPS,
			Host:   environmentHost,
			Path:   holder.associationPath(constants.DATAVERSE_API_VERSION),
		}
		values := url.Values{}
		values.Add("$id", fmt.Sprintf("https://%s/api/data/%s/roles(%s)", environmentHost, constants.DATAVERSE_API_VERSION, roleId))
		apiUrl.RawQuery = values.Encode()

		resp, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
		if err != nil {
			// Dataverse rejects the dissociation when the principal's application registration no
			// longer exists in Entra. The grant cannot outlive the identity it was made to, so this
			// is the outcome a removal wants rather than a failure to report.
			if isApplicationNotFoundInEntra(err) {
				tflog.Debug(ctx, fmt.Sprintf("Security role %s cannot be dissociated from %s because its application is no longer in Entra; treating it as removed", roleId, holder))
				continue
			}
			return err
		}
		if err := client.Api.HandleForbiddenResponse(resp); err != nil {
			return err
		}
		// A 404 here means the association, the role or the principal is already gone, which is the
		// outcome a removal wants; failing would make an out-of-band removal break destroy.
		if resp.HttpResponse.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, fmt.Sprintf("Security role %s was already dissociated from %s", roleId, holder))
			continue
		}
	}

	return nil
}

func (client *client) DeactivateSystemUser(ctx context.Context, environmentId string, systemUserId string) error {
	// Get the application user to find the application ID
	appUser, err := client.getApplicationUserBySystemId(ctx, environmentId, systemUserId)
	if err != nil {
		return err
	}

	if _, err = client.SetApplicationUserDisabledState(ctx, environmentId, systemUserId, appUser.ApplicationId, true); err != nil {
		return err
	}

	return nil
}

func (client *client) DeleteSystemUser(ctx context.Context, environmentId string, systemUserId string) error {
	return client.executeSystemUserDelete(ctx, environmentId, systemUserId, []int{http.StatusNoContent, http.StatusOK})
}

func (client *client) PermanentlyDeleteSystemUser(ctx context.Context, environmentId string, systemUserId string) error {
	if err := client.executeSystemUserDelete(ctx, environmentId, systemUserId, []int{http.StatusNoContent, http.StatusOK, http.StatusNotFound}); err != nil {
		return err
	}

	return client.executeSystemUserDelete(ctx, environmentId, systemUserId, []int{http.StatusNoContent, http.StatusOK, http.StatusNotFound})
}

func (client *client) executeSystemUserDelete(ctx context.Context, environmentId string, systemUserId string, expectedStatusCodes []int) error {
	// Get the environment host
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	// Create the Dataverse Web API URL to delete the system user
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/systemusers(%s)", constants.DATAVERSE_API_VERSION, systemUserId),
	}

	// Make the request
	_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, expectedStatusCodes, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *client) getEnvironment(ctx context.Context, environmentId string) (*environmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	env := environmentIdDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &env)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("environment %s not found", environmentId))
	}

	return &env, nil
}

func (client *client) GetTenantApplications(ctx context.Context) ([]tenantApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/appmanagement/applicationPackages",
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.APPLICATION_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	application := tenantApplicationArrayDto{}

	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application.Value, nil
}

func (client *client) GetApplicationsByEnvironmentId(ctx context.Context, environmentId string) ([]environmentApplicationDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages", environmentId),
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.APPLICATION_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	application := environmentApplicationArrayDto{}

	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &application)
	if err != nil {
		return nil, err
	}

	return application.Value, nil
}

func (client *client) InstallApplicationInEnvironment(ctx context.Context, environmentId string, uniqueName string) (string, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/appmanagement/environments/%s/applicationPackages/%s/install", environmentId, uniqueName),
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.APPLICATION_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	response, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusAccepted}, nil)
	if err != nil {
		return "", err
	}

	applicationId := ""
	if response.HttpResponse.StatusCode == http.StatusAccepted {
		operationLocationHeader := response.GetHeader(constants.HEADER_OPERATION_LOCATION)
		if operationLocationHeader == "" {
			tflog.Error(ctx, "Missing operation location header in response")
			return "", errors.New("missing operation location header in response")
		}
		tflog.Debug(ctx, "Opeartion Location Header: "+operationLocationHeader)

		_, err = url.Parse(operationLocationHeader)
		if err != nil {
			tflog.Error(ctx, "Error parsing location header: "+err.Error())
			return "", err
		}

		for {
			lifecycleResponse := environmentApplicationLifecycleDto{}
			response, err := client.Api.Execute(ctx, nil, "GET", operationLocationHeader, nil, nil, []int{http.StatusOK, http.StatusConflict}, &lifecycleResponse)
			if err != nil {
				return "", err
			}

			if response.HttpResponse.StatusCode == http.StatusConflict {
				tflog.Debug(ctx, "Lifecycle Operation HTTP Status: '"+response.HttpResponse.Status+"'")
				continue
			}

			if lifecycleResponse.Status == "Succeeded" {
				parts := strings.Split(lifecycleResponse.CreatedDateTime, "/")
				if len(parts) == 0 {
					return "", errors.New("can't parse application id from response " + lifecycleResponse.CreatedDateTime)
				}
				applicationId = parts[len(parts)-1]
				tflog.Debug(ctx, "Created Application Id: "+applicationId)
				break
			} else if lifecycleResponse.Status == "Failed" {
				return "", errors.New("application installation failed. status message: " + lifecycleResponse.Error.Message)
			}
		}
	} else if response.HttpResponse.StatusCode == http.StatusCreated {
		appCreatedResponse := environmentApplicationLifecycleCreatedDto{}
		err = response.MarshallTo(&appCreatedResponse)
		if err != nil {
			return "", err
		}
		if appCreatedResponse.Properties.ProvisioningState != "Succeeded" {
			return "", errors.New("application installation failed. provisioning state: " + appCreatedResponse.Properties.ProvisioningState)
		}
		applicationId = appCreatedResponse.Name
	}

	return applicationId, nil
}

func (client *client) getApplicationUserBySystemId(ctx context.Context, environmentId string, systemUserId string) (*applicationUserDto, error) {
	return client.GetPrincipalBySystemUserId(ctx, environmentId, systemUserId)
}
