package module

import "fmt"

import "github.com/fatih/color"
import "github.com/go-ldap/ldap/v3"

func decodeSID(sid []byte) string {
	
	if len(sid) < 8 { return "" }

	revision := sid[0]
	subAuthorityCount := sid[1]

	var identifierAuthority int64
	for i := 2; i < 8; i++ { identifierAuthority = (identifierAuthority << 8) + int64(sid[i]) }

	result := fmt.Sprintf("S-%d-%d", revision, identifierAuthority)

	for i := 0; i < int(subAuthorityCount); i++ {
		start := 8 + i*4
		if start+4 <= len(sid) {
			subAuthority := uint32(sid[start]) | uint32(sid[start+1])<<8 | uint32(sid[start+2])<<16 | uint32(sid[start+3])<<24
			result = fmt.Sprintf("%s-%d", result, subAuthority)
		}
	}
	return result
}

func SID(domaindn string, conn *ldap.Conn){

	filter := "(userAccountControl:1.2.840.113556.1.4.803:=8192)" // valeu crackmapexec

	statusMsg.Printf("(QUERY) %s\n\n", filter)

	searchRequest := ldap.NewSearchRequest(
		domaindn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"objectSid"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	error(err)

	green := color.New(color.FgGreen).SprintFunc()

	for _, queryresult := range sr.Entries{
		fmt.Printf("%s %s\n", green("(SID)"), decodeSID(queryresult.GetRawAttributeValue("objectSid")))
	}

}