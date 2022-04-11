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
	return fmt.Sprintf("%s/%s", name, domain)
}

func getRecordByName(domain, name string, diagnostics *diag.Diagnostics) (record *freenom.DomainRecord, err error) {

	domainInfo, err := freenom.GetDomainInfo(domain)

	if err != nil {
		diagnostics.AddError(
			"Error reading domain info: "+domain,
			err.Error(),
		)
		return
		// return freenom.DomainRecord{}, err
	}

	foundRecord := false
	// var record *freenom.DomainRecord

	for _, r := range domainInfo.Records {
		log.Print("[DEBUG] Record: ", r.Name, r.Type, r.Value, r.Priority, r.TTL)

		if strings.ToLower(r.Name) == strings.ToLower(name) {
			foundRecord = true
			record = r
			break
		}
	}

	if !foundRecord {
		diagnostics.AddError(
			"Record not found",
			"Record not found"+computeID(domain, name))
		err = fmt.Errorf("Record not found")
		return
	}
	return
}
