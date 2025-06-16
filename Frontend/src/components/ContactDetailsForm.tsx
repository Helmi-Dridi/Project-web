import { useEffect, useMemo } from "react";
import { useForm, Controller, type FieldError, type FieldErrorsImpl, type Merge } from "react-hook-form";
import Select from "react-select";
import countryList from "react-select-country-list";
import type { StudentContactDetails } from "../services/profile.service";
import { useStudentContact } from "../hooks/useProfile";

type Props = {
  initialValues: StudentContactDetails;
  onSuccess: () => void;
};

const getErrorMessage = (
  error: FieldError | Merge<FieldError, FieldErrorsImpl<any>> | undefined
): string | null => (typeof error?.message === "string" ? error.message : null);

export default function ContactDetailsForm({ initialValues, onSuccess }: Props) {
  const { upsert, isSaving } = useStudentContact();
  const countries = useMemo(() => countryList().getData(), []);

  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { errors },
  } = useForm<any>({
    defaultValues: {
      ...initialValues,
      country: initialValues?.country
        ? countries.find((c) => c.value === initialValues.country)
        : null,
    },
  });

  useEffect(() => {
    reset({
      ...initialValues,
      country: initialValues?.country
        ? countries.find((c) => c.value === initialValues.country)
        : null,
    });
  }, [initialValues, reset, countries]);

  const onSubmit = (data: any) => {
    const payload: StudentContactDetails = {
      ...data,
      country: data.country?.value || "",
    };
    upsert(payload, { onSuccess });
  };

  const renderInput = (
    label: string,
    name: keyof StudentContactDetails,
    validation: any = {},
    helper?: string
  ) => (
    <div>
      <label className="block text-sm font-medium text-gray-400 mb-1">{label}</label>
      <input
        type="text"
        {...register(name, validation)}
        className="w-full px-3 py-2 rounded-md bg-gray-800 border border-gray-700 text-gray-100 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      {helper && (
        <p className="text-xs text-gray-400 mt-1">{helper}</p>
      )}
      {getErrorMessage(errors[name]) && (
        <p className="text-red-500 text-xs mt-1">{getErrorMessage(errors[name])}</p>
      )}
    </div>
  );

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {renderInput("Email", "email", {
          required: "Email is required",
          pattern: {
            value: /^\S+@\S+\.\S+$/,
            message: "Invalid email format",
          },
        }, "Use a valid email address like name@example.com.")}

        {renderInput("Phone Number", "phoneNumber", {
          required: "Phone number is required",
        }, "Include country code (e.g., +212612345678).")}

        {renderInput("Address Line", "addressLine", {
          required: "Address is required",
        }, "Street address or PO box (e.g., 123 Main St).")}

        {renderInput("City", "city", {
          required: "City is required",
        }, "Enter the city where you currently reside.")}

        {/* Country using react-select-country-list */}
        <div>
          <label className="block text-sm font-medium text-gray-400 mb-1">Country</label>
          <Controller
            control={control}
            name="country"
            rules={{ required: "Country is required" }}
            render={({ field }) => (
              <Select
                {...field}
                options={countries}
                placeholder="Select your country"
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
          <p className="text-xs text-gray-400 mt-1">Choose the country of your current residence.</p>
          {getErrorMessage(errors.country) && (
            <p className="text-red-500 text-xs mt-1">{getErrorMessage(errors.country)}</p>
          )}
        </div>

        {renderInput("Zip Code", "zipCode", {
          required: "Zip code is required",
        }, "Enter your postal or ZIP code (e.g., 10001).")}

        {/* Preferred Contact Method */}
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-400 mb-1">Preferred Contact</label>
          <select
            {...register("preferredContactMethod", {
              required: "Preferred contact method is required",
            })}
            className="w-full px-3 py-2 rounded-md bg-gray-800 border border-gray-700 text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">Select method</option>
            <option value="email">Email</option>
            <option value="whatsapp">WhatsApp</option>
            <option value="call">Call</option>
          </select>
          <p className="text-xs text-gray-400 mt-1">
            How would you prefer we contact you?
          </p>
          {getErrorMessage(errors.preferredContactMethod) && (
            <p className="text-red-500 text-xs mt-1">
              {getErrorMessage(errors.preferredContactMethod)}
            </p>
          )}
        </div>
      </div>

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
