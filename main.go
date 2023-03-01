package main

import (
	"container/list"
	"fmt"
	"time"
)

var playlist *list.List

type tSong struct {
	duration int
}

func play(song tSong, chPause chan string) {
	fmt.Println("Play Song")
	for i := 0; i < song.duration; i++ {
		select {
		case input, ok := <-chPause:
			{
				if ok {
					if input == "pause" {
						chPause = nil
						select {
						case <-chPause:
						}
					}
					if input == "stop" {
						chPause = nil

					}
				}
			}
		default:
			time.After(time.Second)
		}
	}
}

func pause() <-chan bool {
	ch := make(chan bool, 1)

	return ch
}

func addSong(song tSong) {
	playlist.PushBack(song)
}

func readInput(chRead chan string) {
	var somestring string
	for {
		fmt.Scanf("%s\n", &somestring)
		fmt.Println("readInput: ", somestring)
		chRead <- somestring
		fmt.Println("readInput: sent to chRead")
	}
}

func main() {
	chPause := make(chan string)
	chRead := make(chan string)

	playlist = list.New()

	go readInput(chRead)
	if playlist.Len() == 0 {
		var duration int
		fmt.Println("Please add new Song by typing length of song: addSong")
		fmt.Scanln(&duration)
		addSong(tSong{duration: duration})
	}
	if playlist.Len() != 0 {
		go play(playlist.Front().Value.(tSong), chPause)
		for e := playlist.Front(); e != nil; e = e.Next() {
			select {
			case input, ok := <-chRead:
				if ok {
					if input == "addSong" {
						var duration int
						fmt.Println("Type duration of song")
						fmt.Scanln(&duration)
						addSong(tSong{duration: duration})
					}
					if input == "pause" {
						pause()
					}
					if input == "prev" {
						e.Prev()
					}
					if input == "next" {
						e.Next()
					}
				}
			}
		}
	}
}
