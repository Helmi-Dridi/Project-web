import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import type { StudentAcademicBackground } from "../services/profile.service";
import { useStudentAcademic } from "../hooks/useProfile";
import type { FieldError, FieldErrorsImpl, Merge } from "react-hook-form";

type Props = {
  initialValues: StudentAcademicBackground;
  onSuccess: () => void;
};

type FormFields = Omit<StudentAcademicBackground, "certificateFilePath"> & {
  certificate?: File;
};

const getErrorMessage = (
  error: FieldError | Merge<FieldError, FieldErrorsImpl<any>> | undefined
): string | null => (typeof error?.message === "string" ? error.message : null);

export default function AcademicBackgroundForm({
  initialValues,
  onSuccess,
}: Props) {
  const { upsert, isSaving } = useStudentAcademic();
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    formState: { errors },
  } = useForm<FormFields>({
    defaultValues: {
      ...initialValues,
      graduationYear: Number(initialValues.graduationYear) || undefined,
    },
  });

  useEffect(() => {
    reset({
      ...initialValues,
      graduationYear: Number(initialValues.graduationYear) || undefined,
    });
  }, [initialValues, reset]);

  const onSubmit = (data: FormFields) => {
    const formData = new FormData();
    formData.append("qualification", data.qualification);
    formData.append("institutionName", data.institutionName);
    formData.append("graduationYear", String(data.graduationYear));
    formData.append("gpaScore", data.gpaScore || "");
    formData.append("languageTestType", data.languageTestType || "");
    formData.append("languageTestScore", data.languageTestScore || "");
    if (selectedFile) {
      formData.append("certificate", selectedFile);
    }
    upsert(formData, { onSuccess });
  };

  const renderInput = (
    label: string,
    name: keyof FormFields,
    type: string = "text",
    validation: any = {},
    helper?: string
  ) => (
    <div>
      <label className="block text-sm font-medium text-gray-400 mb-1">{label}</label>
      <input
        type={type}
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
        {renderInput(
          "Qualification",
          "qualification",
          "text",
          { required: "Qualification is required" },
          "E.g., Bachelor's in Computer Science, Master's in Business."
        )}
        {renderInput(
          "Institution Name",
          "institutionName",
          "text",
          { required: "Institution name is required" },
          "Enter the name of the university or school."
        )}
        {renderInput(
          "Graduation Year",
          "graduationYear",
          "number",
          {
            required: "Graduation year is required",
            min: { value: 1950, message: "Year must be after 1950" },
            max: { value: new Date().getFullYear(), message: "Year cannot be in the future" },
          },
          "Enter the year you officially graduated (e.g., 2022)."
        )}
        {renderInput(
          "GPA Score",
          "gpaScore",
          "text",
          {},
          "Enter your GPA (e.g., 3.5/4 or 15/20)."
        )}
        {renderInput(
          "Language Test Type",
          "languageTestType",
          "text",
          {},
          "E.g., IELTS, TOEFL, or none if not applicable."
        )}
        {renderInput(
          "Language Test Score",
          "languageTestScore",
          "text",
          {},
          "Enter your score (e.g., 7.5 for IELTS)."
        )}

        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-400 mb-1">
            Certificate File (PDF/Image)
          </label>
          <input
            type="file"
            accept=".pdf,image/*"
            onChange={(e) => {
              const file = e.currentTarget.files?.[0] || null;
              setSelectedFile(file);
              setValue("certificate", file as any);
            }}
            className="w-full px-3 py-2 rounded-md bg-gray-800 border border-gray-700 text-gray-100 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p className="text-xs text-gray-400 mt-1">
            Upload a scanned copy of your degree/diploma or a valid academic certificate.
          </p>
          {selectedFile && (
            <p className="text-sm text-gray-500 mt-1">
              Selected file: <span className="text-gray-300">{selectedFile.name}</span>
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
