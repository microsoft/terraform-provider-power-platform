// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package common

var (
	// ProviderVersion is the version of the released provider, set during build/release process with ldflags.
	// This value can not be const as it is set after the build during the linking process.
	ProviderVersion = "0.0.0-dev" // Default value for development builds
	Commit          = "dev"       // Default value for development builds
	Branch          = "dev"       // Default value for development builds
)
