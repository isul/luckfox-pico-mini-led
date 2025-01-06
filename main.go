package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var gpioPin int

	// GPIO 핀 번호 입력 받기, LuckFox Pico Mini A/B의 경우 34
	fmt.Print("Please enter the GPIO pin number: ")
	_, err := fmt.Scan(&gpioPin)
	if err != nil {
		fmt.Printf("Failed to read input: %v\n", err)
		return
	}

	// GPIO export 파일 열기
	if err := writeToFile("/sys/class/gpio/export", strconv.Itoa(gpioPin)); err != nil {
		fmt.Printf("Failed to export GPIO pin: %v\n", err)
		return
	}

	// GPIO direction 설정
	directionPath := fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioPin)
	if err := writeToFile(directionPath, "out"); err != nil {
		fmt.Printf("Failed to set GPIO direction: %v\n", err)
		return
	}

	// GPIO value 파일 제어
	valuePath := fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioPin)
	for i := 0; i < 3; i++ {
		if err := writeToFile(valuePath, "1"); err != nil {
			fmt.Printf("Failed to write value: %v\n", err)
			return
		}

		// 값 읽기
		value, err := readFromFile(valuePath)
		if err != nil {
			fmt.Printf("Failed to read value: %v\n", err)
			return
		}
		fmt.Printf("value: %s", value)
		time.Sleep(1 * time.Second)

		if err := writeToFile(valuePath, "0"); err != nil {
			fmt.Printf("Failed to write value: %v\n", err)
			return
		}

		// 값 읽기
		value, err = readFromFile(valuePath)
		if err != nil {
			fmt.Printf("Failed to read value: %v\n", err)
			return
		}
		fmt.Printf("value: %s", value)
		time.Sleep(1 * time.Second)
	}

	// GPIO unexport 파일 열기
	if err := writeToFile("/sys/class/gpio/unexport", strconv.Itoa(gpioPin)); err != nil {
		fmt.Printf("Failed to unexport GPIO pin: %v\n", err)
		return
	}
}

func writeToFile(path, data string) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(data)
	if err != nil {
		return err
	}
	return writer.Flush()
}

func readFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := reader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		return "", err
	}
	return data, nil
}
