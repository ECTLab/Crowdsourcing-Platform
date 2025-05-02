document.addEventListener("DOMContentLoaded", main)
const API_BASE = 'http://localhost:3000';

function main() {
    initAuth()

    let [hash_lat, hash_lng, hash_zoom] = window.location.hash.substring(1).split(",").map((s) => parseFloat(s))
    if (!hash_lat) {
        window.location.hash = "#17/35.787878/51.375428"
    }

    maplibregl.setRTLTextPlugin(
        'https://unpkg.com/@mapbox/mapbox-gl-rtl-text@0.2.3/mapbox-gl-rtl-text.min.js',
        false
    );

    let map = new maplibregl.Map({
        container: 'map',
        style: 'https://tap30.services/v20240831/styles/Tapsi_v9.0/style.json',
        center: [hash_lng, hash_lat],
        zoom: hash_zoom,
        hash: true,
        attributionControl: false
    });

    let action = App(map)

    map.on('click', action)
}

function initAuth() {
    if (!localStorage.token) {
        fetch(`${API_BASE}/auth`, {
            method: "GET",
            headers: { "Content-Type": "application/json" },
        })
        .then(res => {
            if (!res.ok) throw new Error("Failed to fetch token");
            return res.json();
        })
        .then(data => {
            if (data.token) {
                localStorage.token = data.token;
                console.log("Token acquired and saved to localStorage");
            } else {
                throw new Error("Token not found in response");
            }
        })
        .catch(err => {
            alert("Authentication failed");
            console.error("Token generation failed:", err);
        });
    }
}


function App(map) {
    let reportBtn = document.getElementById('report')
    let routeBtn = document.getElementById('route')
    let confirmBtn = document.getElementById('confirm')
    let importBtn = document.getElementById('import')

    let reportApp = ReportCreator(map)
    let routeApp = RouteCreator(map)
    let confirmApp = ConfirmCreator(map)
    let importApp = ImportCreator(map)

    let app = reportApp
    reportBtn.disabled = true


    reportBtn.addEventListener('click', () => {
        app.clear()
        reportBtn.disabled = true
        routeBtn.disabled = false
        confirmBtn.disabled = false
        importBtn.disabled = false
        app = reportApp
    })

    routeBtn.addEventListener('click', () => {
        app.clear()
        routeBtn.disabled = true
        reportBtn.disabled = false
        confirmBtn.disabled = false
        importBtn.disabled = false
        app = routeApp
    })

    confirmBtn.addEventListener('click', () => {
        app.clear()
        confirmBtn.disabled = true
        routeBtn.disabled = false
        reportBtn.disabled = false
        importBtn.disabled = false
        app = confirmApp
        app.init()
    })

    importBtn.addEventListener('click', () => {
        app.clear()
        importBtn.disabled = true
        confirmBtn.disabled = false
        routeBtn.disabled = false
        reportBtn.disabled = false
        app = importApp
        app.init()
    })
    function action(e) {
        app.action(e)
    }

    return action
}

function ReportCreator(map) {
    let state = 0
    let markers = []

    function action(e) {
        let color = state ? "#1432ff" : "#000000"
        let marker = createMarker(e.lngLat, color, map)
        markers.push(marker)

        if (state === 1) {
            submitReport(markers)
            markers = []
        }
        state = (state + 1) % 2
    }

    function submitReport(markers) {
        let [m0, m1] = markers
        let m0ll = m0.getLngLat()
        let m1ll = m1.getLngLat()
        let report = {
            type: "POLICE",
            "engagement_location_time": {
                latitude: m0ll.lat,
                longitude: m0ll.lng,
                time: Date.now(),
            },
            "submit_location_time": {
                latitude: m1ll.lat,
                longitude: m1ll.lng,
                time: Date.now() + 5_000,
            },
            geometry: "testing_using_crowdsource_simulator",
            rideId: "testing_using_crowdsource_simulator",
            navigationRequestId: "testing_using_crowdsource_simulator",
        }
        fetch(`${API_BASE}/report`, {
            method: "POST",
            headers: { "Content-Type": "application/json", "api-key":localStorage.token },
            body: JSON.stringify({ report, token: localStorage.token })
        }).then(res => { if (!res.ok) { throw res }; return res.json() }).then(res => {
            console.log("POLICE reported")
            console.log(res)
        }).catch(err => {
            alert("failed reporting")
            console.error(err)
        }).finally(() => {
            m0.remove()
            m1.remove()
        })
    }

    function clear() {
        markers.forEach(m => m.remove())
        state = 0
    }

    return {action, clear}
}

function RouteCreator(map) {
    let from = null

    function action(e) {
        let color = from === null ? "#226611" : "#445511"
        let marker = createMarker(e.lngLat, color, map)

        if (from === null) {
            from = marker
        } else {
            let to = marker
            getRoute(from, to)
            from = null
        }
    }

    function getRoute(from, to) {
        let m0ll = from.getLngLat()
        let m1ll = to.getLngLat()

        let routeRequest = {
            from: {
                lat: m0ll.lat,
                lng: m0ll.lng,
            },
            to: {
                lat: m1ll.lat,
                lng: m1ll.lng,
            },
        }
        fetch(`${API_BASE}/route`, {
            method: "POST",
            headers: { "Content-Type": "application/json" , "api-key":localStorage.token},
            body: JSON.stringify({ from: routeRequest.from, to: routeRequest.to, token: localStorage.token })
        }).then(res => { if (!res.ok) { throw res }; return res.json() }).then(res => {
            console.log("route")
            console.log(res)
            let route = res.routes && res.routes[0]
            let geometry = polyline.decode(route.geometry)
            let leg = route && route.legs && res.routes[0].legs && res.routes[0].legs[0]
            if (!leg) {
                throw "no route"
            }
            let polices = leg.annotation.Police
            for (let i = 1; i < geometry.length - 1; i ++) {
                let u = geometry[i]
                let v = geometry[i + 1]
                let police = polices[i]
                console.log(i, police)
                if(!police) {
                    continue;
                }
                if (police.Exists) {
                    let reportId = police.Id
                    let offset = parseFloat(police.Offset)

                    dLatInMet = 111139*(v[0] - u[0])
                    dLngInMet = 85000*(v[1] - u[1])
                    let distance = Math.sqrt(dLatInMet*dLatInMet+dLngInMet*dLngInMet)
                    let frac = offset / distance

                    let LatLng = {
                        lat: (v[0] - u[0]) * frac + u[0],
                        lng: (v[1] - u[1]) * frac + u[1]
                    }
                    createPoliceReport(reportId, LatLng)
                }
            }
        }).catch(err => {
            alert("failed requesting route")
            console.error(err)
        }).finally(() => {
            from.remove()
            to.remove()
        })

        function createPoliceReport(reportId, LatLng) {
            let police = createPoliceMarker(LatLng, map)
            police.getElement().addEventListener('click', e => {
                if(!ConfirmCreatorActive) { return }

                let confirmed = confirm("report exists ?")

                fetch(`${API_BASE}/confirm`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json", "api-key":localStorage.token },
                    body: JSON.stringify({ reportId, confirmed, token: localStorage.token, "type":"POLICE" })
                }).then(res => { if (!res.ok) { throw res }; return res.json() }).then(res => {
                    console.log("report confirmation submitted")
                    console.log(res)
                }).catch(err => {
                    alert("failed submitting confirmation")
                    console.error(err)
                })
            })

            setTimeout(() => {
                police.remove()
            }, 60_000)
        }
    }

    function clear() {
        from && from.remove()
        from = null
    }

    return {action, clear}
}

var ConfirmCreatorActive = false
function ConfirmCreator(map) {
    let active = () => {
        ConfirmCreatorActive = true
        document.querySelectorAll('.police').forEach(el => el.classList.add("pointer"))
    }
    let deactive = () => {
        ConfirmCreatorActive = false
        document.querySelectorAll('.police').forEach(el => el.classList.remove("pointer"))
    }
    return {action: active, clear: deactive, init: () => {}}
}

function ImportCreator(map) {
    let modal = document.getElementById("import-modal")
    let input = document.getElementById("import-input")
    let submit = document.getElementById("do-import")

    let polices = []
    let markers = []
    let draw = () => {
        for (let police of polices) {
            let agg = !!police[2]
            let marker = null
            if(agg) {
                marker = createPoliceMarker({
                    lat: police[0],
                    lng: police[1],
                }, map)
                let policeElement = marker.getElement()
                let span = document.createElement('span')
                span.innerText = police[3]
                policeElement.appendChild(span)
            } else {
                marker = createMarker({
                    lat: police[0],
                    lng: police[1],
                }, "red", map)
            }
            markers.push(marker)
        }
    }

    let deactive = () => {
        modal.hidden = true
        input.value = ""
        document.getElementById("import").disabled = false
        for (let marker of markers) {
            marker.remove()
        }
    }

    let populate = (listOfPolices) => {
        let v = ""
        if (!listOfPolices) {
            v = input.value
            if (v === "") {
                return
            }
        }

        try {
            if (!listOfPolices) {
                polices = JSON.parse(v)
            } else {
                polices = listOfPolices
            }

            if (!Array.isArray(polices)) {
                throw "not array"
            }
            for (let police of polices) {
                const shapeError = police + " not in shape of [float, float, bool, float]"
                if (police.length !== 4) {
                    throw shapeError
                }

                if (typeof police[0] !== 'number' || typeof police[1] !== 'number' || typeof police[3] !== 'number') {
                    throw shapeError
                }

                if (typeof police[2] !== 'boolean') {
                    throw shapeError
                }
            }
        } catch (err) {
            console.error(err)
            alert("invalid input\n\ninput must be in this shape:\n[[float, float, bool, float], ...]")
            deactive()
            return;
        }
        deactive()
        draw()
    }

    let active = () => {
        fetch(`${API_BASE}/visibility`, {
            method: "POST",
            headers: { "Content-Type": "application/json", "api-key":localStorage.token },
            body: JSON.stringify({token: localStorage.token})
        }).then(res => { if (!res.ok) { throw res }; return res.json() }).then(res => {
            console.log("visibility called")
            console.log(res)
            populate(res)
        }).catch(err => {
            alert("failed reporting")
            console.error(err)
        })
    }
    let active_ = () => {
        modal.hidden = false
        input.value = ""
        submit.addEventListener('click', populate)
    }
    return {action: () => {}, clear: deactive, init: active}
}

function createMarker(LatLng, color, map) {
    let marker = new maplibregl.Marker({ color }).setLngLat(LatLng).addTo(map);
    return marker
}

function createPoliceMarker(LatLng, map) {
    let img = document.createElement('div')
    img.style.backgroundImage = "url(/police.svg)"
    img.className = ConfirmCreatorActive ? 'police pointer' : 'police'

    let marker = new maplibregl.Marker({ element: img }).setLngLat(LatLng).addTo(map);
    return marker
}
//
//async function initAuth() {
//    let token = localStorage.token
//
//    document.getElementById("logout").addEventListener('click', logout)
//
//    function logout() {
//        localStorage.token = ""
//        window.location.replace("/login");
//    }
//}