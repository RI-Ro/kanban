package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

func main() {
	PATH_TICKET := "/tmp/admin_ticket"
	ADMIN_PRINCIPAL := "admin/admin"
	PASSWORD_FOR_TRUST := "very$ecretPASSW0RD21212121"

	DOMAIN_FOR_ADD := flag.String("domain", "", "domain example: .domain.com")
	DC_FOR_ADD := flag.String("dc", "", "domain controller example: dc.domain.com")
	flag.Parse()
	if *DOMAIN_FOR_ADD != "" && *DC_FOR_ADD != "" {
		// Получаем админский keytab и сохраняем его в файл для дальнейшего использования
		fmt.Println("1")
		get_admin_ticket(PATH_TICKET, ADMIN_PRINCIPAL)
		// Получаем админский тикет
		fmt.Println("2")
		get_kinit(PATH_TICKET, ADMIN_PRINCIPAL)
		// Производим добавление в список доверенных доменов
		fmt.Println("3")
		create_trust_domain(*DOMAIN_FOR_ADD, *DC_FOR_ADD, PASSWORD_FOR_TRUST)
		// Перезапускаем ald-client для обновления конфига
		fmt.Println("4")
		restart_ald_client()
		// Удаляем админский keytab
		fmt.Println("5")
		delete_admin_ticket(PATH_TICKET)
	} else {
		fmt.Println("Не заданы обязательные параметры!\nОбязательно должны быть указаны наименование домена и контроллер домена!\nПример:\ntrust-add -domain=.domain.com -dc=dc.domain.com\n")
	}
}

func get_admin_ticket(PATH_TICKET string, ADMIN_PRINCIPAL string) {
	fmt.Println("get_admin_ticket start!")
	cmd := exec.Command("kadmin.local", "-q", fmt.Sprintf("ktadd -norandkey -k %s %s", PATH_TICKET, ADMIN_PRINCIPAL))
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("get_admin_ticket err: %s \n", err.Error())
	}
	if err == nil {
		fmt.Printf("get_admin_ticket output: %s \n", output)
	}
	fmt.Println("get_admin_ticket end!")
}

func get_kinit(PATH_TICKET string, ADMIN_PRINCIPAL string) {
	fmt.Println("get_kinit start!")
	cmd := exec.Command("kinit", "-kt", PATH_TICKET, ADMIN_PRINCIPAL)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf(" get_kinit err: %s \n", err.Error())
	}
	if err == nil {
		fmt.Printf("get_kinit output: %s \n", output)
	}
	fmt.Println("get_kinit end!")
}

func create_trust_domain(DOMAIN_FOR_ADD string, DC_FOR_ADD string, password string) {
	fmt.Println("create_trust_domain start!")
	// 1. Подготавливаем команду
	cmd := exec.Command("ald-admin", "trusted-add", DOMAIN_FOR_ADD, "-c", fmt.Sprintf("--kdc=%s", DC_FOR_ADD),
		"--direction=two-way", fmt.Sprintf("--desc=%s", DOMAIN_FOR_ADD))

	// 2. Запускаем команду в PTY
	ptyFile, err := pty.Start(cmd)
	if err != nil {
		fmt.Printf("ошибка запуска PTY: %w", err)
	}
	defer ptyFile.Close()

	// 3. Настраиваем чтение вывода PTY
	reader := bufio.NewReader(ptyFile)

	// 4. Счётчик отправленных ответов (для выхода из цикла)
	answersSent := 0
	expectedPrompts := []string{"':", "пароль:", "(yes/no)", "[no]"}

	// 5. Буфер для накопления вывода (на случай, если приглашение не заканчивается переводом строки)
	var outputBuffer strings.Builder

	// 6. Канал для сигнала об окончании работы
	done := make(chan error)

	// Горyтина: читаем вывод PTY и реагируем на приглашения
	go func() {
		for {
			// Читаем по одному байту или строке – удобнее посимвольно, чтобы поймать неполные строки
			// Используем ReadString('\n') для простоты, но некоторые приглашения могут не иметь '\n'
			// Более надёжно – читать по байтам и накапливать.
			chunk, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// Команда завершилась
					done <- nil
				} else {
					done <- fmt.Errorf("ошибка чтения: %w", err)
				}
				return
			}

			// Печатаем вывод (опционально, для отладки/логирования)
			fmt.Print(chunk)

			// Добавляем прочитанное в буфер для анализа
			outputBuffer.WriteString(chunk)

			// Проверяем, не появилось ли одно из приглашений
			lowerOutput := strings.ToLower(outputBuffer.String())
			triggered := false

			for _, prompt := range expectedPrompts {
				if strings.Contains(lowerOutput, strings.ToLower(prompt)) {
					// Нашли ключевое слово – отправляем ответ
					if answersSent == 0 {
						// Первый запрос – ничего
						fmt.Fprintln(ptyFile, "\n")
						answersSent++
						triggered = true
						break
					}
					if answersSent == 1 {
						// Первый запрос – пароль
						fmt.Fprintln(ptyFile, password)
						answersSent++
						triggered = true
						break
					} else if answersSent == 2 {
						// Второй запрос – подтверждение пароля
						fmt.Fprintln(ptyFile, password)
						answersSent++
						triggered = true
						break
					} else if answersSent == 3 {
						// Третий запрос – подтверждение yes
						fmt.Fprintln(ptyFile, "yes")
						answersSent++
						triggered = true
						break
					}
				}
			}

			if triggered {
				// После отправки сбрасываем буфер, чтобы не сработало повторно на том же тексте
				outputBuffer.Reset()
			}

			// Если отправили все ответы, можно перестать искать приглашения,
			// но продолжать читать вывод до завершения команды.
			if answersSent >= 4 {
				// Просто сливаем оставшийся вывод
				// Для простоты – продолжаем читать до EOF, ничего не отправляя
			}
		}
	}()

	// 7. Ждём завершения команды и обработки вывода
	err = <-done
	if err != nil {
		fmt.Printf("команда завершилась с ошибкой: %s", err.Error())
	}

	// 8. Дожидаемся, пока процесс действительно завершится (освободим ресурсы)
	if waitErr := cmd.Wait(); waitErr != nil {
		fmt.Printf("команда завершилась с ошибкой: %s", waitErr.Error())
	}
}

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

func delete_admin_ticket(PATH_TICKET string) {
	fmt.Println("delete_admin_ticket start!")
	cmd := exec.Command("kdestroy")
	_, err := cmd.Output()
	if err != nil {
		fmt.Printf("delete_admin_ticket err: %s \n", err.Error())
	}
	err = os.Remove(PATH_TICKET)
	if err != nil {
		fmt.Printf("delete_admin_ticket err: %s \n", err.Error())
	}
	fmt.Println("delete_admin_ticket end!")
}
