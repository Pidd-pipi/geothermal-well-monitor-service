package main

func opsRules14() []OpsRule {
	return []OpsRule{
		opsRule1401(),
		opsRule1402(),
		opsRule1403(),
		opsRule1404(),
		opsRule1405(),
		opsRule1406(),
		opsRule1407(),
		opsRule1408(),
	}
}

func opsRule1401() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1401",
		Name:           "geothermal-well-monitor-service control 1401",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1402() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1402",
		Name:           "geothermal-well-monitor-service control 1402",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1403() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1403",
		Name:           "geothermal-well-monitor-service control 1403",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1404() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1404",
		Name:           "geothermal-well-monitor-service control 1404",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule1405() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1405",
		Name:           "geothermal-well-monitor-service control 1405",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1406() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1406",
		Name:           "geothermal-well-monitor-service control 1406",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1407() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1407",
		Name:           "geothermal-well-monitor-service control 1407",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1408() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1408",
		Name:           "geothermal-well-monitor-service control 1408",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
