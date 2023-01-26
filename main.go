package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html/charset"
)

const FEED_URI = "https://somafm.com/channels.xml"
const SERVER_ID = 6
const RESULT_FILE_NAME = "somafm.pls"

type Channels struct {
	XMLName xml.Name `xml:"channels"`
	Channel []struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Id          string `xml:"id,attr"`
	} `xml:"channel"`
}

func main() {
	var err error
	var channels Channels

	data, err := loadXmlFromServer()
	if err != nil {
		log.Fatalln(err)
	}

	err = parseXml(data, err, channels)
	if err != nil {
		log.Fatalln(err)
	}

	err = fileGenerator(channels)
	if err != nil {
		log.Fatalln(err)
	}

}

func fileGenerator(channels Channels) error {
	fd, err := os.OpenFile(RESULT_FILE_NAME, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer fd.Close()

	i := 0
	fd.WriteString("[playlist]\n")
	for _, channel := range channels.Channel {
		i++

		//TODO extract url pattern
		fd.WriteString(fmt.Sprintf("File%d=https://ice%d.somafm.com/%s-128-aac\n", i, SERVER_ID, channel.Id))
		fd.WriteString(fmt.Sprintf("Title%d=%s: %s\n", i, channel.Title, channel.Description))
		fd.WriteString(fmt.Sprintf("Length%d=-1\n", i))
	}

	fd.WriteString(fmt.Sprintf("NumberOfEntries=%d\n", i))

	return nil
}

func parseXml(data []byte, err error, channels Channels) error {
	reader := bytes.NewReader([]byte(data))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	return decoder.Decode(&channels)
}

func loadXmlFromServer() ([]byte, error) {
	response, err := http.Get(FEED_URI)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
