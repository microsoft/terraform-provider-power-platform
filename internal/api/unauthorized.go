// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package api

import (
	"context"
	"fmt"
	"net/http"
	neturl "net/url"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

// newDataverseOrgs remembers when this provider process created each Dataverse organization, keyed by
// host. It is process wide because Terraform configures a new provider client for every command, so
// the client that created an environment is not the one that later reads its records.
var newDataverseOrgs = struct {
	sync.Mutex
	createdAt map[string]time.Time
}{createdAt: map[string]time.Time{}}

// RecordNewDataverseOrg notes that the Dataverse organization at instanceUrl has just been created.
// For constants.NEW_DATAVERSE_ORG_PERIOD afterwards, a 401 from that organization is replayed for up
// to constants.NEW_DATAVERSE_ORG_UNAUTHORIZED_RETRY_BUDGET instead of constants.MAX_UNAUTHORIZED_RETRIES
// times, because a new organization refuses its caller for a while before it accepts it.
func RecordNewDataverseOrg(instanceUrl string) {
	host := dataverseOrgHost(instanceUrl)
	if host == "" {
		return
	}

	now := time.Now()
	newDataverseOrgs.Lock()
	defer newDataverseOrgs.Unlock()
	for knownHost, createdAt := range newDataverseOrgs.createdAt {
		if now.Sub(createdAt) > constants.NEW_DATAVERSE_ORG_PERIOD {
			delete(newDataverseOrgs.createdAt, knownHost)
		}
	}
	newDataverseOrgs.createdAt[host] = now
}

// IsNewDataverseOrg reports whether the organization serving url was created by this provider process
// less than constants.NEW_DATAVERSE_ORG_PERIOD ago.
func IsNewDataverseOrg(url string) bool {
	_, ok := newDataverseOrgAge(url)
	return ok
}

// newDataverseOrgAge returns how long ago this provider process created the organization serving url,
// and false when it did not create it or created it more than constants.NEW_DATAVERSE_ORG_PERIOD ago.
func newDataverseOrgAge(url string) (time.Duration, bool) {
	host := dataverseOrgHost(url)
	if host == "" {
		return 0, false
	}

	newDataverseOrgs.Lock()
	createdAt, ok := newDataverseOrgs.createdAt[host]
	newDataverseOrgs.Unlock()
	if !ok {
		return 0, false
	}

	age := time.Since(createdAt)
	return age, age <= constants.NEW_DATAVERSE_ORG_PERIOD
}

func dataverseOrgHost(url string) string {
	u, err := neturl.Parse(url)
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Host)
}

// unauthorizedReplays tracks the 401 answers that one Execute call replays.
type unauthorizedReplays struct {
	// count and waited are the replays after a 401 and the delay slept before them.
	count  int
	waited time.Duration
	// newOrg is set when the request targets a Dataverse organization this process created newOrgAge ago.
	newOrg    bool
	newOrgAge time.Duration
	// windowStart is when the current run of 401 answers began, zero when the last answer was not 401.
	windowStart   time.Time
	windowReplays int
}

func (u *unauthorizedReplays) observe(ctx context.Context, method, url string, statusCode int) {
	if statusCode != http.StatusUnauthorized {
		u.endWindow(ctx, method, url, fmt.Sprintf("was answered with HTTP %d", statusCode))
		return
	}
	if u.windowStart.IsZero() {
		u.windowStart = time.Now()
		u.windowReplays = 0
	}
}

// recordRetry notes that an answer with statusCode is replayed after delay, and reports whether the
// replay counts toward constants.MAX_RETRY_COUNT. The 401 replays of a new organization do not, because
// their own time budget bounds them.
func (u *unauthorizedReplays) recordRetry(statusCode int, delay time.Duration) bool {
	if statusCode != http.StatusUnauthorized {
		return true
	}
	u.count++
	u.waited += delay
	u.windowReplays++
	return !u.newOrg
}

// endWindow reports how long the current run of 401 answers lasted and how it ended, so acceptance runs
// measure how long a new organization keeps refusing its caller. It logs a warning because the nightly
// acceptance run logs at WARN, and it fires only for a request that was refused at least once.
func (u *unauthorizedReplays) endWindow(ctx context.Context, method, url, outcome string) {
	if u.windowStart.IsZero() {
		return
	}

	fields := map[string]any{
		"unauthorized_window":  time.Since(u.windowStart).Round(time.Millisecond).String(),
		"unauthorized_replays": u.windowReplays,
		"new_dataverse_org":    u.newOrg,
		"request_method":       method,
		"request_url_no_query": urlWithoutQuery(url),
		"unauthorized_outcome": outcome,
	}
	if u.newOrg {
		fields["dataverse_org_age"] = u.newOrgAge.Round(time.Second).String()
	}
	tflog.Warn(ctx, fmt.Sprintf("Request refused with 401 for %s over %d replays, then it %s", fields["unauthorized_window"], u.windowReplays, outcome), fields)

	u.windowStart = time.Time{}
	u.windowReplays = 0
}

func urlWithoutQuery(url string) string {
	u, err := neturl.Parse(url)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host + u.Path
}
