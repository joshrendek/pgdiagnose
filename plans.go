package main

import (
	"strings"
)

type Plan struct {
	ConnectionLimit int
}

func GetPlan(name string) Plan {
	switch trimNme(name) {
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

func trimNme(name string) string {
	name = strings.TrimLeft(name, "enterprise-")
	name = strings.TrimLeft(name, "premium-")
	name = strings.TrimLeft(name, "standard-")
	name = strings.TrimLeft(name, "hobby-")
	return name
}
