package main

import (
	"container/list"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var playlist list.List

const (
	PAUSED = true
	PLAY   = false
)

type Song struct {
	name     string
	duration int
}

func addSong() {

	var duration int
	var name string
	fmt.Println("Type name of song")
	fmt.Scanf("%s", &name)
	fmt.Println("Type duration of song")
	fmt.Scanf("%d", &duration)
	song := Song{name: name, duration: duration}

	playlist.PushBack(song)
}

func PrevSong(elem *list.Element, playlist *list.List) *list.Element {

	if elem.Prev() == nil {
		elem = (*playlist).Back()
	} else {
		elem = elem.Prev()
	}

	return elem
}

func NextSong(elem *list.Element, playlist *list.List) *list.Element {

	if elem.Next() == nil {
		elem = (*playlist).Front()
	} else {
		elem = elem.Next()
	}

	return elem
}

func Play(song Song, cont chan bool, wg *sync.WaitGroup, ctx context.Context) {

	defer wg.Done()

	fmt.Println("Playing song:", song.name)
	for i := 0; i < song.duration; i++ {
		select {
		case input := <-cont:
			if input == PAUSED {
				select {
				case <-time.After(time.Second):
					fmt.Println("p")
				}
			} else {
				continue
			}
		case <-time.After(time.Second):
			fmt.Println(i)
		case <-ctx.Done():
			return
		}
	}
}

func ReadInput(chRead chan string) {
	var somestring string
	for {
		fmt.Scanf("%s\n", &somestring)
		fmt.Println("readInput: ", somestring)
		chRead <- somestring
	}
}

func main() {
	continueChannel := make(chan bool, 1)
	readChannel := make(chan string)
	gracefulShutdown := make(chan os.Signal, 4)

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	playlist := list.New()
	playlist.PushBack(Song{name: "ad", duration: 2})
	playlist.PushBack(Song{name: "as", duration: 3})
	elem := playlist.Front()
	song := elem.Value.(Song)

	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go Play(song, continueChannel, wg, ctx)
	wg.Wait()
	go ReadInput(readChannel)
	for {
		select {
		case input := <-readChannel:
			switch input {
			case "pause":
				continueChannel <- PAUSED
			case "next":
				wg.Add(1)
				elem = NextSong(elem, playlist)
				go Play(elem.Value.(Song), continueChannel, wg, ctx)
			case "prev":
				wg.Add(1)
				elem = PrevSong(elem, playlist)
				go Play(elem.Value.(Song), continueChannel, wg, ctx)
			case "exit":
				cancel()
				return
			case "addSong":
				addSong()
			case "play":
				continueChannel <- PLAY
			}
		case <-gracefulShutdown:
			fmt.Println("Done")
			cancel()
			return
		}
	}
}
