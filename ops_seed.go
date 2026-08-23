package main

// seedOpsRecords provides the initial operations records for the service.
func seedOpsRecords() []OpsRecord {
	return []OpsRecord{
		{ID: "ops-1001", Subject: "青松岭南坡巡护", Owner: "crew-a", Status: OpsStatusActive, Priority: OpsPriorityHigh, Revision: 1, Labels: map[string]string{"site": "qsl-01", "operator": "crew-a"}},
		{ID: "ops-1002", Subject: "溪谷保护区瞭望", Owner: "crew-b", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Revision: 1, Labels: map[string]string{"site": "xg-02", "operator": "crew-b"}},
	}
}
