import { useEffect } from "react";
import { useForm, type FieldError, type FieldErrorsImpl, type Merge } from "react-hook-form";
import type { StudentSettings } from "../services/profile.service";
import { useStudentSettings } from "../hooks/useProfile";

type Props = {
  initialValues: StudentSettings;
  onSuccess: () => void;
};

const getErrorMessage = (
  error: FieldError | Merge<FieldError, FieldErrorsImpl<any>> | undefined
): string | null => (typeof error?.message === "string" ? error.message : null);

export default function SettingsForm({ initialValues, onSuccess }: Props) {
  const { upsert, isSaving } = useStudentSettings();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<StudentSettings>({
    defaultValues: initialValues,
  });

  useEffect(() => {
    reset(initialValues);
  }, [initialValues, reset]);

  const inputClass =
    "w-full px-3 py-2 rounded-md bg-gray-800 border border-gray-700 text-gray-100 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500";

  const onSubmit = (data: StudentSettings) => {
    upsert(data, { onSuccess });
  };

  const renderInput = (
    label: string,
    name: keyof StudentSettings,
    type: string = "text",
    validation: any = {},
    helper?: string
  ) => (
    <div>
      <label className="block text-sm font-medium text-gray-400 mb-1">{label}</label>
      <input
        type={type}
        {...register(name, validation)}
        className={inputClass}
      />
      {helper && <p className="text-xs text-gray-400 mt-1">{helper}</p>}
      {getErrorMessage(errors[name]) && (
        <p className="text-red-500 text-xs mt-1">{getErrorMessage(errors[name])}</p>
      )}
    </div>
  );

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Language preference */}
      <div>
        <label className="block text-sm font-medium text-gray-400 mb-1">Language Preference</label>
        <select {...register("languagePreference")} className={inputClass}>
          <option value="">Select a language</option>
          <option value="en">English</option>
          <option value="fr">French</option>
          <option value="de">German</option>
          <option value="es">Spanish</option>
        </select>
        <p className="text-xs text-gray-400 mt-1">This will determine the language of emails and dashboard.</p>
      </div>

      {/* Notification & Account Toggles */}
      <div className="space-y-2">
        {[
          { label: "Receive Email Notifications", field: "receiveEmailNotifications" },
          { label: "Receive SMS Notifications", field: "receiveSMSNotifications" },
          { label: "Mark Account as Deleted", field: "accountDeleted" },
        ].map(({ label, field }) => (
          <label key={field} className="flex items-center gap-2">
            <input
              type="checkbox"
              {...register(field as keyof StudentSettings)}
              className="accent-blue-600 scale-110"
            />
            <span className="text-sm text-gray-300">{label}</span>
          </label>
        ))}
      </div>

      {/* Emergency Contact Info */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
        {renderInput("Emergency Contact Name", "emergencyName", "text", {
          required: "Contact name is required",
        }, "Full name of your emergency contact.")}

        {renderInput("Relationship", "emergencyRelationship", "text", {
          required: "Relationship is required",
        }, "e.g., Parent, Sibling, Friend")}

        {renderInput("Emergency Phone", "emergencyPhoneNumber", "text", {
          required: "Phone number is required",
        }, "Include country code (e.g., +212612345678)")}

        {renderInput("Emergency Email", "emergencyEmail", "text", {
          pattern: {
            value: /^\S+@\S+\.\S+$/,
            message: "Invalid email format",
          },
        }, "Optional. Used for email contact in case of emergency.")}
      </div>

      {/* Submit */}
      <div className="flex justify-end">
        <button
          type="submit"
          disabled={isSaving}
          className="bg-blue-600 text-white px-5 py-2 rounded-md hover:bg-blue-700 transition disabled:opacity-50"
        >
          {isSaving ? "Saving..." : "Save"}
        </button>
      </div>
    </form>
  );
}
