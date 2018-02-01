//                           _                                       __  __          __
//     ____ ___  ____ _   __(_)__  ____ ___  ____ _____ _____  ___  / /_/ /_  ____  / /_
//    / __ `__ \/ __ \ | / / / _ \/ __ `__ \/ __ `/ __ `/ __ \/ _ \/ __/ __ \/ __ \/ __/
//   / / / / / / /_/ / |/ / /  __/ / / / / / /_/ / /_/ / / / /  __/ /_/ /_/ / /_/ / /_
//  /_/ /_/ /_/\____/|___/_/\___/_/ /_/ /_/\__,_/\__, /_/ /_/\___/\__/_.___/\____/\__/
//  https://github.com/magunetto/moviemagnetbot /____/
//

package main

import (
	"log"
)

func main() {

	InitModel()
	log.Printf("model inited")

	InitTMDb()
	log.Printf("tmdb inited")

	InitRARBG()
	log.Printf("rarbg inited")

	go RunBot()
	log.Printf("bot started")

	go RunHTTPServer()
	log.Printf("http server started")

	select {}
}
