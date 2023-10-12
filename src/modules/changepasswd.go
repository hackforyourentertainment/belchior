package module

import (
	"os"
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func ChangePASSWD(conn *ldap.Conn, domaindn, usuario string){

	filter := fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", ldap.EscapeFilter(usuario))

	searchReq := ldap.NewSearchRequest(
		domaindn,
		ldap.ScopeWholeSubtree, 0, 0, 0, false,
		filter,
		[]string{"sAMAccountName", "dn", "description", "lastLogon", "pwdLastSet", "memberOf"},
		[]ldap.Control{},
	)


	result, err := conn.Search(searchReq)
	if err != nil{
		fmt.Println(err)
		os.Exit(0)
	}

	statusMsg.Printf("(QUERY) %s\n\n", filter)

	if len(result.Entries) != 1{
		errorMsg.Printf("[-] Usuário '%s' não encontrado. Ele realmente existe no domínio?\n", usuario)
		os.Exit(0)
	}

	for _, queryresult := range result.Entries{

		samacc := queryresult.GetAttributeValue("sAMAccountName")

		fmt.Println(samacc)
	}

	/*modify := ldap.NewModifyRequest("CN=username,OU=YourOrganizationalUnit,DC=domain,DC=com")

	newPassword := []byte("\"Nsecure@123\"")
	utf16Password := encodeUTF16LE(newPassword)

	modify.Replace("unicodePwd", []string{string(utf16Password)})

	err = conn.Modify(modify)
	error(err)*/

}