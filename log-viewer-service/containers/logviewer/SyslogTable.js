// SyslogTable.js
import React, { useEffect, useState } from 'react';

export default function SyslogTable() {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    fetch('/logs/syslog')
      .then(res => res.json())
      .then(data => setLogs(data.events));
  }, []);

  return (
    <table>
      <thead>
        <tr>
          <th>Timestamp</th>
          <th>Severity</th>
          <th>Facility</th>
          <th>Message</th>
        </tr>
      </thead>
      <tbody>
        {logs.map((log, idx) => (
          <tr key={idx}>
            <td>{log.timestamp}</td>
            <td>{log.severity}</td>
            <td>{log.facility}</td>
            <td>{log.message}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
