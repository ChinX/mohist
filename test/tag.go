package rflct

import (
	"bytes"
	"strings"
)

var (
	initialisms = []string{
		"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTPS", "HTTP",
		"IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS",
		"TTL","UUID", "UID", "UI", "ID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS",
	}
	initialismReplacer *strings.Replacer
	difference         = 'A' - 'a'
)

func init() {
	var commonInitialismsForReplacer []string
	for _, initialism := range initialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	initialismReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func snakeCasedName(name string) string {
	value := initialismReplacer.Replace(name)
	buf := &bytes.Buffer{}
	var lastCase, currCase, nextCase bool

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] >= 'A' && value[i+1] <= 'Z'
		if i > 0 {
			if currCase == true {
				if lastCase == true && nextCase == true {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
			}
		} else {
			currCase = true
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])
	return strings.ToLower(buf.String())
}
