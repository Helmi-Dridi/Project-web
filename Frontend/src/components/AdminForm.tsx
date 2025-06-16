// components/AdminForm.tsx
import { useForm } from "react-hook-form";
import type { CreateAdminInput } from "../services/admin.service";

interface AdminFormProps {
  onSubmit: (data: Omit<CreateAdminInput, "companyID">) => void;
  onCancel: () => void;
  isSubmitting: boolean;
}

const AdminForm = ({ onSubmit, onCancel, isSubmitting }: AdminFormProps) => {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<Omit<CreateAdminInput, "companyID">>();

  const handleFormSubmit = (data: Omit<CreateAdminInput, "companyID">) => {
    onSubmit(data);
    reset();
  };

  return (
    <form
      onSubmit={handleSubmit(handleFormSubmit)}
      className="space-y-4 bg-gray-900 p-6 rounded-2xl shadow-xl text-white"
    >
      <h2 className="text-xl font-semibold">Create New Admin</h2>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div>
          <input
            {...register("firstName", { required: true })}
            placeholder="First Name"
            className="w-full px-4 py-2 rounded-md border border-gray-700 bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          {errors.firstName && <p className="text-red-400 text-sm mt-1">First name is required</p>}
        </div>

        <div>
          <input
            {...register("lastName", { required: true })}
            placeholder="Last Name"
            className="w-full px-4 py-2 rounded-md border border-gray-700 bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          {errors.lastName && <p className="text-red-400 text-sm mt-1">Last name is required</p>}
        </div>
      </div>

      <div>
        <input
          {...register("email", { required: true })}
          placeholder="Email"
          type="email"
          className="w-full px-4 py-2 rounded-md border border-gray-700 bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        {errors.email && <p className="text-red-400 text-sm mt-1">Email is required</p>}
      </div>

      <div>
        <input
          {...register("password", { required: true })}
          placeholder="Password"
          type="password"
          className="w-full px-4 py-2 rounded-md border border-gray-700 bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        {errors.password && <p className="text-red-400 text-sm mt-1">Password is required</p>}
      </div>

      <div>
        <input
          {...register("country")}
          placeholder="Country"
          className="w-full px-4 py-2 rounded-md border border-gray-700 bg-gray-800 text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <div className="flex justify-end gap-2 mt-6">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-md transition"
        >
          Cancel
        </button>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition"
          disabled={isSubmitting}
        >
          {isSubmitting ? "Creating..." : "Create Admin"}
        </button>
      </div>
    </form>
  );
};

export default AdminForm;