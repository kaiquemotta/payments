package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
)

var number1, number2, result int

// Função para o passo "Given I have a number"
func iHaveANumber(arg1 int) error {
	number1 = arg1
	return nil
}

// Função para o passo "And I have another number"
func iHaveAnotherNumber(arg1 int) error {
	number2 = arg1
	return nil
}

// Função para o passo "When I add the numbers"
func iAddTheNumbers() error {
	result = number1 + number2
	return nil
}

// Função para o passo "Then the result should be"
func theResultShouldBe(arg1 int) error {
	if result != arg1 {
		return fmt.Errorf("expected %d but got %d", arg1, result)
	}
	return nil
}

// Inicializa o cenário com os passos
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I have a number (\d+)$`, iHaveANumber)
	ctx.Step(`^I have another number (\d+)$`, iHaveAnotherNumber)
	ctx.Step(`^I add the numbers$`, iAddTheNumbers)
	ctx.Step(`^the result should be (\d+)$`, theResultShouldBe)
}

// Função principal para executar os testes
func TestMain(m *testing.M) {
	opts := godog.Options{
		Format: "pretty",                                 // Define o formato da saída
		Paths:  []string{"../features/addition.feature"}, // Caminho para o arquivo de feature
	}
	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
