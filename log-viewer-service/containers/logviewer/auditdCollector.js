// auditdCollector.js
const auditdLogs = [];

function addAuditdLog(entry) {
  // push raw text into the array
  auditdLogs.push(entry);
  // optional: limit size to avoid memory bloat
  if (auditdLogs.length > 1000) {
    auditdLogs.shift();
  }
}

module.exports = { auditdLogs, addAuditdLog };
