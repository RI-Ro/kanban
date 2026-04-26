package main

import (
        "io/ioutil"
	"os/exec"
	"fmt"
	"flag"
	"strings"
)

func main(){
	URL := flag.String("url", "http://localhost:8000/domain_list.txt", "domains list in txt file.")
	PATH_TICKET := flag.String("path_ticket","/tmp/admin_ticket","Path to admin ticket for this script")
	ADMIN_PRINCIPAL := flag.String("admin_principal","admin/admin", "Admin principal")
        flag.Parse()
	PATH_DOMAIN := "/tmp/domain_list.txt"
	get_domain_file(*URL, PATH_DOMAIN)
        domain_list := get_domain_list(PATH_DOMAIN)
        fmt.Printf("%v", domain_list)
        get_admin_ticket(*PATH_TICKET, *ADMIN_PRINCIPAL)
	get_kinit(*PATH_TICKET, *ADMIN_PRINCIPAL)
}


func get_domain_file(URL string, PATH_DOMAIN string){
    cmd := exec.Command("wget", "-O", PATH_DOMAIN, URL)
    output, err := cmd.Output()
    if err != nil {
	fmt.Printf("%v \n", err)
	}
    if err == nil {
       	fmt.Printf("%v \n", output)
	}
}

func get_admin_ticket(PATH_TICKET string, ADMIN_PRINCIPAL string) {
    cmd := exec.Command("kadmin.local", "-q", fmt.Sprintf("ktadd -norandkey -k %s %s", PATH_TICKET, ADMIN_PRINCIPAL))
    output, err := cmd.Output()
    if err != nil {
	fmt.Printf("err: %v \n", err)
	}
    if err == nil {
       	fmt.Printf("%v \n", output)
	}
}

func get_kinit(PATH_TICKET string, ADMIN_PRINCIPAL string) {
    cmd := exec.Command("kinit", "-kt", PATH_TICKET, ADMIN_PRINCIPAL)
    output, err := cmd.Output()
    if err != nil {
	fmt.Printf("err: %v \n", err)
	}
    if err == nil {
       	fmt.Printf("%v \n", output)
	}
}

func get_domain_list(PATH_DOMAIN string) (domain_list []string) {
	// 1. Читаем весь файл
	content, err := ioutil.ReadFile(PATH_DOMAIN)
	if err != nil {
		panic(err)
	}
	// 2. Преобразуем байты в строку, удаляем лишние пробелы и разбиваем по новой строке
	domain_list = strings.Split(strings.TrimSpace(string(content)), "\n")
	// Результат: []string
	return domain_list
}

func get_my_domain_name() {

}
