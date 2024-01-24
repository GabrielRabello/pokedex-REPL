package main

import (
	"bufio"
	// "encoding/json"
	"errors"
	"fmt"
	"github.com/GabrielRabello/pokedex-REPL/dto"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	run         func(*config) error
}

type config struct {
	offsetMul int
}

var locations []dto.Location

const cliName string = "Pokedex"
const apiURL string = "https://pokeapi.co/api/v2/"
const apiResource string = "location-area/"

// func commands() func() map[string]cliCommand {
// 	pg := config{offsetMul: 1}
// 	return func() map[string]cliCommand {
// 		return map[string]cliCommand{
// 			"help": {
// 				name:        "help",
// 				description: "Displays a help message",
// 				run:         commandHelp,
// 			},
// 			"exit": {
// 				name:        "exit",
// 				description: "Exit the Pokedex",
// 				run:         commandExit,
// 			},
// 			"clear": {
// 				name:        "clear",
// 				description: "Clear the screen",
// 				run:         commandClearScreen,
// 			},
// 			"map": {
// 				name:        "map",
// 				description: "Displays 20 locations per page",
// 				run:         commandMap,
// 			},
// 		}
// 	}
// }

var commands = map[string]cliCommand{
	"help": {
		name:        "help",
		description: "Displays a help message",
		run:         commandHelp,
	},
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		run:         commandExit,
	},
	"clear": {
		name:        "clear",
		description: "Clear the screen",
		run:         commandClearScreen,
	},
	"map": {
		name:        "map",
		description: "Displays 20 locations per page",
		run:         commandMap,
	},
}

func commandHelp(_ *config) error {
	const str = "\n" +
		"Welcome to the Pokedex!\n\n" +
		"Usage:\n" +
		"help: Displays a help message\n" +
		"exit: Exit the Pokedex" +
		"clear: Clear the screen"
	_, err := fmt.Println(str)
	return err
}

func commandExit(_ *config) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}

func commandClearScreen(_ *config) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func commandMap(conf *config) error {
	url := apiURL + apiResource + "?offset=" + fmt.Sprint(conf.offsetMul*20) + "&limit=20"
	_, err := getRequest(url)
	if err != nil {
		return err
	}
	return err
}

func getRequest(url string) ([]byte, error) {
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

func handleInput(scanner *bufio.Scanner) (cliCommand, error) {
	scanner.Scan()
	input := scanner.Text()
	if err := scanner.Err(); err != nil {
		return cliCommand{}, errors.New("Erro: " + err.Error())
	}
	input = sanitizeInput(input)
	val, ok := commands[input]
	if !ok {
		return cliCommand{}, errors.New("comando inválido")
	}
	return val, nil
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)
	conf := config{offsetMul: 1}
	res, err := http.Get(apiURL + apiResource + "?offset=" + fmt.Sprint(conf.offsetMul*20) + "&limit=20")
	if err != nil {
		fmt.Print(err.Error())
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s", body)

	// for {
	// 	printPrompt()
	// 	command, err := handleInput(sc)
	// 	if err != nil {
	// 		fmt.Fprintln(os.Stderr, "\nErro: ", err)
	// 		continue
	// 	}
	// 	err = command.run(&conf)
	// 	if err != nil {
	// 		fmt.Fprintln(os.Stderr, "\nErro: ", err)
	// 	}
	// }
}
