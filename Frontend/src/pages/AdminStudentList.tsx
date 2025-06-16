import { useState } from "react";
import { useStudents, useDeleteAdmin } from "../hooks/useAdmin";
import { getCurrentUser } from "../services/authService";
import Modal from "../components/Modal";
import StudentDocumentModal from "../components/StudentDocumentModal";
import {
  useStudentProfileByUserId,
  useStudentContactByUserId,
  useStudentAcademicByUserId,
  useStudentSettingsByUserId,
} from "../hooks/useProfile";
import PersonalInfoDisplay from "../components/PersonalInfoDisplay";
import ContactDetailsDisplay from "../components/ContactDetailsDisplay";
import AcademicBackgroundDisplay from "../components/AcademicBackgroundDisplay";
import SettingsDisplay from "../components/SettingsDisplay";
import { SectionCardBack } from "../components/SectionCardBack";
import { UserIcon } from "@heroicons/react/24/solid";

const AdminStudentList = () => {
  const { data: students, isLoading, isError } = useStudents();
  const deleteUser = useDeleteAdmin();
  const currentUser = getCurrentUser();
  const isCEO = currentUser?.roles?.includes("CEO");

  const [deleteTargetId, setDeleteTargetId] = useState<string | null>(null);
  const [viewDocsUserId, setViewDocsUserId] = useState<string | null>(null);
  const [viewDetailsUserId, setViewDetailsUserId] = useState<string | null>(null);

  const profileQuery = useStudentProfileByUserId(viewDetailsUserId || "");
  const contactQuery = useStudentContactByUserId(viewDetailsUserId || "");
  const academicQuery = useStudentAcademicByUserId(viewDetailsUserId || "");
  const settingsQuery = useStudentSettingsByUserId(viewDetailsUserId || "");

  const handleDelete = () => {
    if (!deleteTargetId) return;
    deleteUser.mutate(deleteTargetId, {
      onSuccess: () => setDeleteTargetId(null),
    });
  };

  return (
    <div className="p-6 max-w-7xl mx-auto space-y-6">
      <h1 className="text-3xl font-bold text-gray-800 flex items-center gap-2">
        <UserIcon className="w-6 h-6 text-blue-600" /> Student Accounts
      </h1>

      {isLoading ? (
        <p className="text-gray-500">Loading students...</p>
      ) : isError ? (
        <p className="text-red-500">Failed to fetch students.</p>
      ) : !students || students.length === 0 ? (
        <p className="text-gray-500">No students found.</p>
      ) : (
        <div className="overflow-x-auto rounded-xl shadow-md border border-gray-200 bg-white">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-600">Firstname</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-600">Lastname</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-600">Email</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-600">Country</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-600">Status</th>
                <th className="px-6 py-3 text-left text-sm font-semibold text-gray-600">Actions</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-100 text-sm">
              {students.map((student) => (
                <tr key={student.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4">{student.firstName}</td>
                  <td className="px-6 py-4">{student.lastName}</td>
                  <td className="px-6 py-4">{student.email}</td>
                  <td className="px-6 py-4">{student.country || "—"}</td>
                  <td className="px-6 py-4">
                    <span className={`inline-block px-2 py-1 rounded text-xs font-medium ${
                      student.status ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"
                    }`}>
                      {student.status ? "Active" : "Inactive"}
                    </span>
                  </td>
                  <td className="px-6 py-4 space-x-2">
                    <button
                      className="btn btn-sm bg-blue-100 text-blue-700 hover:bg-blue-200"
                      onClick={() => setViewDocsUserId(student.id)}
                    >
                      Documents
                    </button>
                    <button
                      className="btn btn-sm bg-indigo-100 text-indigo-700 hover:bg-indigo-200"
                      onClick={() => setViewDetailsUserId(student.id)}
                    >
                      Details
                    </button>
                    {isCEO && (
                      <button
                        className="btn btn-sm bg-red-100 text-red-700 hover:bg-red-200"
                        onClick={() => setDeleteTargetId(student.id)}
                      >
                        Delete
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Confirm Delete Modal */}
      {deleteTargetId && (
        <Modal isOpen={true} onClose={() => setDeleteTargetId(null)}>
          <div>
            <h2 className="text-lg font-bold text-red-600 mb-4">Confirm Delete</h2>
            <p>Are you sure you want to delete this student?</p>
            <div className="mt-4 flex justify-end gap-2">
              <button onClick={() => setDeleteTargetId(null)} className="btn">
                Cancel
              </button>
              <button
                onClick={handleDelete}
                className="btn bg-red-600 text-white"
                disabled={deleteUser.isPending}
              >
                {deleteUser.isPending ? "Deleting..." : "Delete"}
              </button>
            </div>
          </div>
        </Modal>
      )}

      {/* View Documents Modal */}
      {viewDocsUserId && (
        <StudentDocumentModal
          userId={viewDocsUserId}
          isOpen={!!viewDocsUserId}
          onClose={() => setViewDocsUserId(null)}
        />
      )}

      {/* View Details Modal */}
      {viewDetailsUserId && (
        <Modal isOpen={true} onClose={() => setViewDetailsUserId(null)}>
          <div className="max-h-[75vh] overflow-y-auto space-y-6">
            <h2 className="text-xl font-bold mb-4 text-gray-800">📋 Student Full Details</h2>
            <SectionCardBack title="Profile">
              <PersonalInfoDisplay data={profileQuery.data?.data} />
            </SectionCardBack>
            <SectionCardBack title="Contact">
              <ContactDetailsDisplay data={contactQuery.data?.data} />
            </SectionCardBack>
            <SectionCardBack title="Academic Background">
              <AcademicBackgroundDisplay data={academicQuery.data?.data} />
            </SectionCardBack>
            <SectionCardBack title="Settings">
              <SettingsDisplay data={settingsQuery.data?.data} />
            </SectionCardBack>
          </div>
        </Modal>
      )}
    </div>
  );
};

export default AdminStudentList;
