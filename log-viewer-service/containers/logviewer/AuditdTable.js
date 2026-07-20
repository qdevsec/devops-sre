// AuditdTable.js
import React, { useEffect, useState } from 'react';

export default function AuditdTable() {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    fetch('/logs/auditd')
      .then(res => res.json())
      .then(data => setLogs(data.events));
  }, []);

  return (
    <table>
      <thead>
        <tr>
          <th>Timestamp</th>
          <th>Syscall</th>
          <th>User</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        {logs.map((log, idx) => (
          <tr key={idx}>
            <td>{log.timestamp}</td>
            <td>{log.syscall}</td>
            <td>{log.user}</td>
            <td>{log.action}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
