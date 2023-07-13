package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const FEED_URI = "https://somafm.com/channels.xml"
const SERVER_ID = 6
const RESULT_FILE_NAME = "somafm.pls"

type Channels struct {
	XMLName xml.Name  `xml:"channels"`
	Channel []Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Id          string `xml:"id,attr"`
}

func main() {
	var err error
	var channels Channels

	xmlDocument, err := loadXmlFromServer()
	if err != nil {
		log.Fatalln(err)
	}

	err = parseXml(xmlDocument, err, &channels)
	if err != nil {
		log.Fatalln(err)
	}

	err = fileGenerator(&channels)
	if err != nil {
		log.Fatalln(err)
	}
}

func fileGenerator(channels *Channels) error {
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

func parseXml(xmlDocument string, err error, channels *Channels) error {
	//HACK for charset
	data := []byte(strings.Replace(xmlDocument, "encoding=\"ISO-8859-1\"", "", 1))

	return xml.Unmarshal(data, &channels)
}

func loadXmlFromServer() (string, error) {
	response, err := http.Get(FEED_URI)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
