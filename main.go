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
	Username   string
	Password   string
	Aggression int
	Subreddits []string
}

func main() {
	processedPosts = make(map[string]string)

	m := mctg.New(1)
	go huntReddit(m)

	i := 0
	for {
		fmt.Println("Posts Handled: ", len(processedPosts))
		fmt.Println(i, ": ", m.GenerateSentence())
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

	t := time.Now()
	ts := t.Format("2006-01-02 15-04-05")
	fmt.Println(ts)
	f, err := os.OpenFile("corpus-"+ts+".txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	subOpts := geddit.ListingOptions{
		Limit: 1000,
	}
	for {
		for _, subreddit := range settings.Subreddits {
			ds := ""
			submissions, _ := session.SubredditComments(subreddit, subOpts)

			for _, s := range submissions {
				if processedPosts[s.Body] == "" {
					ds += s.Body + " "
					processedPosts[s.Body] = subreddit
				}
			}

			if len(ds) > 0 {
				m.ParseCorpusFromString(ds, false)
				if _, err = f.WriteString(ds); err != nil {
					panic(err)
				}
			}

			time.Sleep(time.Second * time.Duration(settings.Aggression))
		}
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
