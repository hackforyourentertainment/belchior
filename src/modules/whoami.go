package module

import "os"
import "fmt"
import "os/exec"
//import "strings"

import "github.com/fatih/color"
import "github.com/go-ldap/ldap/v3"

var statusMsg *color.Color = color.New(color.FgYellow)
var successMsg *color.Color = color.New(color.FgGreen)
var errorMsg *color.Color = color.New(color.FgRed)

func error(err interface{}){ if err != nil { errorMsg.Println("[-]", err) ; os.Exit(0) } }

func execpythonscript(timewindows string) string {

	//comando := fmt.Sprintf("python3 %s", timewindows)

	out, _ := exec.Command("python3" , "src/py/getwindowstime.py", timewindows).Output()
	
	return string(out)

}

func Whoami(conn *ldap.Conn, domaindn, usuario string){

	filter := fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", ldap.EscapeFilter(usuario))
	//searchReq := ldap.NewSearchRequest(domaindn, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{"sAMAccountName"}, []ldap.Control{})

	searchReq := ldap.NewSearchRequest(
		domaindn,
		ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter,
		[]string{"sAMAccountName", "dn", "description", "lastLogon", "pwdLastSet", "memberOf"},
		[]ldap.Control{},
	)

	result, err := conn.Search(searchReq)
	error(err)

	statusMsg.Printf("(QUERY) %s\n\n", filter)

	if len(result.Entries) != 1{
		errorMsg.Printf("[-] Usuário '%s' não encontrado. Ele realmente existe no domínio?\n", usuario)
		os.Exit(0)
	}

	for _, queryresult := range result.Entries{

		samacc := queryresult.GetAttributeValue("sAMAccountName")
		desc := queryresult.GetAttributeValue("description")
		lastlogon := queryresult.GetAttributeValue("lastLogon")
		pwdlastset := queryresult.GetAttributeValue("pwdLastSet")

		green := color.New(color.FgGreen).SprintFunc()

		fmt.Printf("%s %s\n", green("(DESCRIPTION)"), desc)
		fmt.Printf("%s %s\n", green("(DISTINGNAME)"), queryresult.DN)
		fmt.Printf("%s %s\n", green("(SAMACCOUNTNAME)"), samacc)
		fmt.Printf("%s %s", green("(PWDLASTSET)"), execpythonscript(pwdlastset))
		fmt.Printf("%s %s", green("(LASTLOGON)"), execpythonscript(lastlogon))

		for _, attr := range queryresult.Attributes {
			if attr.Name == "memberOf"{
				for _, grupos := range attr.Values{
					groupLine := fmt.Sprintf("%s %s", green("(MEMBEROF)"), grupos)
					fmt.Println(groupLine)
				}
			}
		}
	}
}