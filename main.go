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
	ids["kinokoruumu416"] = "1662875580"

	streamTweet(api,ids)
}

func fileExport(filename string, output string) {

	file, err := Create("tweets/" + filename + ".txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(([]byte)(output))

}

func streamTweet(api *anaconda.TwitterApi,ids map[string]string) {

	for username := range ids{
		v := url.Values{}
		v.Set("follow", ids[username])

		stream := api.PublicStreamFilter(v)

		fmt.Println("stream started: @" + username)

		go func(){
			for {
				x := <-stream.C
				switch tweet := x.(type) {
				case anaconda.Tweet:
					rep := regexp.MustCompile(`^@.*\s`)
					if !rep.MatchString(tweet.Text) { // @???で誰かに向けたツイート以外を取得
						fileExport(username+"/"+tweet.IdStr, tweet.Text+"\n")
						fmt.Println(tweet.Text)
						fmt.Println("--------")
					}
				default:
					fmt.Println("unkown type(%T): %v \n", x, x)
				}
			}
		}()

	}

}

func getTwitterApi() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(Getenv("TWITTER_CONSUMER_SECRET"))
	api := anaconda.NewTwitterApi(Getenv("TWITTER_ACCESS_TOKEN"), Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	return api
}
