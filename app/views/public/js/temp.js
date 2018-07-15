
let MyDeviceID = null;

const getMySavedPlaylists = function() {
	console.log("getting my playlists");
	return new Promise( (resolve, reject) => {
		let req = new XMLHttpRequest();
        req.open("GET", "/api/playlists/me");
        req.addEventListener("readystatechange", function() {
            if (this.readyState === 4) {
				if (this.status === 204) {
					// TODO: parse json
					console.log(this.responseText);
					resolve(this.responseText);
				} else {
					console.log(this.statusText);
					reject(new Error(String(this.status) + "Unable to transfer player to this player: " + this.statusText));
				}
			}
        })
	});
}


const getSpotifyPlaylists = function(token) {
    console.log("trying to get spotify playlists");
    return new Promise( (resolve, reject) => {
        let req = new XMLHttpRequest();
        req.open("GET", "https://api.spotify.com/v1/me/playlists?limit=50&offset=0");
        req.setRequestHeader("Authorization", "Bearer " + token);
		req.setRequestHeader("Content-Type", "application/json");
		req.setRequestHeader("Accept", "application/json");
        req.addEventListener("readystatechange", function() {
            console.log("READY STATE CHANGE!");
            if (this.readyState === 4) {
				if (this.status === 200) {
					let playlists = JSON.parse(this.responseText);
					resolve(playlists.items);
				} else {
					console.log(this.statusText);
					reject(new Error(String(this.status) + "Unable to transfer player to this player: " + this.statusText));
				}
			}
        });
        req.send();
    });
}

const addToStereodose = function(playlistID) {
	return new Promise( (resolve, reject) => {
        let req = new XMLHttpRequest();
        req.open("POST", "/api/playlists/");
		req.setRequestHeader("Content-Type", "application/json");
        req.addEventListener("readystatechange", function() {
            if (this.readyState === 4) {
				if (this.status === 200) {
					resolve(true);
				} else {
					alert(this.statusText + " " + this.responseText);
					reject(new Error(String(this.status) + "Unable to add to Stereodose " + this.statusText));
				}
			}
		});
		let data = {
			SpotifyID: playlistID
		};
        req.send(JSON.stringify(data));
    });
}

const deleteFromStereodose = function(playlistID) {
	return new Promise( (resolve, reject) => {
        let req = new XMLHttpRequest();
        req.open("DELETE", "/api/playlists/" + playlistID);
        req.addEventListener("readystatechange", function() {
            if (this.readyState === 4) {
				if (this.status === 200) {
					resolve(true);
				} else {
					alert(this.statusText + " " + this.responseText);
					reject(new Error(String(this.status) + "Unable to delete: " + this.statusText));
				}
			}
		});
        req.send();
    });
}

const renderSpotifyPlaylists = function(playlists) {
	console.log(playlists);
	console.log("rendering Spotify Playlists");
	let ol = document.getElementById("my-playlists");
	playlists.forEach((playlist) => {
		let entry = document.createElement('li');
		let button = document.createElement('button');
		let deleteButton = document.createElement('button');
		entry.appendChild(document.createTextNode(playlist.name));
		entry.appendChild(button);
		entry.appendChild(deleteButton);
		button.appendChild(document.createTextNode("Add to Stereodose"));
		deleteButton.appendChild(document.createTextNode("Remove from Stereodose"));

		button.setAttribute("data-spotify-id", playlist.id);
		button.addEventListener("click", function() {
			addToStereodose(this.getAttribute("data-spotify-id"));
		});
		deleteButton.setAttribute("data-spotify-id", playlist.id);
		deleteButton.addEventListener("click", function() {
			deleteFromStereodose(this.getAttribute("data-spotify-id"));
		})
		ol.appendChild(entry);
	});

}

const transferPlayback = function(token, deviceID) {
    console.log("maybe transverring playback");
    return new Promise( (resolve, reject) => {
        let req = new XMLHttpRequest();
        req.open("PUT", "https://api.spotify.com/v1/me/player");
        req.setRequestHeader("Authorization", "Bearer " + token);
	    req.setRequestHeader("Content-Type", "application/json");
        req.addEventListener("readystatechange", function() {
            console.log("READY STATE CHANGE!");
            if (this.readyState === 4) {
				if (this.status === 204) {
					console.log(this.responseText);
					resolve(this.responseText);
				} else {
					console.log(this.statusText);
					reject(new Error(String(this.status) + "Unable to transfer player to this player: " + this.statusText));
				}
			}
        })
        let data = {
            device_ids: [deviceID],
            play: false
        }
        req.send(JSON.stringify(data));
    });
}

const getMyPlaylists = function () {
	return new Promise((resolve, reject) => {
		let myReq = new XMLHttpRequest();
		myReq.open("GET", "/api/playlists/?limit=1000&offset=0");
		myReq.addEventListener("readystatechange", function () {
			if (this.readyState === 4) {
				if (this.status === 200) {
					let playlists = JSON.parse(this.responseText);
					resolve(playlists);
				} else {
					reject(new Error(String(this.status) + "Unable to get playlists: " + this.statusText));
				}
			}
		});
		myReq.send();
	});
}

// getSongs takes playlistID and downloads all of its tracks
const getSongs = function(playlistID) {
	return new Promise( (resolve, reject) => {
		let req = new XMLHttpRequest();
		req.open("GET", "/api/playlists/" + playlistID);
		req.addEventListener("readystatechange", function() {
			if (this.readyState === 4) {
				if (this.status === 200) {
					let tracks = JSON.parse(this.responseText).Tracks;
					resolve(tracks);
				} else {
					reject(new Error(String(this.status) + "Unable to get songs: " + this.statusText));
				}
			}
			
		})
		req.send();
	});
};

const playSong = function() {
	let playlistID = this.getAttribute("data-spotify-playlist-id");
	let trackID = this.getAttribute("data-stereodose-id");
	console.log(MyDeviceID + " " + playlistID + " " + trackID);
	req = new XMLHttpRequest();
	req.open("PUT", "https://api.spotify.com/v1/me/player/play?" + MyDeviceID);
	req.setRequestHeader("Authorization", "Bearer " + Token);
	req.setRequestHeader("Content-Type", "application/json");
	let data = {
		"context_uri": playlistID,
		"offset": { 
			"uri": "spotify:track:" + trackID
		}
	}
	req.addEventListener("readystatechange", function() {
		if (this.readyState === 4) {
			console.log(this.responseText);
			console.log("The song should be playing...");
		}
	});
	data = JSON.stringify(data);
	console.log(data);
	req.send(data);
}

const clearTracks = function() {
	let list = document.getElementById("tracks-ol");
	list.innerHTML = "";
}

const updateTracksView = async function() {
	console.log("updating track view?");
	console.log(this);
	clearTracks();
	let playlistID = this.getAttribute("data-stereodose-id");
	let spotifyPlaylistID = this.getAttribute("data-spotify-id");
	try {
        let tracks = await getSongs(playlistID);
        console.log(tracks);
		let tracksOl = document.getElementById("tracks-ol");
		tracks.forEach( (track) => {
			let entry = document.createElement('li');
			entry.setAttribute("data-spotify-playlist-id", spotifyPlaylistID);
			entry.setAttribute("data-stereodose-id", track.SpotifyID);
			entry.appendChild(document.createTextNode(track.Name));
			entry.addEventListener("click", playSong)
			tracksOl.appendChild(entry);
		});
	} catch(e) {
		console.log(e);
		console.error(e);
	}
};


let Main = async function (Token, DeviceID) {
	let playlists = null;
    MyDeviceID = DeviceID
    transferPlayback(Token, DeviceID)
	console.log("what is this!?!?!?" + DeviceID);
	playlists = await getSpotifyPlaylists(Token);
	renderSpotifyPlaylists(playlists);
	try {
        playlists = await getMyPlaylists();
        console.log(playlists);
		let ol = document.getElementById("playlist-ol");
		playlists.forEach((playlist) => {
			let entry = document.createElement('li');
			entry.appendChild(document.createTextNode(playlist.name));
			entry.setAttribute("data-stereodose-id", playlist.SpotifyID);
			entry.setAttribute("data-spotify-id", playlist.uri);
			entry.addEventListener("click", updateTracksView);
            ol.appendChild(entry);
		});
	} catch (e) {
		console.log("error!");
		console.error(e);
	}
};