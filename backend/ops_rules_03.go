package main

func opsRules03() []OpsRule {
	return []OpsRule{
		opsRule0301(),
		opsRule0302(),
		opsRule0303(),
		opsRule0304(),
		opsRule0305(),
		opsRule0306(),
		opsRule0307(),
		opsRule0308(),
	}
}

func opsRule0301() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0301",
		Name:           "geothermal-well-monitor-service control 0301",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0302() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0302",
		Name:           "geothermal-well-monitor-service control 0302",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0303() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0303",
		Name:           "geothermal-well-monitor-service control 0303",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0304() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0304",
		Name:           "geothermal-well-monitor-service control 0304",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule0305() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0305",
		Name:           "geothermal-well-monitor-service control 0305",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0306() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0306",
		Name:           "geothermal-well-monitor-service control 0306",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0307() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0307",
		Name:           "geothermal-well-monitor-service control 0307",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0308() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0308",
		Name:           "geothermal-well-monitor-service control 0308",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
