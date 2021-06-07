package helper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/laetificat/pricewatcher-api/internal/watcher"
)

// NoSupportedDomainFoundErrorMessage is the standardized error message.
var NoSupportedDomainFoundErrorMessage = "no supported domain is found in the url"

/*
GetSupportedDomains returns the list of supported domains.
*/
func GetSupportedDomains() []string {
	return watcher.SupportedDomains
}

/*
GuessDomain returns a supported domain if it can match one with the given url.
*/
func GuessDomain(url string) (string, error) {
	for _, supportedDomain := range GetSupportedDomains() {
		if matched, _ := regexp.MatchString(".*\\.?"+supportedDomain+".*", url); matched {
			return supportedDomain, nil
		}
	}

	return "", fmt.Errorf(NoSupportedDomainFoundErrorMessage)
}

/*
IsSupported checks if the given domain is present in the list of supported domains.
*/
func IsSupported(domain string) bool {
	supportedTypes := fmt.Sprintf(",%s,", strings.Join(watcher.SupportedDomains, ","))
	return strings.Contains(supportedTypes, fmt.Sprintf(",%s,", domain))
}
