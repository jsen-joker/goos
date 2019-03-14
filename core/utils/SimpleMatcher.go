package utils

import (
	"strings"
)

const SINGLE = "*"
const DOUBLE = "**"

func Matcher(path string, patterns []string) bool {
	for _, v := range patterns {
		if match(path, v) {
			return true
		}
	}
	return false
}

func match(path string, pattern string) bool {

	if path == pattern {
		return true
	}

	ps := strings.Split(path, "/")
	pns := strings.Split(pattern, "/")
	start := 0
	end := len(ps)
	if ps[0] == "" {
		start ++
	}
	if ps[end - 1] == "" {
		end --
	}
	ps = ps[start : end]

	start = 0
	end = len(pns)
	if pns[0] == "" {
		start ++
	}
	if pns[end - 1] == "" {
		end --
	}
	pns = pns[start : end]

	hasDouble := false
	workPs := 0
	workPns := 0
	lenPs := len(ps)
	lenPns := len(pns)

	for true {
		if workPs >= lenPs {
			if workPns >= lenPns {
				return true
			} else {
				if pns[workPns] == DOUBLE {
					workPns++
					continue
				} else {
					return false
				}
			}
		} else {
			if workPns >= lenPns {
				return false
			}
		}


		a := ps[workPs]
		b := pns[workPns]
		if !checkPart(a, b) {
			if hasDouble {
				workPs++
			} else {
				return false
			}
		}
		if b == DOUBLE {
			hasDouble = true
		}

		workPs ++
		workPns ++

		if hasDouble {
			if workPns >= lenPns {
				return true
			}
		}


	}
	return false
}

func checkPart(ps string, pns string) bool {
	if ps == pns {
		return true
	}
	if pns == SINGLE {
		return true
	}
	if pns == DOUBLE {
		return true
	}
	if strings.Contains(pns, "*") && !strings.Contains(pns, "**") {
		return checkArea(ps, pns)
	}
	return false
}

func checkArea(ps string, pns string) bool {
	pnss := strings.Split(pns, "*")

	lenPnss := len(pnss)
	workPnss := 0

	for true {
		if (len(ps) == 0) {
			if workPnss >= lenPnss {
				return true
			}
			return false
		} else if workPnss >= lenPnss {
			if strings.LastIndex(pns, "*") == len(pns) - 1 {
				return true
			}
			return false
		}
		if pnss[workPnss] == "" {
			workPnss ++
			continue
		}

		if strings.Index(ps, pnss[workPnss]) == 0 {
			ps = ps[len(pnss[workPnss]) : ]
			workPnss ++
		} else if workPnss == 0 && strings.Index(pns, "*") == 0 && strings.Index(ps, pnss[workPnss]) > 0 {
			ps = ps[len(pnss[workPnss]) + strings.Index(ps, pnss[workPnss]) : ]
			workPnss ++
		} else if workPnss > 0 && strings.Index(ps, pnss[workPnss]) > 0 {
			ps = ps[len(pnss[workPnss]) + strings.Index(ps, pnss[workPnss]) : ]
			workPnss ++
		} else {
			return false
		}
	}
	return false
}
