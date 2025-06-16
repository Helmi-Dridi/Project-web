import type { StudentAcademicBackground } from "../services/profile.service";

type Props = {
  data?: StudentAcademicBackground;
};

export default function AcademicBackgroundDisplay({ data }: Props) {
  const render = (label: string, value?: string | number | null) => (
    <div className="flex flex-wrap gap-2 text-sm">
      <span className="w-48 font-medium text-gray-400">{label}:</span>
      <span className="text-gray-100 break-words">
        {value ? value : <span className="italic text-gray-500">Not provided</span>}
      </span>
    </div>
  );

  const renderCertificate = (path?: string | null) => {
    if (!path) return render("Certificate File", null);

    const fileUrl = `http://localhost:8080${path}`;
    return (
      <div className="flex flex-wrap gap-2 text-sm">
        <span className="w-48 font-medium text-gray-400">Certificate File:</span>
        <a
          href={fileUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="text-blue-400 hover:text-blue-600 underline break-words"
        >
          View Certificate
        </a>
      </div>
    );
  };

  return (
    <div className="space-y-3">
      {render("Qualification", data?.qualification)}
      {render("Institution", data?.institutionName)}
      {render("Graduation Year", data?.graduationYear)}
      {render("GPA", data?.gpaScore)}
      {render("Language Test Type", data?.languageTestType)}
      {render("Language Test Score", data?.languageTestScore)}
      {renderCertificate(data?.certificateFilePath)}
    </div>
  );
}
