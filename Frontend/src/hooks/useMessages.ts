import { useEffect, useState, useCallback } from "react";
import {
  sendMessage,
  getMessagesWithUser,
  getPaginatedMessages,
  getUnreadMessageCount,
  getConversationPartners,
  searchMessages,
  markMessageAsRead,
  getAdminInbox,
} from "../services/messages.service";
import type { Message } from "../services/messages.service";

export const useMessagesWithUser = (userId: string, enabled = true) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchMessages = useCallback(async () => {
    if (!enabled || !userId) return;

    setLoading(true);
    try {
      const data = await getMessagesWithUser(userId);
      setMessages(data);
    } catch (error) {
      console.error("Failed to fetch messages:", error);
    } finally {
      setLoading(false);
    }
  }, [userId, enabled]);

  const appendMessage = useCallback((msg: Message) => {
  setMessages((prev) => {
    if (!msg.id) return prev; // Prevent bad data
    const alreadyExists = prev.some((m) => m.id === msg.id);
    if (alreadyExists) return prev;
    return [...prev, msg]; // 🔁 This triggers re-render
  });
}, []);


  useEffect(() => {
    fetchMessages();
  }, [fetchMessages]);

  return {
    messages,
    loading,
    refetch: fetchMessages,
    appendMessage,
  };
};
// ✅ Hook to send a message
export const useSendMessage = () => {
  const [sending, setSending] = useState(false);

  const send = async (
    receiverId: string,
    content: string,
    attachment?: string
  ) => {
    setSending(true);
    const msg = await sendMessage(receiverId, content, attachment);
    setSending(false);
    return msg;
  };

  return { send, sending };
};

// ✅ Hook to get paginated messages with a user
export const usePaginatedMessages = (
  userId: string,
  limit = 20,
  offset = 0
) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetch = async () => {
      setLoading(true);
      const data = await getPaginatedMessages(userId, limit, offset);
      setMessages(data);
      setLoading(false);
    };
    fetch();
  }, [userId, limit, offset]);

  return { messages, loading };
};

// ✅ Hook to track unread message count
export const useUnreadCount = () => {
  const [count, setCount] = useState(0);

  const fetchCount = useCallback(async () => {
    const c = await getUnreadMessageCount();
    setCount(c);
  }, []);

  useEffect(() => {
    fetchCount();
  }, [fetchCount]);

  return { count, refetch: fetchCount };
};

// ✅ Hook to get conversation partners
export const useConversationPartners = () => {
  const [partners, setPartners] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetch = async () => {
      const users = await getConversationPartners();
      setPartners(users);
      setLoading(false);
    };
    fetch();
  }, []);

  return { partners, loading };
};

// ✅ Hook for searching messages
export const useSearchMessages = (query: string) => {
  const [results, setResults] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!query) return;
    const fetch = async () => {
      const data = await searchMessages(query);
      setResults(data);
      setLoading(false);
    };
    fetch();
  }, [query]);

  return { results, loading };
};

// ✅ Utility: mark message as read
export const useMarkAsRead = () => {
  return markMessageAsRead;
};

// ✅ Hook to get admin inbox
export const useAdminInbox = () => {
  const [inbox, setInbox] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchInbox = async () => {
      const data = await getAdminInbox();
      setInbox(data);
      setLoading(false);
    };
    fetchInbox();
  }, []);

  return { inbox, loading };
};
