package main

func opsRules09() []OpsRule {
	return []OpsRule{
		opsRule0901(),
		opsRule0902(),
		opsRule0903(),
		opsRule0904(),
		opsRule0905(),
		opsRule0906(),
		opsRule0907(),
		opsRule0908(),
	}
}

func opsRule0901() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0901",
		Name:           "geothermal-well-monitor-service control 0901",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0902() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0902",
		Name:           "geothermal-well-monitor-service control 0902",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0903() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0903",
		Name:           "geothermal-well-monitor-service control 0903",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0904() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0904",
		Name:           "geothermal-well-monitor-service control 0904",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule0905() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0905",
		Name:           "geothermal-well-monitor-service control 0905",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0906() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0906",
		Name:           "geothermal-well-monitor-service control 0906",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0907() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0907",
		Name:           "geothermal-well-monitor-service control 0907",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0908() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0908",
		Name:           "geothermal-well-monitor-service control 0908",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
