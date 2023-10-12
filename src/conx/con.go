package con

import "os"

import "belchior/src/modules"
import "github.com/fatih/color"
import "github.com/go-ldap/ldap/v3"

var statusMsg *color.Color = color.New(color.FgYellow)
var successMsg *color.Color = color.New(color.FgGreen)
var errorMsg *color.Color = color.New(color.FgRed)

func error(err interface{}){ if err != nil { errorMsg.Println("[-]", err) ; os.Exit(0) } }

func GetDomainDN(conn *ldap.Conn) string {
	
	searchRequest := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(&(objectClass=domain)(objectClass=domainDNS))", // query
		[]string{"defaultNamingContext"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	error(err)

	return sr.Entries[0].GetAttributeValue("defaultNamingContext")
}

func Mainconnection(host, usuariofinal, usuario, senha, modulo, targetuser string){
	
	statusMsg.Printf("(LDAP) Tentativa de conexão TCP em %s\n", host)

	conn, err := ldap.DialURL(host)
	error(err)

	defer conn.Close()

	err = conn.Bind(usuariofinal, senha)
	
	if err != nil {
		errorMsg.Printf("(ERROR) %s:%s\n", usuario, senha)
		os.Exit(0)
	}

	successMsg.Printf("(LDAP) %s:%s\n", usuario, senha)

	if modulo == "whoami" {
		if len(targetuser) == 0{ module.Whoami(conn, GetDomainDN(conn), usuario) } else { module.Whoami(conn, GetDomainDN(conn), targetuser) }
	} else if modulo == "passwdmustchange" {
		module.PasswdMustChange(GetDomainDN(conn), conn)
	} else if modulo == "users" {
		module.ListUsers(GetDomainDN(conn), conn)
	} else if modulo == "passpol"{
		module.PasswordPolicy(GetDomainDN(conn), conn)
	} else if modulo == "serviceaccounts"{
		module.ServiceAccounts(GetDomainDN(conn), conn)
	} else if modulo == "listadmins"{
		module.ListAdmins(GetDomainDN(conn), conn)
	} else if modulo == "passwdnotreq"{
		module.PasswdNotReq(GetDomainDN(conn), conn)
	} else if modulo == "managegroup"{
		module.ManageGroups(GetDomainDN(conn), conn)
	} else if modulo == "listdc"{
		module.ListDC(GetDomainDN(conn), conn)
	} else if modulo == "delegateauth"{
		module.DelegateAuth(GetDomainDN(conn), conn)
	} else if modulo == "sid"{
		if len(targetuser) == 0{ module.ChangePASSWD(conn, usuario, GetDomainDN(conn)) } else { module.ChangePASSWD(conn, targetuser, GetDomainDN(conn)) }
	}
}