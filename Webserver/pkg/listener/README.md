# Listener
Entrypoint for packet
- IP:Port 단위로 Listen하며, 패킷을 수신하면 filter chain을 타고 다음 단계로 이동합니다.

```
listeners:
- name: listener_0
  address:
  socket_address:
  address: 0.0.0.0
  port_value: 10000
```

listener는 서비스 할 포트를 기준으로 새로운 서비스를 제공한다고 생각하기
- Virtual Host