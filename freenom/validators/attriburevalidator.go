package validators

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func IsDomain() tfsdk.AttributeValidator {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^((([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]))$`),
		"Invalid domain",
	)
}

func IsIpv4() tfsdk.AttributeValidator {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`),
		"Invalid ipv4",
	)
}

func IsMacAddress() tfsdk.AttributeValidator {
	return stringvalidator.RegexMatches(
		regexp.MustCompile(`^[a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}$`),
		"Invalid mac address",
	)
}
