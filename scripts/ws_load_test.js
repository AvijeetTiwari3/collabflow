import ws from 'k6/ws';
import { check, sleep } from 'k6';
import { Trend, Counter } from 'k6/metrics';

const wsMessageLatency = new Trend('ws_message_latency_ms');
const messagesSent = new Counter('ws_messages_sent');
const messagesReceived = new Counter('ws_messages_received');

export const options = {
  stages: [
    { duration: '30s', target: 25 },
    { duration: '1m', target: 100 },  // 100 concurrent WebSocket sessions
    { duration: '30s', target: 0 },
  ],
  thresholds: {
    'ws_message_latency_ms': ['p(95)<30'], // Sub-30ms propagation target
  },
};

const WS_URL = __ENV.WS_URL || 'ws://localhost:8080/ws';

export default function () {
  const roomId = 'load-test-room-1';
  const url = `${WS_URL}?room=${roomId}`;

  const res = ws.connect(url, {}, function (socket) {
    socket.on('open', () => {
      // Periodically send CRDT and cursor movements
      socket.setInterval(() => {
        const start = Date.now();
        const payload = JSON.stringify({
          type: 'cursor_move',
          room_id: roomId,
          payload: { position: Math.floor(Math.random() * 1000) },
          timestamp: new Date().toISOString(),
        });

        socket.send(payload);
        messagesSent.add(1);
        wsMessageLatency.add(Date.now() - start);
      }, 1000); // Send every 1 second
    });

    socket.on('message', (data) => {
      messagesReceived.add(1);
    });

    socket.on('close', () => {});
    socket.on('error', (e) => {
      console.error('WS Error:', e.error());
    });

    socket.setTimeout(() => {
      socket.close();
    }, 25000);
  });

  check(res, { 'status is 101 Switching Protocols': (r) => r && r.status === 101 });
  sleep(1);
}
