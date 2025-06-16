import { useState } from "react";
import {
  useAdminDetails,
  useCreateAdmin,
  useDeleteAdmin,
  useAssignAdminRole,
  useUsersWithoutRole,
} from "../hooks/useAdmin";
import { useRoles } from "../hooks/useRole";
import Modal from "../components/Modal";
import AdminForm from "../components/AdminForm";
import AssignRoleModal from "../components/AssignRoleModal";
import { getCurrentUser } from "../services/authService";
import toast from "react-hot-toast";
import { UserIcon } from "@heroicons/react/24/solid";

const AdminList = () => {
  const currentUser = getCurrentUser();
  const isCEO = currentUser?.roles?.includes("CEO");

  const { data: adminData, isLoading: loadingAdmins } = useAdminDetails();
  const { data: noRoleData, isLoading: loadingNoRole } = useUsersWithoutRole();
  const { data: rolesData } = useRoles();
  const roles = Array.isArray(rolesData) ? rolesData : [];

  const createAdmin = useCreateAdmin();
  const deleteAdmin = useDeleteAdmin();
  const assignRole = useAssignAdminRole();

  const [isCreateModalOpen, setCreateModalOpen] = useState(false);
  const [selectedUserId, setSelectedUserId] = useState<string | null>(null);
  const [selectedRoleId, setSelectedRoleId] = useState<string>("");
  const [deleteTargetId, setDeleteTargetId] = useState<string | null>(null);

  const handleAssignRole = () => {
    if (!selectedUserId || !selectedRoleId) return;
    assignRole.mutate(
      { userId: selectedUserId, roleId: selectedRoleId },
      {
        onSuccess: () => {
          toast.success("Role assigned successfully");
          setSelectedUserId(null);
          setSelectedRoleId("");
        },
      }
    );
  };

  const handleDelete = () => {
    if (!deleteTargetId) return;
    deleteAdmin.mutate(deleteTargetId, {
      onSuccess: () => {
        toast.success("Admin deleted successfully");
        setDeleteTargetId(null);
      },
    });
  };

  const renderTable = (users: any[], title: string, allowAssign = false) => (
    <div className="mt-10 rounded-2xl bg-white shadow-md p-6">
      <h2 className="text-xl font-semibold flex items-center gap-2 text-gray-800 mb-4">
        <UserIcon className="text-blue-500 w-5 h-5" />
        {title}
      </h2>
      {users.length === 0 ? (
        <p className="text-sm text-gray-500 italic px-4 py-6">No users found in this section.</p>
      ) : (
        <div className="overflow-x-auto rounded border border-gray-200">
          <table className="min-w-full">
            <thead>
              <tr className="bg-gray-100 text-left text-sm font-medium text-gray-700">
                <th className="py-3 px-4">Firstname</th>
                <th className="py-3 px-4">Lastname</th>
                <th className="py-3 px-4">Email</th>
                <th className="py-3 px-4">Country</th>
                <th className="py-3 px-4">Status</th>
                <th className="py-3 px-4">Actions</th>
              </tr>
            </thead>
            <tbody className="text-sm text-gray-700 divide-y">
              {users.map((user, idx) => {
                const invalidId =
                  !user.id || user.id === "00000000-0000-0000-0000-000000000000";
                return (
                  <tr key={idx} className="hover:bg-gray-50">
                    <td className="py-2 px-4">{user.firstName}</td>
                    <td className="py-2 px-4">{user.lastName}</td>
                    <td className="py-2 px-4">{user.email}</td>
                    <td className="py-2 px-4">{user.country || "—"}</td>
                    <td className="py-2 px-4">
                      {user.status ? (
                        <span className="text-green-600 font-medium">Active</span>
                      ) : (
                        <span className="text-red-500 font-medium">Inactive</span>
                      )}
                    </td>
                    <td className="py-2 px-4 space-x-2">
                      {!invalidId && allowAssign && (
                        <button
                          className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 transition"
                          onClick={() => setSelectedUserId(user.id)}
                        >
                          Assign Role
                        </button>
                      )}
                      {!invalidId && isCEO && (
                        <button
                          className="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700 transition"
                          onClick={() => setDeleteTargetId(user.id)}
                        >
                          Delete
                        </button>
                      )}
                      {invalidId && (
                        <span className="text-xs text-gray-400">Invalid ID</span>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-800">👥 Admin Management</h1>
        {isCEO && (
          <button
            onClick={() => setCreateModalOpen(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 shadow"
          >
            + Add Admin
          </button>
        )}
      </div>

      {loadingAdmins || loadingNoRole ? (
        <p className="text-gray-500">Loading users...</p>
      ) : (
        <>
          {renderTable(adminData || [], "Admins")}
          {renderTable(noRoleData || [], "Users Without Roles", true)}
        </>
      )}

      {/* Create Admin Modal */}
      <Modal isOpen={isCreateModalOpen} onClose={() => setCreateModalOpen(false)}>
        <div className="bg-gray-900 text-white p-6 rounded-xl w-full max-w-lg">
          <h2 className="text-xl font-semibold mb-4">Create New Admin</h2>
          <AdminForm
            onSubmit={(data) =>
              createAdmin.mutate(
                { ...data, companyID: "" },
                {
                  onSuccess: () => {
                    toast.success("Admin created successfully");
                    setCreateModalOpen(false);
                  },
                }
              )
            }
            onCancel={() => setCreateModalOpen(false)}
            isSubmitting={createAdmin.isPending}
          />
        </div>
      </Modal>

      {/* Assign Role Modal */}
      {selectedUserId && (
        <Modal isOpen={true} onClose={() => setSelectedUserId(null)}>
          <AssignRoleModal
            roles={roles}
            selectedRoleId={selectedRoleId}
            onChange={setSelectedRoleId}
            onCancel={() => setSelectedUserId(null)}
            onConfirm={handleAssignRole}
          />
        </Modal>
      )}

      {/* Confirm Delete Modal */}
      {deleteTargetId && (
        <Modal isOpen={true} onClose={() => setDeleteTargetId(null)}>
          <div className="bg-white rounded-lg p-6">
            <h2 className="text-lg font-bold text-red-600 mb-4">Confirm Delete</h2>
            <p>Are you sure you want to delete this admin?</p>
            <div className="mt-4 flex justify-end gap-2">
              <button onClick={() => setDeleteTargetId(null)} className="btn">
                Cancel
              </button>
              <button
                onClick={handleDelete}
                className="btn bg-red-600 text-white"
                disabled={deleteAdmin.isPending}
              >
                {deleteAdmin.isPending ? "Deleting..." : "Delete"}
              </button>
            </div>
          </div>
        </Modal>
      )}
    </div>
  );
};

export default AdminList;
