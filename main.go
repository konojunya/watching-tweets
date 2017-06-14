package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"log"
	"net/url"
	. "os"
	"regexp"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func main() {

	users := readJson()

	refresh(users)

	api := getTwitterApi()

	ch := make(chan string)
	for _, user := range users {
		go streamTweet(api, user.Name, user.Id, ch)
	}
	fmt.Println(<-ch)
}

func refresh(users []User) {

	if err := RemoveAll("../tweets"); err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		if err := MkdirAll("../tweets/"+user.Name, 0755); err != nil {
			log.Fatal(err)
		}
		fmt.Println("create dir tweets/" + user.Name)
	}

	fmt.Println("\nProject is refreshed!\n\n")

}

func readJson() []User {
	bytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		panic(err)
	}

	var users []User
	if err := json.Unmarshal(bytes, &users); err != nil {
		panic(err)
	}

	return users
}

func fileExport(path, text string) {

	file, err := Create("../tweets/" + path + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write(([]byte)(text))

}

func streamTweet(api *anaconda.TwitterApi, username, idstr string, ch chan string) {

	v := url.Values{}
	v.Set("follow", idstr)

	stream := api.PublicStreamFilter(v)

	fmt.Println("started stream: @" + username)

	for {
		x := <-stream.C
		switch tweet := x.(type) {
		case anaconda.Tweet:
			rep := regexp.MustCompile(`^@.*\s`)
			if !rep.MatchString(tweet.Text) { // @???で誰かに向けたツイート以外を取得
				fileExport(username+"/"+tweet.IdStr, tweet.Text)
				fmt.Println(tweet.Text + "from @" + username)
				fmt.Println("--------")
			}
		default:
			fmt.Printf("unkown type(%T): %v \n", x, x)
			ch <- "end"
		}
	}
}

func getTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(Getenv("TWITTER_ACCESS_TOKEN"), Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	return api
}
