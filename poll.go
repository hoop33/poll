package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/codegangsta/cli"
)

type PollResponse struct {
	Links struct {
		Poll struct {
			Href string
		}
	}
}

var verbose = false

func main() {

	app := cli.NewApp()
	app.Name = "poll"
	app.Version = "0.0.1"
	app.Usage = "Call an API that returns a poll URL, and then poll that URL until a 200 or error"
	app.Authors = []cli.Author{
		{Name: "Rob Warner", Email: "rwarner@grailbox.com"},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "user, u",
			EnvVar: "REMOTE_USER",
			Usage:  "authorized user",
		},
		cli.StringFlag{
			Name:   "baseurl, b",
			EnvVar: "POLL_BASE_URL",
			Usage:  "the base URL for the initial request",
		},
		cli.IntFlag{
			Name:   "interval, i",
			EnvVar: "POLL_INTERVAL",
			Usage:  "poll interval, in seconds",
		},
		cli.IntFlag{
			Name:   "maxtries, m",
			EnvVar: "POLL_MAX_TRIES",
			Usage:  "the maximum number of times to poll",
		},
		cli.BoolFlag{
			Name:   "verbose, V",
			EnvVar: "POLL_VERBOSE",
			Usage:  "run in verbose mode",
		},
	}

	app.Action = func(c *cli.Context) {
		user := c.String("user")
		baseurl := c.String("baseurl")
		if baseurl == "" {
			baseurl = "http://localhost:8280/v1/"
		}
		interval := c.Int("interval")
		if interval == 0 {
			interval = 1
		}
		maxtries := c.Int("maxtries")
		if maxtries == 0 {
			maxtries = 10
		}
		verbose = c.Bool("verbose")

		if len(c.Args()) > 0 {
			url := fmt.Sprint(baseurl, c.Args()[0])
			sc, location, json, err := getUrl(user, url)
			handleError(err)

			if sc == 202 {
				var pollurl string
				if location != nil {
					pollurl = location.String()
					log(fmt.Sprint("Using location header ", pollurl))
				} else {
					pollurl, err = parse(json)
					handleError(err)
					log(fmt.Sprint("Using URL from JSON ", pollurl))
				}

				s := "seconds"
				if interval == 1 {
					s = "second"
				}
				for i := 0; i < maxtries; i++ {
					log(fmt.Sprint("Sleeping ", interval, " ", s, "..."))
					time.Sleep(time.Second * time.Duration(interval))
					log(fmt.Sprint("Poll #", (i + 1), "..."))
					sc, _, json, err = getUrl(user, pollurl)
					handleError(err)
					if sc != 202 {
						break
					}
				}
			}
      log(fmt.Sprint("Status Code:", sc, "\n"))
			fmt.Println(string(json))
		} else {
			fmt.Println("Missing URL parameter")
		}
	}

	app.Run(os.Args)
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func log(message string) {
	if verbose {
		fmt.Println(message)
	}
}

func getUrl(user string, url string) (int, *url.URL, []byte, error) {
	log(fmt.Sprint("Getting url ", url, " as user ", user))

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("RemoteUser", user)
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}
	location, _ := resp.Location()
	return resp.StatusCode, location, payload, nil
}

func parse(b []byte) (string, error) {
	var pollResponse PollResponse
	err := json.Unmarshal(b, &pollResponse)
	handleError(err)

	if pollResponse.Links.Poll.Href != "" {
		return pollResponse.Links.Poll.Href, nil
	}
	return "", errors.New(fmt.Sprint("No poll url in response:\n\n", string(b)))
}
