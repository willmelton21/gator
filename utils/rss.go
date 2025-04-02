package rss

import (
  "fmt"
  "net/http"
  "io"
  "encoding/xml"
  "html"
  "context"
)
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
    
    var feed RSSFeed 

    //client := &http.Client{}

    req, err := http.NewRequestWithContext(ctx,"GET",feedURL, nil)
      if err != nil {
        return nil, fmt.Errorf("error getting http request %w",err)
      }
      
    req.Header.Set("User-Agent","Gator")
    resp, err := http.DefaultClient.Do(req)
     if err != nil {
        return nil, fmt.Errorf("error getting http response from Do %w",err)
      }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("http returned not ok")
        }

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
       return nil, fmt.Errorf("error reading http body %w",err)
     }
       
    err = xml.Unmarshal(bodyBytes, &feed)
    if err != nil {
        return nil, fmt.Errorf("Error unmarshalling data to struct %w",err)

    }

  feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
  feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

  for i := 0; i < len(feed.Channel.Item); i++ {
        feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
        feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)

    }

  return &feed, nil 
}
