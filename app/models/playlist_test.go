package models

import (
	"testing"
)

func TestStereodosePlaylistService_CreatePlaylistBySpotifyID(t *testing.T) {

	// connectionString := "postgresql://postgres:development@127.0.0.1:5432/stereodose?sslmode=disable"
	// db, err := gorm.Open("postgres", connectionString)
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }
	// db.Debug().DropTable(User{}, Playlist{}, UserImage{}, PlaylistImage{}, Track{})
	// db.AutoMigrate(User{}, Playlist{}, UserImage{}, PlaylistImage{}, Track{})
	// // var store = sessions.NewCookieStore([]byte("something-very-secret"))
	// // DB := NewStereodoseDB(db, store)

	// playlist := &Playlist{
	// 	Name:      "playlist1",
	// 	SpotifyID: "1",
	// 	Tracks: []Track{
	// 		Track{
	// 			Name:      "track1",
	// 			SpotifyID: "asdf",
	// 		},
	// 		Track{
	// 			Name:      "track2",
	// 			SpotifyID: "qwer",
	// 		},
	// 	},
	// }

	// playlist2 := &Playlist{
	// 	Name:      "playlist2",
	// 	SpotifyID: "2",
	// 	Tracks: []Track{
	// 		Track{
	// 			Name:      "track3",
	// 			SpotifyID: "asdf2",
	// 		},
	// 		Track{
	// 			Name:      "track4",
	// 			SpotifyID: "qwer2",
	// 		},
	// 	},
	// }

	// playlist3 := &Playlist{
	// 	Name:      "playlist3",
	// 	SpotifyID: "3",
	// 	Tracks: []Track{
	// 		Track{
	// 			Name:      "track1",
	// 			SpotifyID: "asdf",
	// 		},
	// 	},
	// }

	// err = db.Debug().Save(playlist).Error
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// err = db.Debug().Save(playlist2).Error
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// err = db.Debug().Save(playlist3).Error
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// playlists := []Playlist{}
	// err = db.Debug().Preload("Tracks").Offset("0").Limit("10").Find(&playlists).Error
	// if err != nil {
	// 	t.Error(err.Error())
	// }
	// if len(playlists) == 0 {
	// 	t.Error("playlists length 0")
	// }
	// log.Println(playlists[0].Name)
	// if playlists[0].Tracks[0].SpotifyID != "asdf" {
	// 	t.Error("Expected playlist id to be asdf, got:" + playlists[0].Tracks[0].SpotifyID)
	// }

	// err = db.Debug().(playlist).Error
	// if err != nil {
	// 	t.Error(err.Error())
	// }

	// err = db.Debug().Save(playlist).Error
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }
}
