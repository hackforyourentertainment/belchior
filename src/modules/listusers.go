package module

import "os"
import "fmt"

import "github.com/fatih/color"
import "github.com/go-ldap/ldap/v3"

func outputfile(texto string){
	f, err := os.OpenFile("usuarios.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	error(err)

	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%s", texto)); err != nil {
		fmt.Println(err)
	}

}

func ListUsers(domaindn string, conn *ldap.Conn){

	filter := "(&(objectClass=user)(sAMAccountName=*)(!(userAccountControl:1.2.840.113556.1.4.803:=2)))"

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

	os.Remove("usuarios.txt")

	for _, queryresult := range sr.Entries{
		fmt.Printf("%s %s\n", green("(SAMACCOUNTNAME)"), queryresult.GetAttributeValue("sAMAccountName"))
		outputfile(fmt.Sprintf("%s\n", queryresult.GetAttributeValue("sAMAccountName")))
	}

	statusMsg.Printf("\n(INFO) %s\n", "Lista de usuários salva em 'usuarios.txt'.")
}