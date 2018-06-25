// let myReq = new XMLHttpRequest();
// myReq.open("GET", "/api/playlists/");
// myReq.addEventListener("readystatechange", function () {
// 	if (this.readyState === 4) {
// 		let playlists = JSON.parse(this.responseText);
// 	}
// });
// myReq.send();

const getMyPlaylists = function () {
	return new Promise((resolve, reject) => {
		let myReq = new XMLHttpRequest();
		myReq.open("GET", "/api/playlists/");
		myReq.addEventListener("readystatechange", function () {
			if (this.readyState === 4) {
				let playlists = JSON.parse(this.responseText);
				resolve(playlists);
			}
			if (this.status !== 200) {
				reject(new Error(String(this.status) + "Unable to get playlists: " + this.statusText));
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
				let tracks = JSON.parse(this.responseText).tracks;
				resolve(tracks);
			}
			if (this.status !== 200) {
				reject(new Error(String(this.status) + "Unable to get songs: " + this.statusText));
			}
		})
		req.send();
	});
};

const updateTracksView = async function() {
	console.log("updating track view?");
	console.log(this);
	playlistID = this.getAttribute("data-stereodose-id");
	try {
		let tracks = await getSongs(playlistID);
		let tracksOl = document.getElementById("tracks-ol");
		tracks.forEach( (track) => {
			let entry = document.createElement('li');
			entry.appendChild(document.createTextNode(track.Name));
			tracksOl.appendChild(entry);
		});

	} catch(e) {
		console.error(e);
	}
};


let Main = async function () {
	try {
		let playlists = await getMyPlaylists();
		let ol = document.getElementById("playlist-ol");
		playlists.forEach((playlist) => {
			let entry = document.createElement('li');
			entry.appendChild(document.createTextNode(playlist.name));
			entry.setAttribute("data-stereodose-id", playlist.ID);
			entry.addEventListener("click", updateTracksView);
			ol.appendChild(entry);
		});
	} catch (e) {
		console.log("error!");
		console.error(e);
	}
};