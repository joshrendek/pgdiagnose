package main

import (
	"strings"
)

type Plan struct {
	ConnectionLimit int
}

func GetPlan(name string) Plan {
	switch trimName(name) {
	case "dev", "basic":
		return Plan{20}
	case "crane", "yanari":
		return Plan{60}
	case "kappa":
		return Plan{120}
	case "ronin", "tengu", "fugu":
		return Plan{200}
	case "ika":
		return Plan{400}
	case "baku", "mecha", "ryu":
		return Plan{500}
	}
	return Plan{}
}

func trimName(name string) string {
	name = strings.TrimPrefix(name, "enterprise-")
	name = strings.TrimPrefix(name, "premium-")
	name = strings.TrimPrefix(name, "standard-")
	name = strings.TrimPrefix(name, "hobby-")
	return name
}
