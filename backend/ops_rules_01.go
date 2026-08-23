package main

func opsRules01() []OpsRule {
	return []OpsRule{
		opsRule0101(),
		opsRule0102(),
		opsRule0103(),
		opsRule0104(),
		opsRule0105(),
		opsRule0106(),
		opsRule0107(),
		opsRule0108(),
	}
}

func opsRule0101() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0101",
		Name:           "geothermal-well-monitor-service control 0101",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0102() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0102",
		Name:           "geothermal-well-monitor-service control 0102",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0103() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0103",
		Name:           "geothermal-well-monitor-service control 0103",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0104() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0104",
		Name:           "geothermal-well-monitor-service control 0104",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule0105() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0105",
		Name:           "geothermal-well-monitor-service control 0105",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0106() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0106",
		Name:           "geothermal-well-monitor-service control 0106",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0107() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0107",
		Name:           "geothermal-well-monitor-service control 0107",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0108() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0108",
		Name:           "geothermal-well-monitor-service control 0108",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
