package main

func opsRules12() []OpsRule {
	return []OpsRule{
		opsRule1201(),
		opsRule1202(),
		opsRule1203(),
		opsRule1204(),
		opsRule1205(),
		opsRule1206(),
		opsRule1207(),
		opsRule1208(),
	}
}

func opsRule1201() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1201",
		Name:           "geothermal-well-monitor-service control 1201",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1202() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1202",
		Name:           "geothermal-well-monitor-service control 1202",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1203() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1203",
		Name:           "geothermal-well-monitor-service control 1203",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1204() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1204",
		Name:           "geothermal-well-monitor-service control 1204",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule1205() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1205",
		Name:           "geothermal-well-monitor-service control 1205",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1206() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1206",
		Name:           "geothermal-well-monitor-service control 1206",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1207() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1207",
		Name:           "geothermal-well-monitor-service control 1207",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1208() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1208",
		Name:           "geothermal-well-monitor-service control 1208",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
