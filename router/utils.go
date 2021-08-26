package router

import (
	"strconv"
	"strings"
)

func getStringParams(param string, def string) string {
	if param == "" {
		return def
	}
	return param
}

func getStringSliceParams(param string, def []string) []string {
	if param == "" {
		return def
	}
	return strings.Split(param, ",")
}

func getIntParams(param string, def int) int {
	ret, err := strconv.ParseInt(param, 10, 32)
	if param == "" || err != nil {
		return def
	}
	return int(ret)
}

func getBoolParams(param string, def bool) bool {
	if param == "" {
		return def
	}
	return strings.EqualFold(param, "true")
}
