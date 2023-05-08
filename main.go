package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var sb strings.Builder
var hardestCardMap = map[string]int{}

func main() {

	cardCollection := CardCollection{
		ByCards:       map[string]string{},
		ByDefinitions: map[string]string{},
	}

	if len(os.Args) > 3 {
		log.Fatal("Error! Maximum count of arguments <= 3.")
	}

	importInput := flag.String("import_from", "nothing.txt", "import from file")
	exportInput := flag.String("export_to", "nothing.txt", "export to file")

	flag.Parse()

	if *importInput != "nothing.txt" {
		cardCollection.ImportCardsFromLogs(*importInput)
	}
	if *exportInput != "nothing.txt" {
		cardCollection.ExportCardsFromLogsToDisk(*exportInput)
	}

	for {
		logPrintSave("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")
		action := ReadLine()
		logInputSave(action)

		switch action {
		case "add":
			cardCollection.AddCard()
		case "remove":
			cardCollection.RemoveCard()
		case "import":
			cardCollection.ImportCardsFromDisk()
		case "export":
			cardCollection.ExportCardsToDisk()
		case "ask":
			ask(&cardCollection)
		case "log":
			ExportLogsToDisk()
		case "hardest card":
			HardestCardOutput()
		case "reset stats":
			ResetStats()
		case "exit":
			logPrintSave("Bye bye!")
			return
		}
		fmt.Println()
	}

}

func ResetStats() {
	for k := range hardestCardMap {
		delete(hardestCardMap, k)
	}
	logPrintSave("Card statistics have been reset.")
}

func GetMaxCountOfTries() int {
	maxCount := 0
	for _, i := range hardestCardMap {
		if i > maxCount {
			maxCount = i
		}
	}
	return maxCount
}

func DelEverythingExceptMax() {
	maxCount := GetMaxCountOfTries()
	for key, i := range hardestCardMap {
		if i != maxCount {
			delete(hardestCardMap, key)
		}
	}
}

func HardestCardLenCheck() bool {
	keys := make([]interface{}, 0)
	for k, _ := range hardestCardMap {
		keys = append(keys, k)
	}

	if len(keys) == 0 {
		return true
	} else {
		return false
	}
}

func HardestCardOutput() {
	if HardestCardLenCheck() {
		logPrintSave("There are no cards with errors.")
	} else {
		DelEverythingExceptMax()
		var arrOfKeys []string
		for key, _ := range hardestCardMap {
			arrOfKeys = append(arrOfKeys, key)
		}
		hardestCardKeysOutput := strings.Join(arrOfKeys, ", ")
		logPrintSave(fmt.Sprintf("The hardest card is %q. You have %d errors answering it.", hardestCardKeysOutput, GetMaxCountOfTries()))
	}
}

func ExportLogsToDisk() {
	logPrintSave("File name:")
	filename := ReadLine()
	logInputSave(filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logPrintSave("The log has been saved.")
	file.WriteString(sb.String())
}

func logPrintSave(st string) {
	sb.WriteString(st)
	sb.WriteString("\n")
	fmt.Println(st)
}

func logInputSave(st string) {
	sb.WriteString(st)
	sb.WriteString("\n")
}

func ask(cardCollection *CardCollection) {
	logPrintSave("How many times to ask?")
	askTimes, _ := strconv.Atoi(ReadLine())
	logInputSave(strconv.Itoa(askTimes))
	askedTimes := 0
	for {
		for k, v := range cardCollection.ByCards {
			if askedTimes == askTimes {
				return
			}

			logPrintSave(fmt.Sprintf("Print the definition of \"%s\":\n", k))
			definitionOfUser := ReadLine()
			logInputSave(definitionOfUser)

			if v == definitionOfUser {
				fmt.Println("Correct!")
			} else if correctTermForUserDefinition, ok := cardCollection.ByDefinitions[definitionOfUser]; ok {
				logPrintSave(fmt.Sprintf("Wrong. The right answer is \"%s\", "+
					"but your definition is correct for \"%s\".\n", v, correctTermForUserDefinition))
				hardestCardMap[k] = hardestCardMap[k] + 1
			} else {
				hardestCardMap[k] = hardestCardMap[k] + 1
				logPrintSave(fmt.Sprintf("Wrong. The right answer is \"%s\"\n", v))
			}
			askedTimes++
		}
	}

}

type CardCollection struct {
	ByCards, ByDefinitions map[string]string
	Cards                  []CardDefinition
}

type CardDefinition struct {
	Card, Definition string
}

func (cardCollection *CardCollection) AddCard() {
	logPrintSave("The card")
	var card string
	for {
		card = ReadLine()
		logInputSave(card)
		if _, ok := cardCollection.ByCards[card]; ok {
			logPrintSave(fmt.Sprintf("The card \"%s\" already exists. Try again:\n", card))
		} else {
			break
		}
	}

	logPrintSave("The definition of the card")
	var definition string
	for {
		definition = ReadLine()
		logInputSave(definition)
		if _, ok := cardCollection.ByDefinitions[definition]; ok {
			logPrintSave(fmt.Sprintf("The definition \"%s\" already exists. Try again:\n", definition))
		} else {
			break
		}
	}
	cardCollection.ByCards[card] = definition
	cardCollection.ByDefinitions[definition] = card
	logPrintSave(fmt.Sprintf("The pair (\"%s\":\"%s\") has been added.\n", card, definition))
}

func (cardCollection *CardCollection) RemoveCard() {
	logPrintSave("Which card?")
	card := ReadLine()
	logInputSave(card)
	definition, cardExist := cardCollection.ByCards[card]
	if !cardExist {
		logPrintSave(fmt.Sprintf("Can't remove \"%s\": there is no such card.", card))
		return
	}
	delete(cardCollection.ByCards, card)
	delete(cardCollection.ByDefinitions, definition)
	logPrintSave("The card has been removed.")
}

func (cardCollection *CardCollection) ImportCardsFromDisk() {
	logPrintSave("File name:")
	filename := ReadLine()
	logInputSave(filename)

	file, err := os.Open(filename)
	if err != nil {
		logPrintSave("File not found.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	doubleCount := 0

	var card string
	var definition string
	for i := 0; scanner.Scan(); i++ {
		if i%2 == 0 {
			card = scanner.Text()
			logInputSave(card)
		} else {
			definition = scanner.Text()
			logInputSave(definition)
		}

		if len(card) > 0 && len(definition) > 0 {
			cardCollection.ByCards[card] = definition
			cardCollection.ByDefinitions[definition] = card
			card = ""
			definition = ""
		}
		doubleCount++
	}
	logPrintSave(fmt.Sprintf("%d cards have been loaded.\n", doubleCount/2))
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (cardCollection *CardCollection) ImportCardsFromLogs(st string) {
	filename := st

	file, err := os.Open(filename)
	if err != nil {
		logPrintSave("File not found.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	doubleCount := 0

	var card string
	var definition string
	for i := 0; scanner.Scan(); i++ {
		if i%2 == 0 {
			card = scanner.Text()
			logInputSave(card)
		} else {
			definition = scanner.Text()
			logInputSave(definition)
		}

		if len(card) > 0 && len(definition) > 0 {
			cardCollection.ByCards[card] = definition
			cardCollection.ByDefinitions[definition] = card
			card = ""
			definition = ""
		}
		doubleCount++
	}
	logPrintSave(fmt.Sprintf("%d cards have been loaded.\n", doubleCount/2))
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (cardCollection *CardCollection) ExportCardsToDisk() {
	logPrintSave("File name:")
	filename := ReadLine()
	logInputSave(filename)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for k, v := range cardCollection.ByCards {
		fmt.Fprintf(file, "%s\n%s\n", k, v)
	}
	logPrintSave(fmt.Sprintf("%d cards have been saved.\n", len(cardCollection.ByCards)))
}

func (cardCollection *CardCollection) ExportCardsFromLogsToDisk(st string) {
	filename := st

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for k, v := range cardCollection.ByCards {
		fmt.Fprintf(file, "%s\n%s\n", k, v)
	}
	logPrintSave(fmt.Sprintf("%d cards have been saved.\n", len(cardCollection.ByCards)))
}

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
