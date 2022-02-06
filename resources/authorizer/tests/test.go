package main

import (
	"fmt"
	_ "github.com/jptosso/coraza-libinjection"
	_ "github.com/jptosso/coraza-pcre"
	"github.com/jptosso/coraza-waf/v2"
	"github.com/jptosso/coraza-waf/v2/seclang"
)

// our router
var waf *coraza.Waf

func init() {
	waf = coraza.NewWaf()
	parser, _ := seclang.NewParser(waf)

	/*if err := parser.FromString(`SecAction "id:1,phase:1,deny:403,log"`); err != nil {
		panic(err)
	}*/

	files := []string{
		"wafconfig/coraza.conf",
		"wafconfig/coreruleset/crs-setup.conf",
		"wafconfig/coreruleset/rules/*.conf",
		"wafconfig/coraza-additional-rules.conf",
	}
	for _, f := range files {
		if err := parser.FromFile(f); err != nil {
			panic(err)
		}
	}

}

func main() {

	tx := waf.NewTransaction()

	defer func() {
		tx.ProcessLogging()
		tx.Clean()
	}()

	tx.ProcessConnection("192.168.1.100", 55555, "1.1.1.1", 443)

	tx.ProcessURI(fmt.Sprintf("%s?%s", "/hello", "q=SELECT * from password where user = 'root'"), "GET", "HTTP/1.1")

	tx.AddRequestHeader("X-File-Name", "shell.php")

	it1 := tx.ProcessRequestHeaders()

	fmt.Printf("phase1: %#v\n", it1)

	it2, _ := tx.ProcessRequestBody()

	fmt.Printf("phase2: %#v\n", it2)

	al := tx.AuditLog()
	if len(al.Messages) > 0 {
		for _, log := range al.Messages {
			fmt.Printf("messages: %#v\n", log)
		}

	}

}
