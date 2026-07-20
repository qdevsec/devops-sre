// server.js
const express = require('express');
const path = require('path');
const client = require('prom-client');
const { auditdLogs, addAuditdLog } = require('./auditdCollector');
const { syslogLogs, addSyslogLog } = require('./syslogCollector');

const app = express();
const PORT = 8080;

// Use raw text parser instead of JSON
app.use(express.text({ type: '*/*' }));

// Serve static files from ConfigMap
app.use(express.static('/placeholders'));

// --- Log ingestion endpoints (sidecars POST here) ---
app.post('/logs/auditd', (req, res) => {
  addAuditdLog(req.body);
  console.log("Added auditd:", req.body); // debug
  res.sendStatus(200);
});
app.post('/logs/syslog', (req, res) => {
  addSyslogLog(req.body);
  console.log("Added syslog:", req.body); // debug
  res.sendStatus(200);
});

// --- Log retrieval endpoints (UI GETs here) ---
app.get('/logs/auditd', (req, res) => {
  res.json({ status: "ok", message: "Auditd logs", events: auditdLogs });
});
app.get('/logs/syslog', (req, res) => {
  res.json({ status: "ok", message: "Syslog logs", events: syslogLogs });
});
app.get('/logs/sysmon', (req, res) => {
  res.sendFile(path.join('/placeholders', 'sysmon.json'));
});
app.get('/logs/cloud', (req, res) => {
  res.sendFile(path.join('/placeholders', 'cloud.json'));
});

// --- Application metrics ---
const collectDefaultMetrics = client.collectDefaultMetrics;
collectDefaultMetrics();
app.get('/metrics/app', async (req, res) => {
  res.set('Content-Type', client.register.contentType);
  res.end(await client.register.metrics());
});

// --- Kubernetes metrics (proxy Prometheus) ---
app.get('/metrics/k8s', async (req, res) => {
  res.redirect('http://prometheus-svc.web-ha.svc.cluster.local:9090/metrics');
});

app.listen(PORT, () => {
  console.log(`Log viewer running on port ${PORT}`);
});
