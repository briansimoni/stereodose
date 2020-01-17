package models

import (
	"reflect"
	"testing"

	"github.com/zmb3/spotify"
)

func Test_simpleArtistsToString(t *testing.T) {
	parkwayDrive := spotify.SimpleArtist{
		Name: "Parkway Drive",
	}
	eminem := spotify.SimpleArtist{
		Name: "Eminem",
	}
	ironMaiden := spotify.SimpleArtist{
		Name: "Iron Maiden",
	}

	tests := []struct {
		name    string
		artists []spotify.SimpleArtist
		want    string
	}{
		{name: "Three artists", artists: []spotify.SimpleArtist{parkwayDrive, eminem, ironMaiden}, want: "Parkway Drive, Eminem, Iron Maiden"},
		{name: "One artist", artists: []spotify.SimpleArtist{ironMaiden}, want: "Iron Maiden"},
		{name: "No artist", artists: []spotify.SimpleArtist{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := simpleArtistsToString(tt.artists); got != tt.want {
				t.Errorf("simpleArtistToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_randomPlaylistFromSet(t *testing.T) {

	playlist1 := Playlist{
		Tracks: []Track{
			Track{ SpotifyID: "asdf" },
		},
	}

	wantedPlaylist := &Playlist{
		Tracks: []Track{
			Track{ SpotifyID: "asdf" },
			Track{ SpotifyID: "asdf" },
			Track{ SpotifyID: "asdf" },
			Track{ SpotifyID: "asdf" },
			Track{ SpotifyID: "asdf" },
		},
	}
	type args struct {
		playlists []Playlist
		length    int
	}
	tests := []struct {
		name string
		args args
		want *Playlist
	}{
		{ name: "One playlist. One Track", args: args{[]Playlist{playlist1}, 5}, want: wantedPlaylist },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := randomPlaylistFromSet(tt.args.playlists, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("randomPlaylistFromSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
