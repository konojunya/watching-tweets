package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	. "os"
	"regexp"
)

func main() {

	api := getTwitterApi()

	var ids map[string]string = make(map[string]string)
	ids["ms_rinna"] = "3274075003"
	ids["kawamina_happy"] = "726052257121214468"

	ch := make(chan string)
	for id := range ids {
		go streamTweet(api, id, ids[id], ch)
	}
	fmt.Println(<-ch)
}

func fileExport(path, text string) {

	file, err := Create("tweets/" + path + ".txt")
	if err != nil {
		panic(err)
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
				fileExport(username+"/"+idstr, tweet.Text)
				fmt.Println(tweet.Text + "from @" + username)
				fmt.Println("--------")
			}
		default:
			fmt.Println("unkown type(%T): %v \n", x, x)
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
