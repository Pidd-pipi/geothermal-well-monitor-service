package main

func opsRules07() []OpsRule {
	return []OpsRule{
		opsRule0701(),
		opsRule0702(),
		opsRule0703(),
		opsRule0704(),
		opsRule0705(),
		opsRule0706(),
		opsRule0707(),
		opsRule0708(),
	}
}

func opsRule0701() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 1%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0701",
		Name:           "geothermal-well-monitor-service control 0701",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0702() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 2%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0702",
		Name:           "geothermal-well-monitor-service control 0702",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0703() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 3%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0703",
		Name:           "geothermal-well-monitor-service control 0703",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0704() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 4%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0704",
		Name:           "geothermal-well-monitor-service control 0704",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       true,
	}
}

func opsRule0705() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 5%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0705",
		Name:           "geothermal-well-monitor-service control 0705",
		Severity:       OpsPriorityLow,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0706() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 6%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0706",
		Name:           "geothermal-well-monitor-service control 0706",
		Severity:       OpsPriorityNormal,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0707() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 7%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0707",
		Name:           "geothermal-well-monitor-service control 0707",
		Severity:       OpsPriorityHigh,
		RequiredLabels: labels,
		Terminal:       false,
	}
}

func opsRule0708() OpsRule {
	labels := []string{"site", "operator", "evidence"}
	if 8%2 == 0 {
		labels = append(labels, "reviewed")
	}
	return OpsRule{
		Code:           "OPS-0708",
		Name:           "geothermal-well-monitor-service control 0708",
		Severity:       OpsPriorityCritical,
		RequiredLabels: labels,
		Terminal:       true,
	}
}
