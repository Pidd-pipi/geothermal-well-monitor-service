package main

func opsRules10() []OpsRule {
	return []OpsRule{
		opsRule1001(),
		opsRule1002(),
		opsRule1003(),
		opsRule1004(),
		opsRule1005(),
		opsRule1006(),
		opsRule1007(),
		opsRule1008(),
	}
}

func opsRule1001() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1001",
		Name:           "geothermal-well-monitor-service control 1001",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1002() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1002",
		Name:           "geothermal-well-monitor-service control 1002",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1003() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1003",
		Name:           "geothermal-well-monitor-service control 1003",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1004() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1004",
		Name:           "geothermal-well-monitor-service control 1004",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule1005() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1005",
		Name:           "geothermal-well-monitor-service control 1005",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1006() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1006",
		Name:           "geothermal-well-monitor-service control 1006",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1007() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1007",
		Name:           "geothermal-well-monitor-service control 1007",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1008() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1008",
		Name:           "geothermal-well-monitor-service control 1008",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
