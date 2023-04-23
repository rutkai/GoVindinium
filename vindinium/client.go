package vindinium

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	MoveTimeout  = 15
	StartTimeout = 0
)

type Client struct {
	Server    string
	Key       string
	Mode      string
	BotName   string
	Turns     string
	RandomMap bool
	Debug     bool
	Bot       Bot
	State     *State
	Url       string
}

func NewClient(server, key, mode, botName, turns string, randomMap bool, debug bool) (client *Client) {
	client = &Client{
		Server:    server,
		Key:       key,
		Mode:      mode,
		BotName:   botName,
		Turns:     turns,
		RandomMap: randomMap,
		Debug:     debug,
	}
	client.Setup()
	return
}

func (c *Client) Setup() {
	c.Url = c.Server + "/api/" + c.Mode
	switch c.BotName {
	case "fighter":
		c.Bot = &FighterBot{}
	case "introvert":
		c.Bot = &IntrovertBot{}
	default:
		c.Bot = &RandomBot{}
	}
}

func (c *Client) finished() bool {
	return c.State.Game.Finished
}

func (c *Client) move(dir Direction) error {
	values := make(url.Values)
	values.Set("dir", string(dir))
	return c.post(c.State.PlayUrl, values, MoveTimeout)
}

func (c *Client) post(uri string, values url.Values, seconds int) error {
	if c.Debug {
		fmt.Printf("Making request to to: %s\n", uri)
	}
	timeout := time.Duration(seconds) * time.Second
	dial := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if timeout == 0 {
			return net.Dial(network, addr)
		}
		return net.DialTimeout(network, addr, timeout)
	}

	transport := http.Transport{DialContext: dial}
	client := http.Client{Transport: &transport}

	response, err := client.PostForm(uri, values)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	data, _ := io.ReadAll(response.Body)

	if response.StatusCode >= 500 {
		return errors.New(fmt.Sprintf("Server responded with %s", response.Status))
	} else if response.StatusCode >= 400 {
		return errors.New(fmt.Sprintf("Request error: %s", string(data[:])))
	}

	if err := json.Unmarshal(data, &c.State); err != nil {
		return err
	}
	c.State.Init()

	if c.Debug {
		fmt.Printf("Setting data to:\n%s\n", string(data))
	}

	return nil
}

func (c *Client) Start() error {
	values := make(url.Values)
	values.Set("key", c.Key)
	if c.Mode == "training" {
		values.Set("turns", c.Turns)
		if !c.RandomMap {
			values.Set("map", "m1")
		}
	}

	fmt.Printf("Connecting to %s and waiting for other players to join...\n", c.Server)
	return c.post(c.Url, values, StartTimeout)
}

func (c *Client) Play() error {
	fmt.Printf("Playing at: %s\n", c.State.ViewUrl)
	move := 1
	for c.State.Game.Finished == false {
		fmt.Printf("\rMaking move: %d", move)

		if c.Debug {
			fmt.Printf("\nclient: %+v\n", c)
			fmt.Printf("bot: %+v\n", c.Bot)
			fmt.Printf("state: %+v\n", c.State)
		}

		dir := c.Bot.Move(c.State)
		if err := c.move(dir); err != nil {
			return err
		}

		move++
	}

	fmt.Println("\nFinished game.")
	return nil
}
