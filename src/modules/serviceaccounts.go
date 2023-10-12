package module

import "fmt"

import "github.com/fatih/color"
import "github.com/go-ldap/ldap/v3"

func ServiceAccounts(domaindn string, conn *ldap.Conn){

	filter := "(servicePrincipalName=*)"

	statusMsg.Printf("(QUERY) %s\n\n", filter)

	searchRequest := ldap.NewSearchRequest(
		domaindn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"sAMAccountName"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	error(err)

	green := color.New(color.FgGreen).SprintFunc()

	for _, queryresult := range sr.Entries{
		fmt.Printf("%s %s\n", green("(SAMACCOUNTNAME)"), queryresult.GetAttributeValue("sAMAccountName"))
	}

}