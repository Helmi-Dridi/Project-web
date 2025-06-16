import type { StudentProfile } from "../services/profile.service";

type Props = {
  data?: StudentProfile;
};

export default function PersonalInfoDisplay({ data }: Props) {
  const render = (label: string, value?: string | null) => (
    <div className="flex flex-wrap sm:flex-nowrap gap-2 text-sm">
      <span className="w-48 font-medium text-gray-400">{label}:</span>
      <span className="text-gray-100 break-words">
        {value && value.trim() !== "" ? (
          value
        ) : (
          <span className="text-gray-500 italic">Not provided</span>
        )}
      </span>
    </div>
  );

  if (!data) {
    return <p className="text-gray-500 italic">No personal information available.</p>;
  }

  return (
    <div className="space-y-3">
      {render("Date of Birth", data.dateOfBirth?.slice(0, 10))}
      {render("Gender", data.gender)}
      {render("Nationality", data.nationality)}
      {render("Passport Number", data.passportNumber)}
      {render("National ID", data.nationalId)}
    </div>
  );
}
