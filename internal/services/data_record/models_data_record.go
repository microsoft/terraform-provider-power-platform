// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package data_record

type dataRecordDto struct {
	Id           string `json:"id"`
	OdataContext string `json:"@odata.context"`
	OdataEtag    string `json:"@odata.etag"`
}
