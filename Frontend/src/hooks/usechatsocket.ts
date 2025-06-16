import { useEffect, useRef } from "react";
import { getAuthHeaders, getCurrentUser } from "../services/authService";

type MessageHandler = (data: any) => void;

export const useChatSocket = (onMessage: MessageHandler) => {
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const headers = getAuthHeaders();
    const auth = headers.Authorization;
    const user = getCurrentUser();

    if (!auth || !auth.startsWith("Bearer ") || !user?.workCompanyId) {
      console.error("❌ Missing token or company ID for WebSocket connection.");
      return;
    }

    const token = auth.split(" ")[1];
    const url = `ws://localhost:8080/v1/messages/${user.workCompanyId}/ws?token=${encodeURIComponent(token)}`;
    const socket = new WebSocket(url);

    socketRef.current = socket;

    socket.onopen = () => {
      console.log("✅ WebSocket connected");
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.log("📩 WebSocket message:", data);
        onMessage(data);
      } catch (err) {
        console.error("❌ Failed to parse WebSocket message:", err);
      }
    };

    socket.onerror = (event) => {
      console.error("❌ WebSocket error:", event);
    };

    socket.onclose = (event) => {
      console.log("🛑 WebSocket closed", event);
    };

    // Clean up on unmount
    return () => {
      console.log("🔌 Cleaning up WebSocket...");
      socket.close();
    };
  }, [onMessage]);
};
