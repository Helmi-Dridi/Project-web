import { useState } from "react";
import LoginPage from "../components/LoginPage"; // move your current login form here
import SplashScreen from "../components/SplashScreen";

export default function IndexPage() {
  const [showSplash, setShowSplash] = useState(true);

  return (
    <>
      {showSplash && <SplashScreen onComplete={() => setShowSplash(false)} />}
      <div className={`${showSplash ? "opacity-0" : "opacity-100"} transition-opacity duration-500`}>
        <LoginPage />
      </div>
    </>
  );
}
