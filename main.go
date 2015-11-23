package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	fs := http.FileServer(assetFS())

	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.HandleFunc("/admin", dashboard)

	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

type Stats struct {
	Queries      int
	AdsBlocked   int
	NumOfDomains int
	TodaysAds    string
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	tstr, _ := Asset("templates/index.tmpl")
	tmpl, er := template.New("inx").Parse(string(tstr[:]))
	if er != nil {
		log.Println(er)
	}

	stats := Stats{}
	stats.AdsBlocked, er = adsBlockedToday()

	if er != nil {
		log.Println(er)
	}

	stats.Queries, er = queries()

	if er != nil {
		log.Println(er)
	}

	stats.NumOfDomains, er = numDomains()

	if er != nil {
		log.Println(er)
	}

	if stats.Queries > 0 {
		x := float64(stats.AdsBlocked) / float64(stats.Queries) * 100
		stats.TodaysAds = fmt.Sprintf("%0.2f", x)
	}

	err := tmpl.Execute(w, stats)
	if err != nil {
		log.Println(err)
	}
}

func adsBlockedToday() (int, error) {
	// cat /var/log/pihole.log
	cat := exec.Command("cat", "/var/log/pihole.log")

	// awk  '/\/etc\/pihole\/gravity.list/ {print $6}'
	awk := exec.Command("awk", `/\/etc\/pihole\/gravity.list/ {print $6}`)

	wc := exec.Command("wc", "-l")

	// Run the pipeline
	output, stderr, err := Pipeline(cat, awk, wc)

	if len(stderr) > 0 {
		log.Panicln(stderr)
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(output[:])))
}

func numDomains() (int, error) {
	// wc -l /etc/pihole/gravity.list | awk '{print $1}'

	wc := exec.Command("wc", "-l", "/etc/pihole/gravity.list")

	awk := exec.Command("awk", `{print $1}`)

	// Run the pipeline
	output, stderr, err := Pipeline(wc, awk)

	if len(stderr) > 0 {
		log.Panicln(stderr)
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(output[:])))
}

func queries() (int, error) {
	// cat /var/log/pihole.log | awk '/query/ {print $6}' | wc -l

	// cat /var/log/pihole.log
	cat := exec.Command("cat", "/var/log/pihole.log")

	// awk '/query/ {print $6}'
	awk := exec.Command("awk", `/query/ {print $6}`)

	wc := exec.Command("wc", "-l")

	// Run the pipeline
	output, stderr, err := Pipeline(cat, awk, wc)

	if len(stderr) > 0 {
		log.Panicln(stderr)
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(output[:])))
}

// Pipeline strings together the given exec.Cmd commands in a similar fashion
// to the Unix pipeline.  Each command's standard output is connected to the
// standard input of the next command, and the output of the final command in
// the pipeline is returned, along with the collected standard error of all
// commands and the first error found (if any).
//
// To provide input to the pipeline, assign an io.Reader to the first's Stdin.
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}
