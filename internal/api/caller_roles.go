// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

// IsCallerWithoutRolesMessage reports whether a Dataverse error or validation message says that the
// calling user has not been assigned any security roles, such as "The user with id ... has not been
// assigned any roles. They need a role with the prvImportCustomization privilege".
func IsCallerWithoutRolesMessage(message string) bool {
	return strings.Contains(strings.ToLower(message), "has not been assigned any roles")
}

// ReplayWhileNewDataverseOrgAssignsCallerRoles runs operation, a request to the Dataverse organization
// serving url, and runs it again while it reports that it was refused because its caller has not been
// assigned any security roles yet. A new organization does that until it has assigned the caller its
// roles, so only an organization this provider process created less than constants.NEW_DATAVERSE_ORG_PERIOD
// ago gets the replays, spaced by constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_INTERVAL, until their
// delays add up to constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_BUDGET. Any other refusal or error is
// returned as it is. operation must be safe to repeat, and it returns callerWithoutRoles together with
// its error when the error is that refusal.
func (client *Client) ReplayWhileNewDataverseOrgAssignsCallerRoles(ctx context.Context, url string, operation func() (callerWithoutRoles bool, err error)) error {
	newOrgAge, newOrg := newDataverseOrgAge(url)
	replays := 0
	var waited time.Duration
	var windowStart time.Time
	for {
		callerWithoutRoles, err := operation()
		if err == nil || !callerWithoutRoles {
			if replays > 0 {
				outcome := "was accepted"
				if err != nil {
					outcome = "failed for another reason"
				}
				logCallerRolesWindow(ctx, url, windowStart, replays, newOrgAge, outcome)
			}
			return err
		}
		if !newOrg {
			return err
		}
		if windowStart.IsZero() {
			windowStart = time.Now()
		}
		if waited >= constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_BUDGET {
			logCallerRolesWindow(ctx, url, windowStart, replays, newOrgAge, "was given up while the caller had no roles")
			return fmt.Errorf("the caller still had no security roles after %d retries over %s, the most the provider waits for a Dataverse organization it created %s ago to assign them: %w",
				replays, waited, newOrgAge.Round(time.Second), err)
		}

		delay := constants.NEW_DATAVERSE_ORG_CALLER_ROLES_RETRY_INTERVAL
		tflog.Debug(ctx, fmt.Sprintf("Dataverse organization created %s ago has not assigned the caller any roles yet, retrying %s after %s", newOrgAge.Round(time.Second), urlWithoutQuery(url), delay))
		if sleepErr := client.SleepWithContext(ctx, delay); sleepErr != nil {
			logCallerRolesWindow(ctx, url, windowStart, replays, newOrgAge, "was interrupted while the caller had no roles: "+sleepErr.Error())
			return sleepErr
		}
		waited += delay
		replays++
	}
}

// logCallerRolesWindow reports how long a new organization refused an operation because its caller had
// no roles, so acceptance runs, which log at WARN, measure how long the role assignment takes.
func logCallerRolesWindow(ctx context.Context, url string, windowStart time.Time, replays int, newOrgAge time.Duration, outcome string) {
	window := time.Since(windowStart).Round(time.Millisecond).String()
	tflog.Warn(ctx, fmt.Sprintf("Operation refused because the caller had no security roles for %s over %d replays, then it %s", window, replays, outcome), map[string]any{
		"caller_without_roles_window":  window,
		"caller_without_roles_replays": replays,
		"caller_without_roles_outcome": outcome,
		"dataverse_org_age":            newOrgAge.Round(time.Second).String(),
		"request_url_no_query":         urlWithoutQuery(url),
	})
}
