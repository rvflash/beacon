// Copyright (c) 2018 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Analytics contains data to log sent by beacon.
type Analytics struct {
	Language  string `json:"lang,omitempty"`
	UserAgent string `json:"ua,omitempty"`
}

// Proof of concept
// Simple http server to test the navigator.sendBeacon() method to asynchronously send data on page unload.
// > https://developer.mozilla.org/en-US/docs/Web/API/Navigator/sendBeacon
func main() {
	// Handles data to logs on server side.
	// Also try to do some stuffs in JavaScript, but...
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		d := json.NewDecoder(r.Body)
		var l Analytics
		if err := d.Decode(&l); err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Logs events (here, just adds a cookie).
		cookie(w, "srv", "rv", false)
		fmt.Printf("hi %q guy with %q\n", l.Language, l.UserAgent)
		// Tries to continue the tracking with JavaScript events.
		http.ServeFile(w, r, "tracker.html")
	})

	// Handles events to fake a JavaScript call to an other audience manager.
	http.HandleFunc("/audience", func(w http.ResponseWriter, r *http.Request) {
		var os, by string
		c, err := r.Cookie("cli")
		if err == nil {
			os = c.Value
		}
		if c, err = r.Cookie("srv"); err == nil {
			by = c.Value
		}
		w.WriteHeader(http.StatusOK)
		// Logs an other events.
		fmt.Printf("a new hit of %q, powered by %q\n", os, by)
	})

	// Handles "just an other page on the same domain"
	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Cool, thanks!")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Resets the cookie and redirects to the landing page.
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		// Deletes both cookies, those created by client and server.
		cookie(w, "cli", "", true)
		cookie(w, "srv", "", true)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// Handles the rest of the traffic.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// No route: tooth!
			http.NotFound(w, r)
			return
		}
		// Landing page
		http.ServeFile(w, r, "index.html")
	})

	// Launches the server on the given port (see the environment var named BEACON_PORT).
	addr := port(os.Getenv("BEACON_PORT"))
	fmt.Printf("listening on localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func port(s string) string {
	if _, err := strconv.Atoi(s); err != nil {
		// Empty string or invalid port number.
		// Default port used as fallback.
		return ":8080"
	}
	return ":" + s
}

func cookie(w http.ResponseWriter, k, v string, delete bool) {
	// Duration
	days := 7
	if delete {
		days = -7
	}
	cookie := http.Cookie{
		Name:    k,
		Value:   v,
		Expires: time.Now().Add(time.Duration(days) * 24 * time.Hour),
	}
	http.SetCookie(w, &cookie)
}
