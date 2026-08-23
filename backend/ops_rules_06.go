package main

func opsRules06() []OpsRule {
	return []OpsRule{
		opsRule0601(),
		opsRule0602(),
		opsRule0603(),
		opsRule0604(),
		opsRule0605(),
		opsRule0606(),
		opsRule0607(),
		opsRule0608(),
	}
}

func opsRule0601() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0601",
		Name:           "geothermal-well-monitor-service control 0601",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0602() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0602",
		Name:           "geothermal-well-monitor-service control 0602",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0603() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0603",
		Name:           "geothermal-well-monitor-service control 0603",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0604() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0604",
		Name:           "geothermal-well-monitor-service control 0604",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule0605() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0605",
		Name:           "geothermal-well-monitor-service control 0605",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0606() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0606",
		Name:           "geothermal-well-monitor-service control 0606",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0607() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0607",
		Name:           "geothermal-well-monitor-service control 0607",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0608() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0608",
		Name:           "geothermal-well-monitor-service control 0608",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
