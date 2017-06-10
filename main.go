package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	. "os"
)

func main() {

	streamTweet()

}

func fileExport(filename string,output string) {

	file, err := Create("tweets/rinna/"+filename+".txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(([]byte)(output))

}

func streamTweet() {
	api := GetTwitterApi()

	v := url.Values{}
	v.Set("follow", "3274075003")

	stream := api.PublicStreamFilter(v)

	fmt.Println("streaming start!")

	for {
		x := <-stream.C
		switch tweet := x.(type) {
		case anaconda.Tweet:
			fileExport(tweet.IdStr,tweet.Text + "\n")
			fmt.Println(tweet.Text)
			fmt.Println("--------")
		case anaconda.StatusDeletionNotice:
			// pass
		default:
			fmt.Println("unkown type(%T): %v \n", x, x)
		}
	}
}

func GetTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(Getenv("TWITTER_ACCESS_TOKEN"), Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	return api
}
