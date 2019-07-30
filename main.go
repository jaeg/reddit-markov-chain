package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jaeg/markov-chain-text-generator/mctg"
	geddit "github.com/jzelinskie/geddit"
)

var processedPosts map[string]string

type Settings struct {
	Username string
	Password string
}

func main() {
	processedPosts = make(map[string]string)

	m := mctg.New(1)
	go huntReddit(m)

	i := 0
	for {
		fmt.Println(i, m.GenerateSentence())
		i++
		time.Sleep(time.Second)
	}

}

func huntReddit(m *mctg.MCTG) {
	settings, _ := readSettingsFile("settings.json")
	session, _ := geddit.NewLoginSession(
		settings.Username,
		settings.Password,
		"gedditAgent v1",
	)

	subOpts := geddit.ListingOptions{
		Limit: 100,
	}
	for {
		ds := ""
		subreddit := "all"
		submissions, _ := session.SubredditComments(subreddit, subOpts)

		for _, s := range submissions {
			if processedPosts[s.Body] == "" {
				ds += s.Body + " "
				processedPosts[s.Body] = subreddit
				fmt.Println("Added", s.Body)
			}
		}

		m.ParseCorpusFromString(ds, false)
		time.Sleep(time.Second * 5)
	}
}

func readSettingsFile(path string) (settings Settings, err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	b, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(b, &settings)
	return
}
