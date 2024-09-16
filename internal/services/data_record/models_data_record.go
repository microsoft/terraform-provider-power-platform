// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package data_record

type DataRecordDto struct {
	Id           string `json:"id"`
	OdataContext string `json:"@odata.context"`
	OdataEtag    string `json:"@odata.etag"`
}
