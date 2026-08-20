// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
)

func NewManagedSolutionClient(apiClient *api.Client) Client {
	return Client{
		Api:            apiClient,
		SolutionClient: solution.NewSolutionClient(apiClient),
	}
}

type Client struct {
	Api            *api.Client
	SolutionClient solution.Client
}

func (client *Client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
	return client.SolutionClient.DataverseExists(ctx, environmentId)
}

func (client *Client) GetSolutionUniqueName(ctx context.Context, environmentId, name string) (*solution.SolutionDto, error) {
	return client.SolutionClient.GetSolutionUniqueName(ctx, environmentId, name)
}

func (client *Client) GetSolutionById(ctx context.Context, environmentId, solutionId string) (*solution.SolutionDto, error) {
	return client.SolutionClient.GetSolutionById(ctx, environmentId, solutionId)
}

func (client *Client) GetSolutions(ctx context.Context, environmentId string) ([]solution.SolutionDto, error) {
	return client.SolutionClient.GetSolutions(ctx, environmentId)
}

func (client *Client) GetInstalledSolutions(ctx context.Context, environmentId string) ([]solution.SolutionDto, error) {
	var solutions []solution.SolutionDto
	err := client.retryWhileDataverseRejectsCaller(ctx, environmentId, func() error {
		var readErr error
		solutions, readErr = client.SolutionClient.GetInstalledSolutions(ctx, environmentId)
		return readErr
	})
	if err != nil {
		return nil, err
	}
	return solutions, nil
}

// retryWhileDataverseRejectsCaller replays a read while Dataverse answers with 403, which a freshly
// provisioned organization does until it has provisioned the caller as an application user.
func (client *Client) retryWhileDataverseRejectsCaller(ctx context.Context, environmentId string, read func() error) error {
	deadline := time.Now().Add(constants.DATAVERSE_CALLER_PROVISIONING_POLL_TIMEOUT)
	for {
		err := read()
		var accessDenied customerrors.AccessDeniedError
		if err == nil || !errors.As(err, &accessDenied) || time.Now().After(deadline) {
			return err
		}

		tflog.Debug(ctx, fmt.Sprintf("Dataverse in environment '%s' denied access to the caller, retrying while the environment finishes provisioning", environmentId))
		if sleepErr := client.Api.SleepWithContext(ctx, constants.DATAVERSE_CALLER_PROVISIONING_POLL_INTERVAL); sleepErr != nil {
			return sleepErr
		}
	}
}

func (client *Client) DeleteSolution(ctx context.Context, environmentId, solutionId string) error {
	return client.SolutionClient.DeleteSolution(ctx, environmentId, solutionId)
}

// importOperation selects the Dataverse import API used to realize the managed solution:
// initial installation or in-place staged upgrade of an existing managed install.
type importOperation string

const (
	importOperationInstall         importOperation = "ImportSolutionAsync"
	importOperationStageAndUpgrade importOperation = "StageAndUpgradeAsync"
)

func (client *Client) ApplyManagedSolution(ctx context.Context, environmentId string, content []byte, componentParameters []any, operation importOperation, skipProductUpdateDependencies bool) (*solution.SolutionDto, error) {
	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, errors.New("solution content is nil")
	}

	stageRequestBody := stageSolutionImportDto{
		CustomizationFile: base64.StdEncoding.EncodeToString(content),
	}

	stageURL := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/StageSolution",
	}

	stageResponse := stageSolutionImportResponseDto{}
	resp, err := client.Api.Execute(ctx, nil, "POST", stageURL.String(), nil, stageRequestBody, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &stageResponse)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	if stageResponse.StageSolutionResults.StageSolutionStatus != "Passed" {
		e := fmt.Errorf("solution failed with status: '%s'", stageResponse.StageSolutionResults.StageSolutionStatus)
		for _, missingDependency := range stageResponse.StageSolutionResults.MissingDependencies {
			e = errors.Join(fmt.Errorf("missing dependency: '%s'", missingDependency.RequiredComponentSchemaName), e)
		}
		for _, validation := range stageResponse.StageSolutionResults.SolutionValidationResults {
			e = errors.Join(fmt.Errorf("solution validation failed: %s", validation.Message), e)
		}
		return nil, e
	}

	importRequestBody := importSolutionDto{
		PublishWorkflows:                 true,
		OverwriteUnmanagedCustomizations: true,
		SkipProductUpdateDependencies:    skipProductUpdateDependencies,
		ComponentParameters:              componentParameters,
		SolutionParameters: importSolutionParametersDto{
			StageSolutionUploadId: stageResponse.StageSolutionResults.StageSolutionUploadId,
		},
	}

	importURL := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/" + string(operation),
	}
	importResponse := importSolutionResponseDto{}
	resp, err = client.Api.ExecuteWithoutRetry(ctx, nil, "POST", importURL.String(), nil, importRequestBody, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &importResponse)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
		return nil, err
	}

	pollURL := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.2/asyncoperations(%s)", importResponse.AsyncOperationId),
	}

	for {
		asyncResponse := asyncSolutionPullResponseDto{}
		resp, err = client.Api.Execute(ctx, nil, "GET", pollURL.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &asyncResponse)
		if err != nil {
			return nil, err
		}
		if err := client.Api.HandleForbiddenResponse(resp); err != nil {
			return nil, err
		}
		if err := client.Api.HandleNotFoundResponse(resp); err != nil {
			return nil, err
		}

		if asyncResponse.CompletedOn != "" {
			if err := client.validateSolutionImportResult(ctx, environmentHost, importResponse.ImportJobKey); err != nil {
				return nil, err
			}
			return client.SolutionClient.GetSolutionUniqueName(ctx, environmentId, stageResponse.StageSolutionResults.SolutionDetails.SolutionUniqueName)
		}

		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return nil, err
		}
	}
}

func (client *Client) PublishAllCustomizations(ctx context.Context, environmentId string) error {
	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	publishURL := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/PublishAllXml",
	}
	resp, err := client.Api.Execute(ctx, nil, "POST", publishURL.String(), nil, struct{}{}, []int{http.StatusOK, http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
	if err != nil {
		return err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return err
	}
	return client.Api.HandleNotFoundResponse(resp)
}

func (client *Client) validateSolutionImportResult(ctx context.Context, environmentHost, importJobKey string) error {
	apiURL := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.0/RetrieveSolutionImportResult(ImportJobId=%s)", importJobKey),
	}

	response := validateSolutionImportResponseDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiURL.String(), nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return err
	}
	if response.SolutionOperationResult.Status != "Passed" {
		messages := make([]string, 0, len(response.SolutionOperationResult.ErrorMessages))
		for _, message := range response.SolutionOperationResult.ErrorMessages {
			messages = append(messages, fmt.Sprint(message))
		}
		return fmt.Errorf("solution import failed: %s", strings.Join(messages, "; "))
	}
	return nil
}

func (client *Client) getConnections(ctx context.Context, environmentId string) ([]connectionDto, error) {
	apiURL := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildEnvironmentHostUri(environmentId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   "/connectivity/connections",
	}
	values := url.Values{}
	values.Add("api-version", "1")
	apiURL.RawQuery = values.Encode()

	connections := connectionArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiURL.String(), nil, nil, []int{http.StatusOK}, &connections)
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}

	return connections.Value, nil
}

func (client *Client) GetUnmanagedEnvironmentVariableDefinitions(ctx context.Context, environmentId string, schemaNames []string) (map[string]bool, error) {
	result := make(map[string]bool, len(schemaNames))

	for _, schemaName := range schemaNames {
		values := url.Values{}
		values.Add("$select", "schemaname,ismanaged")
		values.Add("$filter", fmt.Sprintf("schemaname eq '%s'", escapeODataString(schemaName)))

		response := environmentVariableDefinitionsResponseDto{}
		if err := client.retryWhileDataverseRejectsCaller(ctx, environmentId, func() error {
			return client.SolutionClient.GetTableData(ctx, environmentId, "environmentvariabledefinitions", values.Encode(), &response)
		}); err != nil {
			return nil, err
		}

		for _, definition := range response.Value {
			if !definition.IsManaged {
				result[schemaName] = true
			}
		}
	}

	return result, nil
}

func inspectSolutionPackage(zipPath string) (*solutionPackage, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	pkg := &solutionPackage{
		ConnectionReferences: make(map[string]packageConnectionReference),
		EnvironmentVariables: make(map[string]packageEnvironmentVariable),
		Dependencies:         make(map[string]string),
	}

	for _, file := range reader.File {
		switch {
		case file.Name == "solution.xml":
			if err := processSolutionXMLFile(file, pkg); err != nil {
				return nil, err
			}

		case file.Name == "customizations.xml":
			if err := processCustomizationsXMLFile(file, pkg); err != nil {
				return nil, err
			}

		case strings.HasPrefix(file.Name, "environmentvariabledefinitions/") && strings.HasSuffix(file.Name, "/environmentvariabledefinition.xml"):
			if err := processEnvironmentVariableDefinitionFile(file, pkg); err != nil {
				return nil, err
			}

		case strings.HasPrefix(file.Name, "environmentvariabledefinitions/") && strings.HasSuffix(file.Name, "/environmentvariablevalues.json"):
			if err := processEnvironmentVariableValueFile(file, pkg); err != nil {
				return nil, err
			}

		default:
			continue
		}
	}

	if pkg.UniqueName == "" {
		return nil, errors.New("solution package did not contain solution.xml with a unique name")
	}
	if pkg.Version == "" {
		return nil, errors.New("solution package did not contain solution.xml with a version")
	}

	return pkg, nil
}

func processSolutionXMLFile(file *zip.File, pkg *solutionPackage) error {
	content, err := readZipFile(file)
	if err != nil {
		return err
	}

	solutionXML := solutionPackageXML{}
	if err := xml.Unmarshal(content, &solutionXML); err != nil {
		return err
	}

	pkg.UniqueName = solutionXML.SolutionManifest.UniqueName
	pkg.Version = solutionXML.SolutionManifest.Version
	pkg.IsManaged = solutionXML.SolutionManifest.Managed == 1
	if len(solutionXML.SolutionManifest.LocalizedNames) > 0 {
		pkg.DisplayName = solutionXML.SolutionManifest.LocalizedNames[0].Description
	}

	return collectDependencies(solutionXML.SolutionManifest.MissingDependencies.MissingDependencies, pkg)
}

func collectDependencies(missingDependencies []solutionMissingDependencyXML, pkg *solutionPackage) error {
	for _, missingDependency := range missingDependencies {
		for _, required := range missingDependency.Required {
			solutionRef := strings.TrimSpace(required.Solution)
			if solutionRef == "" {
				return fmt.Errorf("invalid solution dependency reference for schema %q: solution attribute is empty", required.SchemaName)
			}
			if strings.EqualFold(solutionRef, "Active") {
				return fmt.Errorf("invalid solution dependency reference for schema %q: solution attribute is %q", required.SchemaName, required.Solution)
			}

			name, version, err := parseSolutionRef(solutionRef)
			if err != nil {
				return err
			}

			if err := setHighestDependencyVersion(pkg.Dependencies, name, version); err != nil {
				return err
			}
		}
	}

	return nil
}

func setHighestDependencyVersion(dependencies map[string]string, name, version string) error {
	existingVersion, exists := dependencies[name]
	if !exists {
		dependencies[name] = version
		return nil
	}

	cmp, err := compareVersionStrings(version, existingVersion)
	if err != nil {
		return err
	}
	if cmp > 0 {
		dependencies[name] = version
	}

	return nil
}

func processCustomizationsXMLFile(file *zip.File, pkg *solutionPackage) error {
	content, err := readZipFile(file)
	if err != nil {
		return err
	}

	customizationsXML := solutionCustomizationsXML{}
	if err := xml.Unmarshal(content, &customizationsXML); err != nil {
		return err
	}

	for _, reference := range customizationsXML.ConnectionReferences {
		pkg.ConnectionReferences[reference.LogicalName] = packageConnectionReference(reference)
	}

	return nil
}

func processEnvironmentVariableDefinitionFile(file *zip.File, pkg *solutionPackage) error {
	content, err := readZipFile(file)
	if err != nil {
		return err
	}

	definitionXML := environmentVariableDefinitionXML{}
	if err := xml.Unmarshal(content, &definitionXML); err != nil {
		return err
	}

	existing := pkg.EnvironmentVariables[definitionXML.SchemaName]
	existing.SchemaName = definitionXML.SchemaName
	existing.HasDefaultValue = definitionXML.DefaultValue != nil
	pkg.EnvironmentVariables[definitionXML.SchemaName] = existing

	return nil
}

func processEnvironmentVariableValueFile(file *zip.File, pkg *solutionPackage) error {
	content, err := readZipFile(file)
	if err != nil {
		return err
	}

	valueJSON := environmentVariableValuesJSON{}
	if err := json.Unmarshal(content, &valueJSON); err != nil {
		return err
	}

	schemaName := filepath.Base(filepath.Dir(file.Name))
	value := pkg.EnvironmentVariables[schemaName]
	value.SchemaName = schemaName
	if valueJSON.EnvironmentVariableValues.EnvironmentVariableValue != nil {
		value.ContainsPackagedValue = true
	}
	pkg.EnvironmentVariables[schemaName] = value

	return nil
}

func readZipFile(file *zip.File) ([]byte, error) {
	handle, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	return io.ReadAll(handle)
}

func parseSolutionRef(input string) (name string, version string, err error) {
	start := strings.LastIndex(input, "(")
	end := strings.LastIndex(input, ")")

	if start == -1 || end == -1 || end < start {
		return "", "", fmt.Errorf("invalid solution reference: %s", input)
	}

	name = strings.TrimSpace(input[:start])
	version = strings.TrimSpace(input[start+1 : end])

	return name, version, nil
}

// compareVersionStringsIgnoringBuild compares only major.minor.patch, deliberately
// tolerating build-segment drift: dependency declarations captured from a package's
// MissingDependencies frequently pin a build number that floats independently of the
// installed dependency's build (e.g. PCF support solutions), and requiring an exact
// or higher build would fail deployments that are actually compatible.
func compareVersionStringsIgnoringBuild(left, right string) (int, error) {
	leftParts, err := normalizeVersionParts(left)
	if err != nil {
		return 0, err
	}
	rightParts, err := normalizeVersionParts(right)
	if err != nil {
		return 0, err
	}

	for len(leftParts) < 3 {
		leftParts = append(leftParts, 0)
	}
	for len(rightParts) < 3 {
		rightParts = append(rightParts, 0)
	}

	for i := 0; i < 3; i++ {
		if leftParts[i] < rightParts[i] {
			return -1, nil
		}
		if leftParts[i] > rightParts[i] {
			return 1, nil
		}
	}

	return 0, nil
}

func compareVersionStrings(left, right string) (int, error) {
	leftParts, err := normalizeVersionParts(left)
	if err != nil {
		return 0, err
	}
	rightParts, err := normalizeVersionParts(right)
	if err != nil {
		return 0, err
	}

	maxLen := len(leftParts)
	if len(rightParts) > maxLen {
		maxLen = len(rightParts)
	}

	for len(leftParts) < maxLen {
		leftParts = append(leftParts, 0)
	}
	for len(rightParts) < maxLen {
		rightParts = append(rightParts, 0)
	}

	for i := 0; i < maxLen; i++ {
		if leftParts[i] < rightParts[i] {
			return -1, nil
		}
		if leftParts[i] > rightParts[i] {
			return 1, nil
		}
	}

	return 0, nil
}

func normalizeVersionParts(raw string) ([]int, error) {
	if raw == "" {
		return nil, errors.New("version cannot be empty")
	}

	parts := strings.Split(raw, ".")
	normalized := make([]int, 0, len(parts))
	for _, part := range parts {
		value, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", raw, err)
		}
		normalized = append(normalized, value)
	}

	return normalized, nil
}

func buildConnectionParameters(packageRefs map[string]packageConnectionReference, configuredRefs map[string]string, connections []connectionDto) ([]any, error) {
	connectionIndex := make(map[string]connectionDto, len(connections))
	for _, conn := range connections {
		connectionIndex[conn.Name] = conn
	}

	parameters := make([]any, 0, len(packageRefs))
	for logicalName, packageRef := range packageRefs {
		connectionID, exists := configuredRefs[logicalName]
		if !exists {
			return nil, fmt.Errorf("connection reference %q is required by the solution package but was not configured", logicalName)
		}

		connection, exists := connectionIndex[connectionID]
		if !exists {
			return nil, fmt.Errorf("connection reference %q points to unknown connection id %q", logicalName, connectionID)
		}

		if !strings.EqualFold(connection.Properties.ApiId, packageRef.ConnectorID) {
			return nil, fmt.Errorf("connection reference %q expects connector %q but connection %q resolves to %q", logicalName, packageRef.ConnectorID, connectionID, connection.Properties.ApiId)
		}

		parameters = append(parameters, importSolutionConnectionReferenceDto{
			Type:                           "Microsoft.Dynamics.CRM.connectionreference",
			ConnectionReferenceLogicalName: logicalName,
			ConnectionReferenceDisplayName: "",
			Description:                    "",
			ConnectorId:                    connection.Properties.ApiId,
			ConnectionId:                   connectionID,
		})
	}

	if len(parameters) == 0 {
		return nil, nil
	}

	return parameters, nil
}

func validateConnectionReferences(packageRefs map[string]packageConnectionReference, configuredRefs map[string]string) error {
	packageKeys := make([]string, 0, len(packageRefs))
	for key := range packageRefs {
		packageKeys = append(packageKeys, key)
	}
	slices.Sort(packageKeys)

	configuredKeys := make([]string, 0, len(configuredRefs))
	for key := range configuredRefs {
		configuredKeys = append(configuredKeys, key)
	}
	slices.Sort(configuredKeys)

	missing := make([]string, 0)
	for _, key := range packageKeys {
		if _, exists := configuredRefs[key]; !exists {
			missing = append(missing, key)
		}
	}

	unexpected := make([]string, 0)
	for _, key := range configuredKeys {
		if _, exists := packageRefs[key]; !exists {
			unexpected = append(unexpected, key)
		}
	}

	if len(missing) == 0 && len(unexpected) == 0 {
		return nil
	}

	messages := make([]string, 0, 2)
	if len(missing) > 0 {
		messages = append(messages, fmt.Sprintf("missing required connection references: %s", strings.Join(missing, ", ")))
	}
	if len(unexpected) > 0 {
		messages = append(messages, fmt.Sprintf("unexpected connection references: %s", strings.Join(unexpected, ", ")))
	}
	return fmt.Errorf("connection reference validation failed: %s", strings.Join(messages, "; "))
}

func validateEnvironmentVariablePackaging(packageVars map[string]packageEnvironmentVariable, existingUnmanaged map[string]bool) error {
	problems := make([]string, 0)
	keys := make([]string, 0, len(packageVars))
	for key := range packageVars {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		if packageVars[key].ContainsPackagedValue {
			problems = append(problems, fmt.Sprintf("%s contains a packaged current value; environment variable values are owned by the deployment context, never the solution", key))
			continue
		}
		if existingUnmanaged[key] {
			problems = append(problems, fmt.Sprintf("%s already exists as an unmanaged definition in the target environment; importing this package would silently absorb it into the solution's managed layer and uninstalling would delete it. Make the package reference-only for this variable, or deliberately remove the environment-side definition first", key))
		}
	}

	if len(problems) == 0 {
		return nil
	}

	return fmt.Errorf("environment variable packaging validation failed: %s", strings.Join(problems, "; "))
}

func validateDependencies(required map[string]string, installed []solution.SolutionDto) error {
	installedByName := make(map[string]solution.SolutionDto, len(installed))
	for _, item := range installed {
		installedByName[strings.ToLower(item.Name)] = item
	}

	problems := make([]string, 0)
	keys := make([]string, 0, len(required))
	for key := range required {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, dependencyName := range keys {
		if isBuiltInSolutionDependency(dependencyName) {
			continue
		}

		requiredVersion := required[dependencyName]
		installedSolution, exists := installedByName[strings.ToLower(dependencyName)]
		if !exists {
			problems = append(problems, fmt.Sprintf("required dependency %q >= %s is not installed", dependencyName, requiredVersion))
			continue
		}

		cmp, err := compareVersionStringsIgnoringBuild(installedSolution.Version, requiredVersion)
		if err != nil {
			return err
		}
		if cmp < 0 {
			problems = append(problems, fmt.Sprintf("required dependency %q >= %s but %s is installed", dependencyName, requiredVersion, installedSolution.Version))
		}
	}

	if len(problems) == 0 {
		return nil
	}

	return fmt.Errorf("dependency validation failed: %s", strings.Join(problems, "; "))
}

func resolveSourceToPath(ctx context.Context, source *SourceModel) (string, func(), error) {
	if source == nil {
		return "", nil, errors.New("source is required")
	}

	if !source.Path.IsNull() && !source.Path.IsUnknown() && source.Path.ValueString() != "" {
		return source.Path.ValueString(), func() {}, nil
	}

	if source.URL.IsNull() || source.URL.IsUnknown() || source.URL.ValueString() == "" {
		return "", nil, errors.New("source.path or source.url must be set")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, source.URL.ValueString(), nil)
	if err != nil {
		return "", nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("failed to download solution package from %s: unexpected status %d", source.URL.ValueString(), response.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "managed-solution-*.zip")
	if err != nil {
		return "", nil, err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, response.Body); err != nil {
		_ = os.Remove(tempFile.Name())
		return "", nil, err
	}

	return tempFile.Name(), func() {
		_ = os.Remove(tempFile.Name())
	}, nil
}

func readSourceContent(sourcePath string) ([]byte, error) {
	return os.ReadFile(sourcePath)
}

func escapeODataString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func cloneSourceModel(source *SourceModel) *SourceModel {
	if source == nil {
		return nil
	}
	return &SourceModel{
		Path: source.Path,
		URL:  source.URL,
	}
}
