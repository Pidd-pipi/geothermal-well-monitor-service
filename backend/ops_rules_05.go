package main

func opsRules05() []OpsRule {
	return []OpsRule{
		opsRule0501(),
		opsRule0502(),
		opsRule0503(),
		opsRule0504(),
		opsRule0505(),
		opsRule0506(),
		opsRule0507(),
		opsRule0508(),
	}
}

func opsRule0501() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0501",
		Name:           "geothermal-well-monitor-service control 0501",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0502() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0502",
		Name:           "geothermal-well-monitor-service control 0502",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0503() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0503",
		Name:           "geothermal-well-monitor-service control 0503",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0504() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0504",
		Name:           "geothermal-well-monitor-service control 0504",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule0505() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0505",
		Name:           "geothermal-well-monitor-service control 0505",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0506() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0506",
		Name:           "geothermal-well-monitor-service control 0506",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0507() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0507",
		Name:           "geothermal-well-monitor-service control 0507",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0508() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0508",
		Name:           "geothermal-well-monitor-service control 0508",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
