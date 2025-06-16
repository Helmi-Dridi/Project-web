import { useEffect, useState } from "react";

interface SplashScreenProps {
  onComplete: () => void;
}

export default function SplashScreen({ onComplete }: SplashScreenProps) {
  const [isSliding, setIsSliding] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsSliding(true);
      setTimeout(onComplete, 800); // Wait for animation to finish
    }, 3000);

    return () => clearTimeout(timer);
  }, [onComplete]);

  return (
    <div
      className={`fixed inset-0 bg-gradient-to-br from-slate-50 via-blue-100 to-white z-50 flex flex-col items-center justify-center transition-transform duration-700 ease-in-out ${
        isSliding ? "-translate-y-full" : "translate-y-0"
      }`}
    >
      <img
        src="/lovable-uploads/Scolarev.jpg"
        alt="ScholarRev Logo"
        className="h-32 w-auto mb-4 animate-pulse"
      />
      <h1 className="text-2xl md:text-3xl font-semibold text-slate-700 animate-fade-in">
        Welcome to ScholarRev
      </h1>
    </div>
  );
}
