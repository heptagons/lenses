package dom


// s0 function first argument is the id to locate a DOM element.
// Second argument is the text to set into element's innerHTML.
const JSs0 = `
function s0(i, h) {
	const e = document.getElementById(i); if (e) {
		e.innerHTML = h;
	}
}
`

// s1 function first argument is the id to locate a DOM element.
// Second argument is the value to set into value attribute of element's.
const JSs1 = `
function s1(i, v) {
	const e = document.getElementById(i); if (e) {
		e.value = v;
	}
}
`

// s2 first argument is an array of classes of DOM elements wanted to be shown with 
// style block unless elements are TR type which then be used table-row style.
// Second argument is an array of classes of DOM elements wanted to be hidden style display none.
// Third argument is a map with keys as id to locate DOMs.
// Each element key 'h' is the innerHTML to set.
// Each element key 'v' is the value to set id data-value attr.
// Each element key 'c' is the class to be set in attribute class
// Each element key 'i' is the and element 1 is the data-value to set.
const JSs2 = `
function s2(ons, offs, set) {
	ons.forEach(on => {
		const elems = document.getElementsByClassName(on); for (var i=0; i < elems.length; i++) {
			if (elems[i].tagName == "TR") {
				elems[i].style.display = "table-row";
			} else {
				elems[i].style.display = "block";
			}
			
		}
	})
	offs.forEach(off => {
		const elems = document.getElementsByClassName(off); for (var i=0; i < elems.length; i++) {
			elems[i].style.display = "none";
		}
	})
	Object.entries(set).map(entry => {
		const i = entry[0]; const v = entry[1]; const e = document.getElementById(i); if (e) {
			if (v.h != undefined) {
				e.innerHTML = v.h;
			}
			if (v.v != undefined) {
				e.dataset.value = v.v;
			}
			if (v.c != undefined) {
				e.className = v.c;
			}
			if (v.i != undefined) {
				const input = document.getElementById(v.i); if (input) {
					e.innerHTML = input.value;
					e.dataset.value = input.value;
				}
			}
		}
	});
}
`

// function a0() is ajax
const JSa0 = `
function a0(path, id, args) {
	if (args) {
		Object.entries(args).map(t => {
			const k = t[0]; const v = t[1]; const e = document.getElementById(v); if (e) {
				path += "&" + k + "=" + e.value;
			}
		})
	}
	const xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			s0(id, this.responseText);
		}
	}
	s0(id, "...")
	xhttp.open("GET",%s + path, true);
	xhttp.send();
}
`

// function a1() is ajaxValue
const JSa1 = `
function a1(path, id, value, ons, offs, set) {
	const e = document.getElementById(value); if (e) {
		path += "&value=" + e.dataset.value;
	}
	a0(path, id);
	s2(ons, offs, set);
}`

// function a2() is ajaxValues
const JSa2 = `
function a2(path, id, base, n) {
	path += "&values="
	for (let i=0; i < parseInt(n); i++) {
		if (i > 0) {
			path += ",";
		}
		const e = document.getElementById(base + "-" + i); if (e) {
			path += e.dataset.value;
		}
	}
	a0(path, id);
}
`
