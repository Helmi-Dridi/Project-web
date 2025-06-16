import { useRef } from "react";
import { useAuth } from "../context/AuthContext";
import { Camera } from "lucide-react";

type Props = {
  onImageUpload: (file: File) => Promise<string>;
};

export default function UserProfileCard({ onImageUpload }: Props) {
  const { user, setUser } = useAuth();
  const inputRef = useRef<HTMLInputElement | null>(null);

  const handleImageChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !user) return;

    try {
      const newUrl = await onImageUpload(file);
      const updatedUser = { ...user, profilePicture: newUrl };
      localStorage.setItem("user", JSON.stringify(updatedUser));
      setUser(updatedUser);
    } catch (err) {
      console.error("Failed to upload avatar:", err);
      alert("Image upload failed.");
    }
  };

  return (
    <section className="bg-gray-800 border border-gray-700 rounded-2xl shadow p-6 w-full flex items-center gap-6">
      <div className="relative group w-24 h-24">
        <img
          src={user?.profilePicture || "/default-avatar.png"}
          alt="Profile"
          className="w-24 h-24 rounded-full object-cover border-2 border-gray-700"
        />
        <div
          onClick={() => inputRef.current?.click()}
          className="absolute inset-0 bg-black/40 flex items-center justify-center rounded-full opacity-0 group-hover:opacity-100 transition cursor-pointer"
          title="Change profile picture"
        >
          <Camera className="w-5 h-5 text-white" />
        </div>
        <input
          ref={inputRef}
          type="file"
          accept="image/*"
          onChange={handleImageChange}
          className="hidden"
        />
      </div>

      <div>
        <h2 className="text-2xl font-semibold text-white">
          {user?.name || "Unknown User"}
        </h2>
        <p className="text-gray-400 text-sm">
          This is your public profile image and display name.
        </p>
      </div>
    </section>
  );
}
