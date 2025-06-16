import { useParams, useNavigate } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import { useMessagesWithUser, useSendMessage } from "../hooks/useMessages";
import { useChatSocket } from "../hooks/usechatsocket";
import { useStudents } from "../hooks/useAdmin";
import { getCurrentUser } from "../services/authService";

export default function AdminChatPage() {
  const { userId } = useParams<{ userId: string }>();
  const navigate = useNavigate();
  const currentUser = getCurrentUser();

  const { data: students = [], isLoading } = useStudents();
  const [newMessage, setNewMessage] = useState("");
  const [typingUser, setTypingUser] = useState<string | null>(null);
  const chatEndRef = useRef<HTMLDivElement>(null);

  const { messages, loading, appendMessage } = useMessagesWithUser(userId ?? "", !!userId);
  const { send } = useSendMessage();

  // Auto-redirect to first student if none selected
  useEffect(() => {
    if (!userId && !isLoading && students.length > 0) {
      navigate(`/admin/chat/${students[0].id}`, { replace: true });
    }
  }, [userId, students, isLoading, navigate]);

  // Scroll to latest message
  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // ✅ WebSocket message and typing handler
  useChatSocket((msg) => {
    if (!userId || !currentUser) return;

    console.log("🔄 Incoming message:", msg, "Current user:", currentUser.ID, "Chatting with:", userId);

    if (msg.type === "message") {
      const isRelated =
  (String(msg.sender_id) === String(currentUser.ID) && String(msg.receiver_id) === String(userId)) ||
  (String(msg.sender_id) === String(userId) && String(msg.receiver_id) === String(currentUser.ID));

      if (isRelated) appendMessage(msg);
    }

    if (msg.type === "typing" && msg.from === userId) {
      setTypingUser(userId);
      setTimeout(() => setTypingUser(null), 2000);
    }
  });

  const handleSend = async () => {
    if (!newMessage.trim() || !userId) return;
    const msg = await send(userId, newMessage);
    appendMessage(msg);
    setNewMessage("");
  };

  const selectedStudent = students.find((u) => u.id === userId);
  const getFullName = (u: any) => `${u.firstName ?? ""} ${u.lastName ?? ""}`.trim();

  if (isLoading) {
    return <div className="p-4 text-gray-500">Loading students...</div>;
  }

  if (!userId || !selectedStudent) {
    return <div className="m-auto text-gray-500 text-lg">Select a student to start chatting.</div>;
  }

  return (
    <div className="flex h-screen">
      {/* Sidebar */}
      <aside className="w-72 bg-white border-r shadow-md">
        <div className="p-5 font-bold text-xl border-b">Students</div>
        <div className="overflow-y-auto h-full divide-y">
          {students.map((student) => (
            <button
              key={student.id}
              onClick={() => navigate(`/admin/chat/${student.id}`)}
              className={`w-full px-4 py-3 text-left flex items-center gap-3 hover:bg-blue-50 transition ${
                student.id === userId ? "bg-blue-100 font-semibold" : ""
              }`}
            >
              <div className="h-10 w-10 rounded-full bg-blue-500 text-white flex items-center justify-center uppercase font-bold">
                {(student.firstName?.[0] ?? "?") + (student.lastName?.[0] ?? "")}
              </div>
              <div>
                <div className="text-sm">{getFullName(student)}</div>
                <div className="text-xs text-gray-500">{student.email}</div>
              </div>
            </button>
          ))}
        </div>
      </aside>

      {/* Chat Window */}
      <main className="flex-1 flex flex-col bg-gray-50">
        <header className="p-4 border-b bg-white flex items-center gap-4 shadow-sm">
          <div className="h-10 w-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
            {(selectedStudent.firstName?.[0] ?? "?") + (selectedStudent.lastName?.[0] ?? "")}
          </div>
          <div>
            <div className="font-semibold text-base">{getFullName(selectedStudent)}</div>
            <div className="text-xs text-gray-500">{selectedStudent.email}</div>
          </div>
        </header>

        <section className="flex-1 overflow-y-auto px-6 py-4 space-y-3">
          {loading && <p className="text-sm text-gray-500">Loading chat...</p>}
          {messages.map((m) => {
            const isSentByMe = m.sender_id === currentUser?.ID;
            return (
              <div
                key={m.id}
                className={`max-w-lg px-4 py-2 rounded-xl shadow ${
                  isSentByMe
                    ? "bg-blue-600 text-white ml-auto text-right"
                    : "bg-white text-left"
                }`}
              >
                {m.content}
              </div>
            );
          })}

          {typingUser && (
            <p className="text-sm text-gray-500 italic animate-pulse">
              {getFullName(selectedStudent)} is typing...
            </p>
          )}
          <div ref={chatEndRef} />
        </section>

        <footer className="p-4 border-t bg-white flex items-center gap-3">
          <input
            className="flex-1 border border-gray-300 rounded-full px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            onKeyDown={(e) => e.key === "Enter" && handleSend()}
            placeholder="Type your message..."
          />
          <button
            onClick={handleSend}
            className="bg-blue-600 hover:bg-blue-700 transition text-white px-5 py-2 rounded-full"
          >
            Send
          </button>
        </footer>
      </main>
    </div>
  );
}
