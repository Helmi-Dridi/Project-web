import type { StudentSettings } from "../services/profile.service";

type Props = {
  data?: StudentSettings;
};

export default function SettingsDisplay({ data }: Props) {
  const render = (label: string, value?: string | boolean | null) => (
    <div className="flex flex-wrap text-sm gap-2">
      <span className="w-48 font-medium text-gray-400">{label}:</span>
      <span className="text-gray-100 break-words">
        {typeof value === "boolean"
          ? value
            ? "Yes"
            : "No"
          : value || <span className="italic text-gray-500">Not provided</span>}
      </span>
    </div>
  );

  if (!data) {
    return <p className="text-gray-500 italic">No settings information available.</p>;
  }

  return (
    <div className="space-y-3">
      {render("Language Preference", data.languagePreference)}
      {render("Email Notifications", data.receiveEmailNotifications)}
      {render("SMS Notifications", data.receiveSMSNotifications)}
      {render("Account Deleted", data.accountDeleted)}
      {render("Emergency Name", data.emergencyName)}
      {render("Relationship", data.emergencyRelationship)}
      {render("Phone", data.emergencyPhoneNumber)}
      {render("Email", data.emergencyEmail)}
    </div>
  );
}
