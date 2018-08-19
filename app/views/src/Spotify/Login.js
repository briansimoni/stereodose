export default function () {
	function getCookie(name) {
		var dc = document.cookie;
		var prefix = name + "=";
		var begin = dc.indexOf("; " + prefix);
		if (begin === -1) {
			begin = dc.indexOf(prefix);
			if (begin !== 0) return null;
		}
		else {
			begin += 2;
			var end = document.cookie.indexOf(";", begin);
			if (end === -1) {
				end = dc.length;
			}
		}
		// because unescape has been deprecated, replaced with decodeURI
		//return unescape(dc.substring(begin + prefix.length, end));
		return decodeURI(dc.substring(begin + prefix.length, end));
	}
	let cookie = getCookie("_stereodose-session");
	if (!cookie) {
		// throw new Error("No cookie boi");
		window.location = "/auth/login";
		return;
	}

	return new Promise((resolve, reject) => {
		let req = new XMLHttpRequest();
		req.open("GET", "/auth/refresh");
		req.addEventListener("readystatechange", function () {
			if (this.readyState === 4) {
				if (this.status === 200) {
					let data = JSON.parse(this.responseText);
					resolve(data.access_token);
				} else {
					reject(new Error("failed to get the token boi " + this.responseText));
				}
			}

		})
		req.send();
	});
}
