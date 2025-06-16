import { useState } from "react";
import { useDocumentsByUserId, useVerifyDocument } from "../hooks/hookdocument";
import { FaFileAlt, FaCheckCircle, FaTimesCircle, FaCommentDots } from "react-icons/fa";

type Props = {
  userId: string;
  onClose: () => void;
};

const StudentDocumentViewer = ({ userId, onClose }: Props) => {
  const { data: documents, isLoading } = useDocumentsByUserId(userId);
  const verifyDocument = useVerifyDocument();
  const [remarksMap, setRemarksMap] = useState<Record<string, string>>({});
  const [verified, setVerified] = useState<Record<string, boolean>>({});

  const handleRemarksChange = (id: string, value: string) => {
    setRemarksMap((prev) => ({ ...prev, [id]: value }));
  };

  const handleVerification = (id: string, status: "approved" | "rejected") => {
    const remarks = remarksMap[id] || "";
    verifyDocument.mutate(
      { id, payload: { verification: status, remarks } },
      {
        onSuccess: () => setVerified((prev) => ({ ...prev, [id]: true })),
      }
    );
  };

  return (
    <div className="max-h-[80vh] overflow-y-auto bg-gray-900 p-6 rounded-xl shadow-inner">
      <div className="flex justify-between items-center mb-6 border-b border-gray-700 pb-4">
        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
          <FaFileAlt className="text-blue-500" />
          Student Documents
        </h2>
        <button
          onClick={onClose}
          className="text-gray-400 hover:text-white text-sm"
        >
          ✖ Close
        </button>
      </div>

      {isLoading ? (
        <p className="text-gray-400">Loading documents...</p>
      ) : !documents || documents.length === 0 ? (
        <p className="text-gray-400 italic">No documents uploaded by this student.</p>
      ) : (
        <ul className="space-y-6">
          {documents.map((doc) => (
            <li
              key={doc.id}
              className="bg-gray-800 border border-gray-700 rounded-xl p-5 space-y-3 shadow-sm"
            >
              <div className="flex justify-between items-center">
                <p className="text-white font-medium capitalize flex items-center gap-2">
                  <FaFileAlt className="text-blue-400" />
                  {doc.documentType.replace("_", " ")}
                </p>
                <span
                  className={`text-xs font-semibold px-2 py-1 rounded capitalize ${
                    doc.verification === "approved"
                      ? "bg-green-600 text-white"
                      : doc.verification === "rejected"
                      ? "bg-red-600 text-white"
                      : "bg-yellow-500 text-black"
                  }`}
                >
                  {doc.verification}
                </span>
              </div>

              <div className="text-sm text-gray-300">
                <strong>File:</strong>{" "}
                <a
                  href={`http://localhost:8080/${doc.filePath}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-400 underline hover:text-blue-300"
                >
                  Download
                </a>
              </div>

              {doc.remarks && (
                <div className="flex items-start gap-2 text-sm text-gray-400 italic">
                  <FaCommentDots className="mt-0.5" />
                  {doc.remarks}
                </div>
              )}

              {doc.verification === "pending" && !verified[doc.id] && (
                <div className="space-y-3 pt-3">
                  <textarea
                    value={remarksMap[doc.id] || ""}
                    onChange={(e) => handleRemarksChange(doc.id, e.target.value)}
                    placeholder="Optional remarks for student..."
                    className="w-full p-2 rounded bg-gray-700 border border-gray-600 text-sm text-white placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-blue-500"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleVerification(doc.id, "approved")}
                      className="flex items-center gap-1 px-3 py-1 text-sm rounded bg-green-600 text-white hover:bg-green-700 transition"
                    >
                      <FaCheckCircle size={14} /> Approve
                    </button>
                    <button
                      onClick={() => handleVerification(doc.id, "rejected")}
                      className="flex items-center gap-1 px-3 py-1 text-sm rounded bg-red-600 text-white hover:bg-red-700 transition"
                    >
                      <FaTimesCircle size={14} /> Reject
                    </button>
                  </div>
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default StudentDocumentViewer;
