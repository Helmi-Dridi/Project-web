import { useRef, useState, useEffect } from "react";
import {
  FaFileAlt,

  FaUpload,
  FaEye,
} from "react-icons/fa";
import classNames from "classnames";
import { useDocuments, useUploadDocument } from "../hooks/hookdocument";
import { Dialog} from "@headlessui/react";
import type { StudentDocumentTable } from "../services/documentService";

const documentLabels: Record<string, string> = {
  passport: "Passport",
  transcript: "Transcripts",
  cv: "CV",
  motivation_letter: "Motivation Letter",
  language_test: "Language Certificate",
};

const BASE_URL = "http://localhost:8080";

export default function DocumentsPage() {
  const { data } = useDocuments(1, 20);
  const { mutate: uploadDocument } = useUploadDocument();
  const fileInputRefs = useRef<Record<string, HTMLInputElement | null>>({});
  const [previewDoc, setPreviewDoc] = useState<string | null>(null);
  const [pendingUpload, setPendingUpload] = useState<{ docType: string; file: File } | null>(null);
  const [reuploadedDocs, setReuploadedDocs] = useState<Record<string, boolean>>({});

  const uploadedDocs: Record<string, StudentDocumentTable> = {};
  data?.items?.forEach((doc) => {
    uploadedDocs[doc.documentType] = doc;
  });

  const allTypes = Object.keys(documentLabels);
  const [selectedFilter, setSelectedFilter] = useState("All");

  const handleFileSelect = (docType: string) => {
    const input = fileInputRefs.current[docType];
    if (!input?.files?.length) return;
    const file = input.files[0];
    const fileURL = URL.createObjectURL(file);
    setPendingUpload({ docType, file });
    setPreviewDoc(fileURL);
  };

  const confirmUpload = () => {
    if (!pendingUpload) return;
    const formData = new FormData();
    formData.append("documentType", pendingUpload.docType);
    formData.append("file", pendingUpload.file);
    uploadDocument(formData, {
      onSuccess: () => {
        setReuploadedDocs((prev) => ({
          ...prev,
          [pendingUpload.docType]: true,
        }));
        setPendingUpload(null);
        setPreviewDoc(null);
      },
    });
  };

  useEffect(() => {
    return () => {
      if (previewDoc?.startsWith("blob:")) {
        URL.revokeObjectURL(previewDoc);
      }
    };
  }, [previewDoc]);

  const filteredTypes = allTypes.filter((type) => {
    if (selectedFilter === "All") return true;
    return uploadedDocs[type]?.verification === selectedFilter.toLowerCase();
  });

  return (
    <div className="max-w-6xl mx-auto px-4 py-10">
      <div className="bg-gray-900 border border-gray-700 rounded-3xl shadow-lg p-8 space-y-10">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold text-white flex items-center gap-2">
            <FaFileAlt className="text-blue-500" />
            My Documents
          </h1>
          <div className="flex gap-3">
            {["All", "Approved", "Pending", "Rejected"].map((status) => (
              <button
                key={status}
                onClick={() => setSelectedFilter(status)}
                className={classNames(
                  "px-3 py-1 text-sm rounded-full",
                  selectedFilter === status
                    ? "bg-blue-600 text-white"
                    : "bg-gray-700 text-gray-300 hover:bg-gray-600"
                )}
              >
                {status}
              </button>
            ))}
          </div>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-8">
          {filteredTypes.map((type) => {
            const doc = uploadedDocs[type];
            return (
              <div
                key={type}
                className={classNames(
                  "bg-gray-800 text-white rounded-2xl p-6 shadow-md border transition-all duration-200 min-h-[280px] hover:ring-2 hover:ring-blue-500 flex flex-col",
                  {
                    "border-green-400": doc?.verification === "approved",
                    "border-yellow-400": doc?.verification === "pending",
                    "border-red-400": doc?.verification === "rejected",
                    "border-gray-600": !doc,
                  }
                )}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2 text-lg font-semibold">
                    <FaFileAlt className="text-blue-500" />
                    {documentLabels[type]}
                  </div>
                  {doc?.verification === "approved" && (
                    <span className="text-xs bg-green-600 text-white px-2 py-1 rounded">Approved</span>
                  )}
                  {doc?.verification === "pending" && (
                    <span className="text-xs bg-yellow-500 text-black px-2 py-1 rounded">Pending</span>
                  )}
                  {doc?.verification === "rejected" && (
                    <span className="text-xs bg-red-600 text-white px-2 py-1 rounded">Rejected</span>
                  )}
                </div>

                <p className="text-sm text-gray-400 mt-2">
                  Upload your official {documentLabels[type]} as PDF or image.
                </p>

                {doc?.remarks && (
                  <p className="text-sm text-red-400 mt-1">🛈 {doc.remarks}</p>
                )}

                <div className="mt-auto flex gap-2 pt-4">
                  {doc?.filePath && doc.verification === "pending" && (
                    <button
                      className="flex-1 px-3 py-2 text-sm rounded-md bg-gray-700 hover:bg-gray-600 transition flex items-center justify-center gap-2"
                      onClick={() => setPreviewDoc(`${BASE_URL}/${doc.filePath}`)}
                    >
                      <FaEye className="w-4 h-4" />
                      View
                    </button>
                  )}

                  {doc?.verification !== "pending" && doc?.verification !== "approved" && !reuploadedDocs[type] && (
                    <label className="flex-1 cursor-pointer">
                      <input
                        type="file"
                        accept="application/pdf,image/*"
                        className="hidden"
                        ref={(ref) => {
                          fileInputRefs.current[type] = ref;
                        }}
                        onChange={() => handleFileSelect(type)}
                      />
                      <div className="px-3 py-2 text-sm rounded-md bg-blue-600 text-white hover:bg-blue-700 flex items-center justify-center gap-2">
                        <FaUpload className="w-4 h-4" />
                        {doc?.verification === "rejected" ? "Re-upload" : "Upload"}
                      </div>
                    </label>
                  )}
                </div>
              </div>
            );
          })}
        </div>

        <Dialog open={!!previewDoc} onClose={() => setPreviewDoc(null)} className="fixed z-50 inset-0">
          <div className="fixed inset-0 bg-black/60" aria-hidden="true" />
          <div className="fixed inset-0 flex items-center justify-center px-4">
            <Dialog.Panel className="bg-gray-800 text-white rounded-lg p-6 max-w-3xl w-full shadow-xl">
              <Dialog.Title className="text-lg font-semibold mb-4">
                {pendingUpload ? "Confirm Upload" : "Document Preview"}
              </Dialog.Title>

              {previewDoc && (
                previewDoc.endsWith(".pdf") ? (
                  <iframe src={previewDoc} className="w-full h-[500px] border rounded" />
                ) : (
                  <img src={previewDoc} className="w-full max-h-[500px] object-contain border rounded" />
                )
              )}

              <div className="mt-4 flex justify-end gap-3">
                <button
                  onClick={() => setPreviewDoc(null)}
                  className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-500"
                >
                  Cancel
                </button>
                {pendingUpload && (
                  <button
                    onClick={confirmUpload}
                    className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                  >
                    Confirm Upload
                  </button>
                )}
              </div>
            </Dialog.Panel>
          </div>
        </Dialog>
      </div>
    </div>
  );
}
