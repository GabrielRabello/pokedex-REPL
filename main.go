package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GabrielRabello/pokedex-REPL/dto"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ICommand interface {
	Run() error
}

type Command struct {
	name        string
	description string
	run         func() error
}

func (c *Command) Run() error {
	return c.run()
}

type CommandBatch struct {
	name        string
	description string
	run         func(*dto.BatchResult) error
	result      dto.BatchResult
}

func (c *CommandBatch) Run() error {
	return c.run(&c.result)
}

const cliName string = "Pokedex"
const apiURL string = "https://pokeapi.co/api/v2/"
const apiResource string = "location-area/"

var commands = map[string]ICommand{
	"help": &Command{
		name:        "help",
		description: "Displays a help message",
		run:         commandHelp,
	},
	"exit": &Command{
		name:        "exit",
		description: "Exit the Pokedex",
		run:         commandExit,
	},
	"clear": &Command{
		name:        "clear",
		description: "Clear the screen",
		run:         commandClearScreen,
	},
	"map": &CommandBatch{
		name:        "map",
		description: "Displays 20 locations per page",
		run:         commandMap,
		result:      dto.BatchResult{},
	},
}

func commandHelp() error {
	const str = "\n" +
		"Welcome to the Pokedex!\n\n" +
		"Usage:\n" +
		"help: Displays a help message\n" +
		"exit: Exit the Pokedex" +
		"clear: Clear the screen"
	_, err := fmt.Println(str)
	return err
}

func commandExit() error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}

func commandClearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func commandMap(conf *dto.BatchResult) error {
	if conf.Next == nil {
		s := apiURL + apiResource
		conf.Next = &s
	}
	body, err := GET(*conf.Next)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &conf)
	if err != nil {
		return fmt.Errorf("error unmarshalling response body: %w", err)
	}
	for _, v := range conf.Results {
		fmt.Println(v.Name)
	}
	return nil
}

func GET(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return body, nil
}
func printPrompt() {
	fmt.Print(cliName, "> ")
}

func sanitizeInput(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	return str
}

func handleInput(scanner *bufio.Scanner) (ICommand, error) {
	scanner.Scan()
	input := scanner.Text()
	if err := scanner.Err(); err != nil {
		return nil, errors.New("Erro: " + err.Error())
	}
	input = sanitizeInput(input)
	val, ok := commands[input]
	if !ok {
		return nil, errors.New("comando inválido")
	}
	return val, nil
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)
	for {
		printPrompt()
		command, err := handleInput(sc)
		if err != nil {
			fmt.Fprintln(os.Stderr, "\nErro: ", err)
			continue
		}
		err = command.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "\nErro: ", err)
		}
	}
}
