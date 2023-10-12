package module

import "fmt"
import "time"
import "strconv"

import "github.com/fatih/color"
import "github.com/go-ldap/ldap/v3"

func convert(tempo string) time.Duration {

	lockout, err := strconv.Atoi(tempo)
	error(err)

	lockoutDuration := int64(lockout)
	lockoutDurationSec := (lockoutDuration / -10000000)

	lockoutDurationHuman := time.Duration(lockoutDurationSec) * time.Second
	return lockoutDurationHuman
}

func PasswordPolicy(domaindn string, conn *ldap.Conn){

	filter := "(objectClass=domainDNS)"

	statusMsg.Printf("(QUERY) %s\n\n", filter)

	searchRequest := ldap.NewSearchRequest(
		domaindn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"minPwdLength", "lockoutThreshold", "lockoutDuration"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	error(err)

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	for _, queryresult := range sr.Entries{
		fmt.Printf("%s %s\n", green("(MINPWDLENGHT)"), queryresult.GetAttributeValue("minPwdLength"))
		fmt.Printf("%s %s => %s\n", green("(LOCKOUTTHRESHOLD)"), queryresult.GetAttributeValue("lockoutThreshold"), red("INFO para password spraying!"))
		fmt.Printf("%s %v\n", green("(LOCKOUTDURATION)"), convert(queryresult.GetAttributeValue("lockoutDuration")))
	}
}