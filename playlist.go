package main

import (
	"container/list"
	"fmt"
	"os"
	"time"
)

var playlist *list.List

type tSong struct {
	duration int
}

func play(song tSong, chPause <-chan bool, chCont <-chan bool, chExit <-chan bool, chDone chan<- bool) {
	for i := 0; i < song.duration; i++ {
		fmt.Println(i)
		select {

		case <-chPause:
			select {
			case <-chCont:
				{
				}
			case <-chExit:
				return
			}

		case <-chExit:
			return

		default:
			{
				time.Sleep(time.Second)
			}

		}
	}

	chDone <- true
}

func addSong() {

	var duration int
	fmt.Println("Type duration of song")
	fmt.Scanf("%d", &duration)
	song := tSong{duration: duration}

	playlist.PushBack(song)
}

func readInput(chRead chan string) {
	var somestring string
	for {
		fmt.Scanf("%s\n", &somestring)
		fmt.Println("readInput: ", somestring)
		chRead <- somestring
	}
}

func main() {
	chPause := make(chan bool)
	chContinue := make(chan bool)
	chExit := make(chan bool)
	chRead := make(chan string)
	chDone := make(chan bool)

	playlist = list.New()

	if playlist.Len() == 0 {
		addSong()
	}
	elem := playlist.Front()
	song := elem.Value.(tSong)

	go readInput(chRead)
	go play(song, chPause, chContinue, chExit, chDone)
	if playlist.Len() != 0 {
		for {
			select {
			case input, ok := <-chRead:
				{
					if ok {
						switch input {
						case "next":
							{
								elem = elem.Next()

								if elem == nil {
									elem = playlist.Front()
									select {
									case <-chDone:
										{
										}
									case chExit <- true:
										{
										}
									}
									go play(elem.Value.(tSong), chPause, chContinue, chExit, chDone)
								}
								if elem.Next() != nil {
									elem = elem.Next()
									select {
									case <-chDone:
										{
										}
									case chExit <- true:
										{
										}
									}
									go play(elem.Value.(tSong), chPause, chContinue, chExit, chDone)
								}
							}
						case "prev":
							{
								elem = elem.Prev()

								if elem == nil {
									elem = playlist.Back()
									select {
									case <-chDone:
										{
										}
									case chExit <- true:
										{
										}
									}
									go play(elem.Value.(tSong), chPause, chContinue, chExit, chDone)
								}
								if elem != nil {
									select {
									case <-chDone:
										{
										}
									case chExit <- true:
										{
										}
									}
									go play(elem.Value.(tSong), chPause, chContinue, chExit, chDone)
								}
							}
						case "pause":
							{
								chPause <- true
							}
						case "play":
							{
								chContinue <- true
							}
						case "addSong":
							{
								go addSong()
							}
						case "exit":
							{
								os.Exit(3)
							}
						}
					}
				}
			case done := <-chDone:
				elem = elem.Next()

				if done {
					fmt.Println("New Track!")
				}

				if elem == nil {
					elem = playlist.Front()
					go play(elem.Value.(tSong), chPause, chContinue, chExit, chDone)
				}
				if elem.Next() != nil {
					elem = elem.Next()
					go play(elem.Value.(tSong), chPause, chContinue, chExit, chDone)
				}
			}
		}
	}
}
