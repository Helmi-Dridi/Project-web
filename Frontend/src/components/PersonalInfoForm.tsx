import { useEffect, useMemo } from "react";
import { useForm, Controller } from "react-hook-form";
import Select from "react-select";
import type { StudentProfile } from "../services/profile.service";
import { useStudentProfile } from "../hooks/useProfile";
import countryList from "react-select-country-list";

type Props = {
  initialValues: StudentProfile;
  onSuccess: () => void;
};

export default function PersonalInfoForm({ initialValues, onSuccess }: Props) {
  const { upsert, isSaving } = useStudentProfile();

  const countryOptions = useMemo(() => countryList().getData(), []);

  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors },
  } = useForm<any>({
    defaultValues: {
      ...initialValues,
      nationality: initialValues?.nationality
        ? countryOptions.find((c) => c.value === initialValues.nationality)
        : null,
    },
  });

  useEffect(() => {
    reset({
      ...initialValues,
      nationality: initialValues?.nationality
        ? countryOptions.find((c) => c.value === initialValues.nationality)
        : null,
    });
  }, [initialValues, reset, countryOptions]);

  const onSubmit = (data: any) => {
    const payload: StudentProfile = {
      ...data,
      nationality: data.nationality?.value || "",
      dateOfBirth: data.dateOfBirth
        ? new Date(data.dateOfBirth).toISOString()
        : undefined,
      step: true,
    };
    upsert(payload, { onSuccess });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Date of Birth */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">Date of Birth</label>
          <input
            type="date"
            {...register("dateOfBirth", { required: "Date of birth is required" })}
            className="w-full rounded-md bg-gray-800 border border-gray-700 text-gray-100 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p className="text-xs text-gray-400 mt-1">Format: YYYY-MM-DD</p>
          {errors.dateOfBirth?.message && (
            <p className="text-red-500 text-xs mt-1">{errors.dateOfBirth.message.toString()}</p>
          )}
        </div>

        {/* Gender */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">Gender</label>
          <select
            {...register("gender", { required: "Please select your gender" })}
            className="w-full rounded-md bg-gray-800 border border-gray-700 text-gray-100 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Select gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
          <p className="text-xs text-gray-400 mt-1">As listed on your official documents.</p>
          {errors.gender?.message && (
            <p className="text-red-500 text-xs mt-1">{errors.gender.message.toString()}</p>
          )}
        </div>

        {/* Nationality */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">Nationality</label>
          <Controller
            control={control}
            name="nationality"
            rules={{ required: "Nationality is required" }}
            render={({ field }) => (
              <Select
                {...field}
                options={countryOptions}
                placeholder="Select your nationality"
                classNamePrefix="react-select"
                isSearchable
                styles={{
                  control: (base) => ({
                    ...base,
                    backgroundColor: "#1f2937",
                    borderColor: "#374151",
                    color: "white",
                  }),
                  menu: (base) => ({
                    ...base,
                    backgroundColor: "#1f2937",
                  }),
                  singleValue: (base) => ({
                    ...base,
                    color: "white",
                  }),
                }}
              />
            )}
          />
          <p className="text-xs text-gray-400 mt-1">Choose from 190+ recognized countries.</p>
          {errors.nationality?.message && (
            <p className="text-red-500 text-xs mt-1">{errors.nationality.message.toString()}</p>
          )}
        </div>

        {/* Passport Number */}
        <div>
          <label className="block text-sm font-medium text-gray-300 mb-1">Passport Number</label>
          <input
            type="text"
            {...register("passportNumber", { required: "Passport number is required" })}
            className="w-full rounded-md bg-gray-800 border border-gray-700 text-gray-100 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p className="text-xs text-gray-400 mt-1">As printed on your passport.</p>
          {errors.passportNumber?.message && (
            <p className="text-red-500 text-xs mt-1">{errors.passportNumber.message.toString()}</p>
          )}
        </div>

        {/* National ID */}
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-300 mb-1">National ID</label>
          <input
            type="text"
            {...register("nationalId", { required: "National ID is required" })}
            className="w-full rounded-md bg-gray-800 border border-gray-700 text-gray-100 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p className="text-xs text-gray-400 mt-1">Government-issued ID if available.</p>
          {errors.nationalId?.message && (
            <p className="text-red-500 text-xs mt-1">{errors.nationalId.message.toString()}</p>
          )}
        </div>
      </div>

      <div className="flex justify-end">
        <button
          type="submit"
          disabled={isSaving}
          className="bg-blue-600 text-white px-6 py-2 rounded-md hover:bg-blue-700 transition disabled:opacity-50"
        >
          {isSaving ? "Saving..." : "Save"}
        </button>
      </div>
    </form>
  );
}
