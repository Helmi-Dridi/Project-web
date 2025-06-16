// in ChatSocket.ts
import { useEffect, useRef } from "react";

export function useChatSocket(onMessage: (data: any) => void) {
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const token = localStorage.getItem("token");
    const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);
    wsRef.current = ws;

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      onMessage(data);
    };

    ws.onclose = () => console.log("🔌 WebSocket closed");
    ws.onerror = (err) => console.error("❌ WebSocket error:", err);

    return () => ws.close();
  }, [onMessage]);

  const send = (msg: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    }
  };

  return { send };
}
