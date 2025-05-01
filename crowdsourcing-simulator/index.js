const express = require("express")
const cors = require("cors")
const path = require("path")
const bodyParser = require("body-parser")
const fetch   = require("node-fetch");

const app = express()

app.use(cors({
  origin: "*",
  methods: ["GET","POST","OPTIONS"],
  allowedHeaders: ["Content-Type","api-key"]
}))
app.use(bodyParser.json())

app.use("/", express.static(path.join(__dirname, "app")))
app.use("/login", express.static(path.join(__dirname, "app", "login.html")))

app.post("/report", (req, res) => {
  const { token, report } = req.body
  fetch("http://localhost:8001/crowdsourcing/in-ride-report", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "api-key": token
    },
    body: JSON.stringify(report)
  })
    .then(r => r.ok ? r.json() : Promise.reject(r))
    .then(r => res.send(r))
    .catch(err => res.status(err.status || 500).send(err))
})

app.post("/confirm", (req, res) => {
  const { token, reportId, confirmed } = req.body
  fetch(`http://localhost:8001/crowdsourcing/in-ride-report/${reportId}/confirm`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "api-key": token
    },
    body: JSON.stringify({ confirmed, type: "POLICE" })
  })
    .then(r => r.ok ? r.json() : Promise.reject(r))
    .then(r => res.send(r))
    .catch(err => res.status(err.status || 500).send(err))
})

app.post("/route", (req, res) => {
  const { token, from, to } = req.body
  fetch("http://localhost:8080/navigation/get-route", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "api-key": token
    },
    body: JSON.stringify({
      vehicleType: "car",
      origin: { latitude: from.lat, longitude: from.lng },
      destination: { latitude: to.lat, longitude: to.lng },
    })
  })
    .then(r => r.ok ? r.json() : Promise.reject(r))
    .then(r => res.send(r))
    .catch(err => res.status(err.status || 500).send(err))
})

app.post("/visibility", (req, res) => {
  const { token } = req.body
  fetch("http://localhost:8001/crowdsourcing/reports/visibility", {
    method: "GET",
    headers: { "api-key": token }
  })
    .then(r => r.ok ? r.json() : Promise.reject(r))
    .then(r => res.send(r))
    .catch(err => res.status(err.status || 500).send(err))
})

const port = parseInt(process.argv[2], 10) || 3000
app.listen(port, "0.0.0.0", () =>
  console.log(`running crowdsourcing simulator server on 0.0.0.0:${port}`)
)
