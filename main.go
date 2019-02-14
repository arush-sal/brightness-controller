package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const dirLocation = "/sys/class/backlight"

func main() {
	var inc = flag.Int("inc", 0, "Increase brightness by <percentage>")
	var desc = flag.Int("desc", 0, "Decrease brightness by <percentage>")

	flag.Parse()

	if f := flag.NFlag(); f != 1 {
		flag.Usage()
		os.Exit(126)
	}

	if *inc > 0 {
		err := increaseBrightness(*inc)
		errorCheck(err)
	}
	if *desc > 0 {
		err := decreaseBrightness(*desc)
		errorCheck(err)
	}
}

func increaseBrightness(step int) error {
	videoCard, err := findVideoCard()
	if err != nil {
		return err
	}

	for _, card := range videoCard {
		currentBrightnessLocation := dirLocation + "/" + card + "/brightness"
		maxBrightnessLocation := dirLocation + "/" + card + "/max_brightness"

		maxBrightness, err := ioutil.ReadFile(maxBrightnessLocation)
		if err != nil {
			return err
		}

		currentBrightness, err := ioutil.ReadFile(currentBrightnessLocation)
		if err != nil {
			return err
		}

		currentValue, err := strconv.Atoi(strings.TrimSuffix(string(currentBrightness), "\n"))
		errorCheck(err)
		maxValue, err := strconv.Atoi(strings.TrimSuffix(string(maxBrightness), "\n"))
		errorCheck(err)

		newValue := currentValue + ((maxValue / 100) * step)

		if newValue > maxValue {
			newValue = maxValue
		}

		errorCheck(ioutil.WriteFile(currentBrightnessLocation, []byte(strconv.Itoa(newValue)), 0644))
	}
	return nil
}

func decreaseBrightness(step int) error {
	videoCard, err := findVideoCard()
	if err != nil {
		return err
	}

	for _, card := range videoCard {
		currentBrightnessLocation := dirLocation + "/" + card + "/brightness"
		maxBrightnessLocation := dirLocation + "/" + card + "/max_brightness"

		maxBrightness, err := ioutil.ReadFile(maxBrightnessLocation)
		if err != nil {
			return err
		}

		currentBrightness, err := ioutil.ReadFile(currentBrightnessLocation)
		if err != nil {
			return err
		}

		currentValue, err := strconv.Atoi(strings.TrimSuffix(string(currentBrightness), "\n"))
		errorCheck(err)
		maxValue, err := strconv.Atoi(strings.TrimSuffix(string(maxBrightness), "\n"))
		errorCheck(err)

		newValue := currentValue - ((maxValue / 100) * step)

		if newValue < 0 {
			newValue = 1
		}

		errorCheck(ioutil.WriteFile(currentBrightnessLocation, []byte(strconv.Itoa(newValue)), 0644))
	}

	return nil
}

func findVideoCard() ([]string, error) {
	var ans int

	backlightDir, err := os.Open(dirLocation)
	errorCheck(err)
	defer backlightDir.Close()

	videoCard, err := backlightDir.Readdirnames(0)
	errorCheck(err)

	if len(videoCard) > 1 {
		fmt.Println("Multiple backlight devices found.", "Please select your device...")

		for index, card := range videoCard {
			fmt.Printf("[%d]\t%s\n", index, card)
		}
		fmt.Printf("[%d]\tall\n", len(videoCard))

		fmt.Scanf("%d", &ans)

		if ans == len(videoCard) {
			return videoCard, nil
		}
	}
	selectedVideoCard := []string{videoCard[ans]}
	return selectedVideoCard, nil
}

func errorCheck(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
