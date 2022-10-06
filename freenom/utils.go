package freenom

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/tzwsoho/go-freenom/freenom"
)

func parseID(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid id: %s", id)
	}
	return parts[1], parts[0], nil // return name and domain
}

func computeID(domain, name string) string {
	return fmt.Sprintf("%s/%s", strings.ToLower(name), domain)
}

func computeFQDN(domain, name string) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(name), domain)
}

func getRecordByName(domain, name string, diagnostics *diag.Diagnostics) (record *freenom.DomainRecord, err error) {

	domainInfo, err := freenom.GetDomainInfo(domain)

	if err != nil {
		diagnostics.AddError(
			"Error reading domain info: "+domain,
			err.Error(),
		)
		return
	}

	foundRecord := false

	for _, r := range domainInfo.Records {
		log.Print("[DEBUG] Record: ", r.Name, r.Type, r.Value, r.Priority, r.TTL)

		if strings.EqualFold(r.Name, name) {
			foundRecord = true
			record = r
			break
		}
	}

	if !foundRecord {
		diagnostics.AddError(
			"Record not found",
			"Record not found "+computeID(domain, name),
		)
		err = fmt.Errorf("record not found")
		return
	}
	return
}

func getAllRecordsByDomainName(domain string, diagnostics *diag.Diagnostics) (records []*freenom.DomainRecord, err error) {

	domainInfo, err := freenom.GetDomainInfo(domain)

	if err != nil {
		diagnostics.AddError(
			"Error reading domain info: "+domain,
			err.Error(),
		)
		return
	}

	for _, r := range domainInfo.Records {
		log.Print("[DEBUG] Record: ", r.Name, r.Type, r.Value, r.Priority, r.TTL)
		records = append(records, r)
	}
	return
}

func getAllRecordsByDomainNameAndValue(domain string, value string, diagnostics *diag.Diagnostics) (records []*freenom.DomainRecord, err error) {

	domainInfo, err := freenom.GetDomainInfo(domain)

	if err != nil {
		diagnostics.AddError(
			"Error reading domain info: "+domain,
			err.Error(),
		)
		return
	}

	for _, r := range domainInfo.Records {
		log.Print("[DEBUG] Record: ", r.Name, r.Type, r.Value, r.Priority, r.TTL)
		if r.Value == value {
			records = append(records, r)
		}
	}
	return
}
