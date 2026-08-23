package main

func opsRules11() []OpsRule {
	return []OpsRule{
		opsRule1101(),
		opsRule1102(),
		opsRule1103(),
		opsRule1104(),
		opsRule1105(),
		opsRule1106(),
		opsRule1107(),
		opsRule1108(),
	}
}

func opsRule1101() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1101",
		Name:           "geothermal-well-monitor-service control 1101",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1102() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1102",
		Name:           "geothermal-well-monitor-service control 1102",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1103() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1103",
		Name:           "geothermal-well-monitor-service control 1103",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1104() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1104",
		Name:           "geothermal-well-monitor-service control 1104",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule1105() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1105",
		Name:           "geothermal-well-monitor-service control 1105",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1106() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1106",
		Name:           "geothermal-well-monitor-service control 1106",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1107() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1107",
		Name:           "geothermal-well-monitor-service control 1107",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule1108() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-1108",
		Name:           "geothermal-well-monitor-service control 1108",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
