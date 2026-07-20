// syslogCollector.js
const syslogLogs = [];

function addSyslogLog(entry) {
  syslogLogs.push(entry);
  if (syslogLogs.length > 1000) {
    syslogLogs.shift();
  }
}

module.exports = { syslogLogs, addSyslogLog };
