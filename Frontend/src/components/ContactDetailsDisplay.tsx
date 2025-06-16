import type { StudentContactDetails } from "../services/profile.service";

type Props = {
  data?: StudentContactDetails;
};

export default function ContactDetailsDisplay({ data }: Props) {
  const render = (label: string, value?: string | null) => (
    <div className="flex flex-col sm:flex-row sm:items-start text-sm gap-1 sm:gap-2">
      <span className="w-48 font-medium text-gray-400">{label}:</span>
      <span className="text-gray-100 break-words">
        {value ? value : <span className="italic text-gray-500">Not provided</span>}
      </span>
    </div>
  );

  if (!data) {
    return <p className="text-gray-500 italic">No contact information available.</p>;
  }

  return (
    <div className="space-y-3">
      {render("Email", data.email)}
      {render("Phone Number", data.phoneNumber)}
      {render("Address Line", data.addressLine)}
      {render("City", data.city)}
      {render("Country", data.country)}
      {render("Zip Code", data.zipCode)}
      {render("Preferred Contact", data.preferredContactMethod)}
    </div>
  );
}
