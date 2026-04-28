package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	PATH_TICKET := "/tmp/admin_ticket"
	ADMIN_PRINCIPAL := "admin/admin"
	FILE_SCRIPT := "/tmp/ald_add.sh"
	PASSWORD_FOR_TRUST := "LSKfdjpwejvmvdpokwjeoifho!@QW#E$R"

	DOMAIN_FOR_ADD := flag.String("domain", "", "domain example: .domain.com")
	DC_FOR_ADD := flag.String("dc", "", "domain controller example: dc.domain.com")
	flag.Parse()

	// Удаляем скрипт в любом случае
	defer delete_script(FILE_SCRIPT)
	// Удаляем админский keytab в любом случае
	defer delete_admin_ticket(PATH_TICKET)

	if *DOMAIN_FOR_ADD != "" && *DC_FOR_ADD != "" {
		// Получаем админский keytab и сохраняем его в файл для дальнейшего использования
		get_admin_ticket(PATH_TICKET, ADMIN_PRINCIPAL)
		// Получаем админский тикет
		get_kinit(PATH_TICKET, ADMIN_PRINCIPAL)
		// Создаем файл скрипта, чтобы не вводить пароль руками.
		// Пароль будет "забит гвоздями" в скрипте, чтобы обеспечить одинаковость паролей
		create_script_file(FILE_SCRIPT)
		// Производим добавление в список доверенных доменов
		create_trust_domain(*DOMAIN_FOR_ADD, *DC_FOR_ADD, PASSWORD_FOR_TRUST, FILE_SCRIPT)
	} else {
		red := "\033[31m"
		green := "\033[32m"
		reset := "\033[0m"
		fmt.Println(red + "Не заданы обязательные параметры!\nОбязательно должны быть указаны наименование домена и контроллер домена!\nПример:\n" + green + "trust-add -domain=.domain.com -dc=dc.domain.com\n" + reset)
	}
}

func get_admin_ticket(PATH_TICKET string, ADMIN_PRINCIPAL string) {
	cmd := exec.Command("kadmin.local", "-q", fmt.Sprintf("ktadd -norandkey -k %s %s", PATH_TICKET, ADMIN_PRINCIPAL))
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("get_admin_ticket err: %s \n", err.Error())
	}
}

func get_kinit(PATH_TICKET string, ADMIN_PRINCIPAL string) {
	cmd := exec.Command("kinit", "-kt", PATH_TICKET, ADMIN_PRINCIPAL)
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("get_kinit err: %s \n", err.Error())
	}
}

func get_trust() {
	cmd := exec.Command("ald-admin", "trusted-list", "-c")
	result, err := cmd.Output()
	if err != nil {
		fmt.Printf("get_trust err: %s\n", err.Error())
	}
	if err == nil {
		fmt.Printf("%s\n", result)
	}
}

func create_trust_domain(DOMAIN_FOR_ADD string, DC_FOR_ADD string, password string, FILE_SCRIPT string) {
	red := "\033[31m"
	green := "\033[32m"
	reset := "\033[0m"
	// 1. Подготавливаем команду
	cmd := exec.Command(FILE_SCRIPT, DOMAIN_FOR_ADD, DC_FOR_ADD, password)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(red + "При построении доверительных отношений произошла ошибка.\nДля более подробной информации о возникшей ошибке постройте доверительные отношения командой:\nald-admin trusted-add .domain.name" + reset)
		fmt.Printf("create_trust_domain err: %s \n", err.Error())
		return
	}
	if err == nil {
		fmt.Printf("create_trust_domain output: %s \n", output)
		fmt.Println("Список построенных доверительных отношений:\n")
		get_trust()
		fmt.Println(green + "\nОбязательно осуществите перезапуск службы ald-client командой: \nald-client restart\n" + reset)
	}
}

/*
func restart_ald_client() {
	fmt.Println("restart_ald_client start!")
	cmd := exec.Command("/usr/sbin/ald-client", "restart")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("restart_ald_client err: %s \n", err.Error())
	}
	if err == nil {
		fmt.Printf("restart_ald_client output: %s \n", output)
	}
	fmt.Println("restart_ald_client end!")
}
*/

func create_script_file(FILE_SCRIPT string) {
	_ = os.Remove(FILE_SCRIPT)
	text := `#! /usr/bin/expect -f
set DOMAIN [lindex $argv 0]
set DC [lindex $argv 1]
set PASSWORD [lindex $argv 2]

#set timeout 1
spawn ald-admin trusted-add  $DOMAIN -c --kdc=$DC --direction=two-way --desc=$DOMAIN
expect "':"
send "$PASSWORD\r"
expect "пароль:"
send "$PASSWORD\r"
expect "(yes/no)"
send "yes\r"
expect eof
`
	f, err := os.OpenFile(FILE_SCRIPT,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		fmt.Println(err.Error())
	}
}

func delete_script(FILE_SCRIPT string) {
	_ = os.Remove(FILE_SCRIPT)
}

func delete_admin_ticket(PATH_TICKET string) {
	cmd := exec.Command("kdestroy")
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("delete_admin_ticket err: %s \n", err.Error())
	}
	_ = os.Remove(PATH_TICKET)
}
