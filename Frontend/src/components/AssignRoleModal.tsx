import type { Role } from "../services/role.service";

interface AssignRoleModalProps {
  roles: Role[] | undefined | null;
  selectedRoleId: string;
  onChange: (roleId: string) => void;
  onConfirm: () => void;
  onCancel: () => void;
}

const AssignRoleModal = ({
  roles,
  selectedRoleId,
  onChange,
  onConfirm,
  onCancel,
}: AssignRoleModalProps) => {
  if (!Array.isArray(roles)) {
    console.error("AssignRoleModal: 'roles' is not an array", roles);
    return <p className="text-red-500">Failed to load roles.</p>;
  }

  return (
    <div className="text-gray-900 dark:text-gray-100">
      <h2 className="text-2xl font-bold mb-4">Assign Role</h2>

      <select
  value={selectedRoleId}
  onChange={(e) => onChange(e.target.value)}
  className="w-full p-3 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500 transition"
>
  <option value="">Select a role</option>
  {roles
    .filter((role) => role.name.toLowerCase() !== "ceo") // Filter out "CEO"
    .map((role) => (
      <option key={role.id} value={role.id}>
        {role.name}
      </option>
    ))}
</select>


      <div className="mt-6 flex justify-end gap-3">
        <button
          onClick={onCancel}
          className="px-4 py-2 rounded-md border border-gray-300 dark:border-gray-500 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-100 hover:bg-gray-200 dark:hover:bg-gray-600 transition"
        >
          Cancel
        </button>
        <button
          onClick={onConfirm}
          disabled={!selectedRoleId}
          className="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 disabled:opacity-50 transition"
        >
          Assign
        </button>
      </div>
    </div>
  );
};

export default AssignRoleModal;
