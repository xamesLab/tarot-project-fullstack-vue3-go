package models

import (
	"regexp"
	"strings"
)

var userRolePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{1,31}$`)

func NormalizeRole(raw string) string {
	return strings.TrimSpace(strings.ToLower(raw))
}

func IsValidRole(raw string) bool {
	return userRolePattern.MatchString(NormalizeRole(raw))
}

func DisplayNameOrLogin(displayName, login string) string {
	name := strings.TrimSpace(displayName)
	if name == "" {
		return login
	}
	return name
}

func NormalizeUserStatus(status UserStatus) UserStatus {
	if strings.TrimSpace(string(status)) == "" {
		return UserStatusActive
	}
	return status
}
