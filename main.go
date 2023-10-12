package main

import "os"
import "fmt"
import "bufio"
import "regexp"
import "strings"
import "io/ioutil"

import "belchior/src/conx"
import "github.com/fatih/color"
import "github.com/jessevdk/go-flags"

func error(err interface{}){ if err != nil { fmt.Println(err) ; os.Exit(0) } }
func banner() {  arq, _ := ioutil.ReadFile("ui/ascii.txt")  ; fmt.Println(string(arq)) }

func listmodules() {
	arq, err := os.Open("ui/modules.txt")
	error(err)

	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	for scanner.Scan(){
		line := scanner.Text()
		line = strings.Replace(line, "[MODULE]", green("[MODULE]"), -1)
		//fmt.Println(line)

		re := regexp.MustCompile(`\]\s*([^()]*)\s*\(`)
		matches := re.FindStringSubmatch(line)

		re_espaco := regexp.MustCompile(`\s`)
		noSpaces := re_espaco.ReplaceAllString(matches[1], "")

		line2 := strings.Replace(line, noSpaces, yellow(noSpaces), -1)
		fmt.Println(line2)
	}
}

func main(){
	banner()

	var opts struct {
		Target string `short:"t" long:"target" description:"Definition: TARGET to authenticate in LDAP" required:"true"`
		User string `short:"u" long:"username" description:"Definition: USERNAME to authenticate in LDAP" required:"false"`
		Password string `short:"p" long:"password" description:"Definition: PASSWORD to authenticate in LDAP" required:"false"`
		Module string `short:"M" long:"module" description:"Definition: MODULE to use in LDAP" required:"false"`
		TargetUser string `short:"T" long:"target-user" description:"Definition: TARGET USER to use in LDAP" required:"false"`
		ListModules bool `short:"L" long:"list-modules" description:"Definition: LIST modules to use in LDAP" required:"false"`
	}

	_, err := flags.Parse(&opts) ; if err != nil { os.Exit(0) }

	if opts.ListModules == true {
		listmodules() ; os.Exit(0)
	}

	if (strings.Contains(opts.Target, "ldap://") == false){
		fmt.Println("[!] Use este formato: 'ldap://target.com'")
		os.Exit(0)
	}

	slashIndex := strings.Index(opts.User, "/") ; atIndex := strings.LastIndex(opts.User, "/")

	domain := opts.User[:slashIndex]
	username := opts.User[atIndex+1:]
	usuario_final := fmt.Sprintf("%s@%s", username, domain)

	con.Mainconnection(opts.Target, usuario_final, username, opts.Password, opts.Module, opts.TargetUser)
}

// (objectCategory=person)(objectClass=user)(pwdLastSet=0)(!useraccountcontrol:1.2.840.113556.1.4.803:=2)
// user must change password at next logon
//  Set-ADAccountPassword -Identity 'CN=Administrator,CN=Users,DC=htb,DC=local' -Reset -NewPassword (ConvertTo-SecureString -AsPlainText "SenhaFoda123!@#" -Force)

// "(&(objectCategory=computer)(userAccountControl:1.2.840.113556.1.4.803:=8192))"
// o de cima lista todos os dc

// "(&(objectCategory=person)(objectClass=user)(memberOf=CN=YourGroup,OU=Users,DC=domain,DC=com))"
// o de cima lista todos os usuarios de um grupo

// (&(objectCategory=person)(objectClass=user)(userAccountControl:1.2.840.113556.1.4.803:=16777216))
// o de cima lista todos os usuarios que podem delegar autenticação
