import { useParams, useNavigate } from "react-router-dom";
import { useMessagesWithUser, useSendMessage } from "../hooks/useMessages";
import { useEffect, useState, useRef } from "react";
import { useChatSocket } from "../hooks/usechatsocket";
import { useAdminDetails } from "../hooks/useAdmin";
import { getCurrentUser } from "../services/authService";
import { format } from "date-fns";
import { Menu, Paperclip, Smile } from "lucide-react";

export default function ChatPage() {
  const { userId } = useParams<{ userId: string }>();
  const navigate = useNavigate();

  const { data: admins = [], isLoading: loadingAdmins } = useAdminDetails();
  const [newMessage, setNewMessage] = useState("");
  const [typingUser, setTypingUser] = useState<string | null>(null);
  const [search, setSearch] = useState("");
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const chatEndRef = useRef<HTMLDivElement>(null);

  const { messages, loading, appendMessage } = useMessagesWithUser(userId ?? "", !!userId);
  const { send } = useSendMessage();
  const currentUser = getCurrentUser();

  // ✅ WebSocket for receiving messages and typing
  useChatSocket((msg) => {
    if (!userId || !currentUser) return;

    console.log("🔄 Incoming message:", msg, "Current user ID:", currentUser.ID, "Chatting with:", userId);

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

  useEffect(() => {
    if (!userId && !loadingAdmins && admins.length > 0) {
      navigate(`/chat/${admins[0].id}`, { replace: true });
    }
  }, [userId, loadingAdmins, admins, navigate]);

  useEffect(() => {
    chatEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const handleSend = async () => {
    if (!newMessage.trim() || !userId) return;
    const msg = await send(userId, newMessage);
    appendMessage(msg);
    setNewMessage("");
  };

  const selectedAdmin = admins.find((admin) => admin.id === userId);
  const filteredAdmins = admins.filter((admin) =>
    `${admin.firstName} ${admin.lastName}`
      .toLowerCase()
      .includes(search.toLowerCase()) ||
    admin.email.toLowerCase().includes(search.toLowerCase())
  );

  if (loadingAdmins) {
    return <div className="p-4 text-gray-500">Loading admins...</div>;
  }

  if (!userId || !selectedAdmin) {
    return <div className="m-auto text-gray-500 text-lg">Select an admin to start chatting.</div>;
  }

  return (
    <div className="flex h-screen relative">
      {/* Sidebar */}
      <aside
        className={`bg-white border-r shadow-md md:w-72 w-full md:static absolute z-20 h-full transform transition-transform ${sidebarOpen ? "translate-x-0" : "-translate-x-full md:translate-x-0"}`}
      >
        <div className="p-5 font-bold text-xl border-b">Admins</div>
        <div className="p-4 border-b">
          <input
            type="text"
            placeholder="Search..."
            className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <div className="overflow-y-auto h-full divide-y">
          {filteredAdmins.map((admin) => (
            <button
              key={admin.id}
              onClick={() => navigate(`/chat/${admin.id}`)}
              className={`w-full px-4 py-3 text-left flex items-center gap-3 hover:bg-blue-50 transition ${
                admin.id === userId ? "bg-blue-100 font-semibold" : ""
              }`}
            >
              <div className="h-10 w-10 rounded-full bg-blue-500 text-white flex items-center justify-center uppercase font-bold">
                {admin.firstName[0]}
              </div>
              <div>
                <div className="text-sm">
                  {admin.firstName} {admin.lastName}
                </div>
                <div className="text-xs text-gray-500">{admin.email}</div>
              </div>
            </button>
          ))}
        </div>
      </aside>

      {/* Chat window */}
      <main className="flex-1 flex flex-col bg-gray-50">
        <header className="p-4 border-b bg-white flex items-center gap-4 shadow-sm">
          <button
            className="md:hidden p-2 rounded-md hover:bg-gray-100"
            onClick={() => setSidebarOpen((p) => !p)}
            aria-label="Toggle sidebar"
          >
            <Menu className="w-5 h-5" />
          </button>
          <div className="h-10 w-10 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
            {selectedAdmin.firstName[0]}
          </div>
          <div>
            <div className="font-semibold text-base">
              {selectedAdmin.firstName} {selectedAdmin.lastName}
            </div>
            <div className="text-xs text-gray-500">{selectedAdmin.email}</div>
          </div>
        </header>

        <section className="flex-1 overflow-y-auto px-6 py-4 space-y-3">
          {loading && <p className="text-sm text-gray-500">Loading chat...</p>}
          {messages.map((m) => {
            const isSentByMe = m.sender_id === currentUser?.ID;
            return (
              <div
                key={m.id}
                className={`max-w-lg px-4 py-2 rounded-xl shadow relative ${
                  isSentByMe
                    ? "bg-gradient-to-br from-blue-500 to-blue-700 text-white ml-auto"
                    : "bg-gray-100 text-gray-800"
                }`}
              >
                <p>{m.content}</p>
                <div
                  className={`text-xs mt-1 ${
                    isSentByMe ? "text-right text-blue-100" : "text-gray-500"
                  }`}
                >
                  {format(new Date(m.created_at), "p")}
                </div>
              </div>
            );
          })}
          {typingUser && (
            <p className="text-sm text-gray-500 italic animate-pulse">
              {selectedAdmin.firstName} is typing...
            </p>
          )}
          <div ref={chatEndRef} />
        </section>

        <footer className="p-4 border-t bg-white flex items-center gap-3">
          <button
            type="button"
            aria-label="Add emoji"
            className="text-gray-500 hover:text-blue-600"
          >
            <Smile className="w-5 h-5" />
          </button>
          <button
            type="button"
            aria-label="Attach file"
            className="text-gray-500 hover:text-blue-600"
          >
            <Paperclip className="w-5 h-5" />
          </button>
          <input
            className="flex-1 border border-gray-300 rounded-full px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500 bg-gray-50"
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
