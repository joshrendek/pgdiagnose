package main

func CheckLoad(load *float64) []Check {

	var loadCheck Check

	if load == nil {
		reason := make(map[string]string)
		reason["error"] = "Load check not supported on this plan"
		loadCheck = Check{"Load", "skipped", reason}
	} else if *load > 2 {
		loadCheck = Check{"Load", "red", load}
	} else if *load > 1 {
		loadCheck = Check{"Load", "yellow", load}
	} else {
		loadCheck = Check{"Load", "green", nil}
	}

	v := make([]Check, 1)
	v[0] = loadCheck

	return v
}
