package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCommandWithSudo(password, name string, args ...string) (string, error) {
	cmd := exec.Command("sudo", append([]string{"-S", name}, args...)...)
	cmd.Stdin = strings.NewReader(password + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	
	wifiInterface := "wlx00367608058b"


	if wifiInterface == "" {
		fmt.Println("Не найден интерфейс беспроводной сети. Проверьте, что ваш беспроводной адаптер правильно установлен.")
		os.Exit(1)
	}

	fmt.Println("Обнаружен интерфейс:", wifiInterface)

	
	fmt.Print("Введите пароль для sudo: ")
	reader := bufio.NewReader(os.Stdin)
	sudoPassword, _ := reader.ReadString('\n')
	sudoPassword = strings.TrimSpace(sudoPassword)

	// Включение интерфейса
	output, err := runCommandWithSudo(sudoPassword, "ip", "link", "set", wifiInterface, "up")
	if err != nil {
		fmt.Printf("Ошибка включения интерфейса: %v\nВывод: %s\n", err, output)
		os.Exit(1)
	}

	// Сканирование доступных сетей
	fmt.Println("Сканирование доступных сетей...")
	output, err = runCommandWithSudo(sudoPassword, "nmcli", "dev", "wifi", "rescan")
	if err != nil {
		fmt.Printf("Ошибка сканирования сетей: %v\nВывод: %s\n", err, output)
		os.Exit(1)
	}

	output, err = runCommandWithSudo(sudoPassword, "nmcli", "-t", "-f", "SSID", "dev", "wifi", "list")
	if err != nil {
		fmt.Printf("Ошибка получения списка сетей: %v\nВывод: %s\n", err, output)
		os.Exit(1)
	}

	networks := strings.Split(output, "\n")
	if len(networks) == 0 {
		fmt.Println("Не удалось найти доступные сети. Проверьте, что ваш беспроводной адаптер включен и находится в зоне действия сети.")
		os.Exit(1)
	}

	// Выводим список сетей для выбора пользователем
	fmt.Println("Доступные сети:")
	for i, network := range networks {
		if network != "" {
			fmt.Printf("%d) %s\n", i, network)
		}
	}

	// Просим пользователя выбрать сеть
	fmt.Println("Выберите номер сети:")
	networkNumberStr, _ := reader.ReadString('\n')
	networkNumberStr = strings.TrimSpace(networkNumberStr)
	networkNumber, err := strconv.Atoi(networkNumberStr)
	if err != nil || networkNumber < 0 || networkNumber >= len(networks) || networks[networkNumber] == "" {
		fmt.Println("Некорректный ввод. Пожалуйста, запустите программу снова и выберите правильный номер сети.")
		os.Exit(1)
	}

	selectedSSID := networks[networkNumber]

	// Просим пользователя ввести пароль для сети
	fmt.Printf("Введите пароль для сети '%s':\n", selectedSSID)
	networkPassword, _ := reader.ReadString('\n')
	networkPassword = strings.TrimSpace(networkPassword)

	// Подключаемся к выбранной сети
	output, err = runCommandWithSudo(sudoPassword, "nmcli", "dev", "wifi", "connect", selectedSSID, "password", networkPassword, "ifname", wifiInterface)
	if err != nil {
		fmt.Printf("Ошибка подключения к сети: %v\nВывод: %s\n", err, output)
		os.Exit(1)
	}

	fmt.Printf("Подключение к сети '%s' выполнено.\n", selectedSSID)
}
