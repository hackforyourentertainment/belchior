package module

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/go-ldap/ldap/v3"
)

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.New(color.FgYellow).SprintFunc()(prompt))
	userInput, _ := reader.ReadString('\n')
	return strings.TrimSpace(userInput)
}

func getSearchResults(conn *ldap.Conn, domaindn string, searchFilter string) *ldap.SearchResult {
	searchRequest := ldap.NewSearchRequest(
		domaindn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{"dn"},
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	error(err)

	return searchResult
}

func ManageGroups(domaindn string, conn *ldap.Conn) {
	targetuser := getUserInput("(LDAP) Insira o usuário: ")
	userSearchFilter := fmt.Sprintf("(&(objectClass=user)(objectCategory=person)(sAMAccountName=%s))", targetuser)
	userSearchResult := getSearchResults(conn, domaindn, userSearchFilter)

	if len(userSearchResult.Entries) != 1 {
		errorMsg.Printf("(ERR) Não foi possível encontrar o usuário'%s'\n(ERR) %s\n", targetuser, userSearchFilter)
		os.Exit(0)
	}

	successMsg.Printf("(LDAP) %s\n", userSearchResult.Entries[0].DN)
	userDN := userSearchResult.Entries[0].DN

	group := getUserInput("(LDAP) Insira o grupo: ")
	groupSearchFilter := fmt.Sprintf("(&(objectClass=group)(cn=%s))", group)
	groupSearchResult := getSearchResults(conn, domaindn, groupSearchFilter)

	if len(groupSearchResult.Entries) != 1 {
		errorMsg.Printf("(ERR) Não foi possível encontrar o grupo '%s'\n", group)
		os.Exit(0)
	}

	successMsg.Printf("(LDAP) %s\n", groupSearchResult.Entries[0].DN)
	groupDN := groupSearchResult.Entries[0].DN

	option := getUserInput("(LDAP) Insira a opção [add/rm]: ")

	if option == "add" {
		modify := ldap.NewModifyRequest(groupDN, nil)
		modify.Add("member", []string{userDN})
		err := conn.Modify(modify)

		error(err)

		successMsg.Printf("(LDAP) Usuário '%s' adicionado ao grupo '%s'!\n", targetuser, group)
	} else if option == "rm" {
		modify := ldap.NewModifyRequest(groupDN, nil)
		modify.Delete("member", []string{userDN})
		err := conn.Modify(modify)

		error(err)
		successMsg.Printf("(LDAP) User '%s' removido do grupo '%s'!\n", targetuser, group)
	} else {
		errorMsg.Printf("(ERR) Opção '%s' inválida.\n", option)
	}
}
