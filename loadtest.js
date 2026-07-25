import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
  vus: 50,
  duration: '30s',
};

const TOKEN = __ENV.TOKEN || 'dev-token';
const BASE_URL = __ENV.BASE_URL || 'ws://localhost:8080/ws';

export default function () {
  const machineId = `k6-agent-${__VU}`;
  const url = `${BASE_URL}?machine_id=${machineId}`;

  const params = {
    headers: {
      Authorization: `Bearer ${TOKEN}`,
    },
  };

  const res = ws.connect(url, params, function (socket) {
    socket.on('open', () => {
      socket.setInterval(() => {
        const payload = {
          cpu: [{ name: "cpu0", usage_percent: Math.random() * 100 }],
          memory: { total: 1024, available: 512, used: 512 },
          disk: {
            disk_space: { path: "/", total: 100, available: 50, used: 50, usage_percent: 50 },
            disk_io_usage: []
          },
          network: [],
          processes: [],
          agent_ip: '127.0.0.1'
        };
        socket.send(JSON.stringify(payload));
      }, 5000);
    });

    socket.on('close', () => {
    });

    socket.on('error', (e) => {
      if (e.error() != 'websocket: close sent') {
        console.log('An unexpected error occured: ', e.error());
      }
    });

    socket.setTimeout(function () {
      socket.close();
    }, 30000);
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}
