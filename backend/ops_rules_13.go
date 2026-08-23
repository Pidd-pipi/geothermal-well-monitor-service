package main

func opsRules13() []OpsRule {
	return []OpsRule{
		opsRule1301(),
		opsRule1302(),
		opsRule1303(),
		opsRule1304(),
		opsRule1305(),
		opsRule1306(),
		opsRule1307(),
		opsRule1308(),
	}
}

func opsRule1301() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1301",
		Name:           "geothermal-well-monitor-service control 1301",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1302() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1302",
		Name:           "geothermal-well-monitor-service control 1302",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1303() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1303",
		Name:           "geothermal-well-monitor-service control 1303",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1304() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1304",
		Name:           "geothermal-well-monitor-service control 1304",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule1305() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1305",
		Name:           "geothermal-well-monitor-service control 1305",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1306() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1306",
		Name:           "geothermal-well-monitor-service control 1306",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1307() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1307",
		Name:           "geothermal-well-monitor-service control 1307",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1308() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1308",
		Name:           "geothermal-well-monitor-service control 1308",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
