package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/playwright-community/playwright-go"
)

func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v does not equal %v", actual, expected)
	}
}

func setup_webcalc() {
	// Get playwright
	err := playwright.Install()
	if err != nil {
		log.Fatal("Can install playwright or browsers", err)
	}

	http.HandleFunc("/calculate", calculatorHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func setup_playwright() playwright.Page {
	// Start up playwright browser
	pw_options := playwright.RunOptions{
		SkipInstallBrowsers: false,
		Verbose:             true,
	}
	pw, err := playwright.Run(&pw_options)
	assertErrorToNilf("could not launch playwright: %w", err)
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Timeout:  playwright.Float(2000),
	})
	assertErrorToNilf("could not launch Chromium: %w", err)
	context, err := browser.NewContext()
	assertErrorToNilf("could not create context: %w", err)
	page, err := context.NewPage()
	assertErrorToNilf("could not create page: %w", err)
	return page
}

func fillAndSubmit(page playwright.Page, first, second, operation string) {
	err := page.Locator("input[name=first]").Fill(first)
	assertErrorToNilf("could not fill first input: %w", err)
	err = page.Locator("input[name=second]").Fill(second)
	assertErrorToNilf("could not fill second input: %w", err)
	operationOptions := playwright.SelectOptionValues{Values: &[]string{operation}}
	_, err = page.Locator("select[name=operation]").SelectOption(operationOptions)
	assertErrorToNilf(fmt.Sprintf("could not select %s: %%w", operation), err)
	err = page.Locator("button[type=submit]").Click()
	assertErrorToNilf("could not click submit: %w", err)
}

func gotoFrontPage(page playwright.Page) {
	_, err := page.Goto("http://localhost:8080")
	assertErrorToNilf("could not goto: %w", err)
}

func findResult(page playwright.Page) string {
	result, err := page.Locator("div#result").TextContent()
	assertErrorToNilf("could not get result: %w", err)
	return result
}

func Test_Operations(t *testing.T) {
	go setup_webcalc()
	page := setup_playwright()

	var TT = []struct {
		Name      string
		First     string
		Second    string
		Operation string
		Expexted  string
	}{
		{"Adding", "5", "3", "add", "8"},
		{"Multiplying", "2", "8", "multiply", "16"},
		{"Division", "12", "3", "divide", "4"},
		{"Subtraction", "10", "7", "subtract", "3"},
	}

	for _, tt := range TT {
		gotoFrontPage(page)
		fmt.Println("Testing", tt.Name)
		fillAndSubmit(page, tt.First, tt.Second, tt.Operation)
		result := findResult(page)
		assertEqual(t, tt.Expexted, result)
	}
}
