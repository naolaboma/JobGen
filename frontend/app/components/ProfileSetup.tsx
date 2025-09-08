"use client";

import { useEffect, useState } from "react";

export default function ProfileSetup() {
  const [file, setFile] = useState<File | null>(null);
  const [skills, setSkills] = useState<string[]>([]);
  const [uploading, setUploading] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  const token =
    "..";

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await fetch("http://localhost:8080/api/v1/users/profile", {
          method: "GET",
          headers: {
            Authorization: `Bearer ${token}`,
            Accept: "application/json",
          },
        });

        const data = await res.json();
        const user = data?.data?.user;

        if (user?.skills) {
          setSkills(user.skills);
        }
      } catch (error) {
        console.error("Failed to load user profile:", error);
        setMessage("Failed to load user profile.");
      }
    };

    fetchProfile();
  }, [token]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];

      const maxSizeMB = 5;
      if (selectedFile.size > maxSizeMB * 1024 * 1024) {
        setMessage(`File too large (max ${maxSizeMB}MB).`);
        return;
      }

      setFile(selectedFile);
      setMessage(null);
    }
  };

  const uploadCV = async () => {
    if (!file) {
      setMessage("Please select a file first.");
      return;
    }

    try {
      setUploading(true);
      setMessage(null);

      const formData = new FormData();
      formData.append("file", file);

      const res = await fetch("http://localhost:8080/api/v1/cv/parse", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
          Accept: "application/json",
        },
        body: formData,
      });

      const data = await res.json();

      console.log("CV parsing response:", data);
      setMessage("CV parsing job started successfully.");
    } catch (error: any) {
      console.error("Upload failed:", error);
      setMessage("CV parsing failed.");
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-white px-4 py-8">
      <h1 className="text-2xl font-bold text-teal-700 mb-6 text-center">
        Set Up Your Profile
      </h1>

      <label
        htmlFor="cv-upload"
        className="flex flex-col items-center justify-center w-full max-w-md h-48 border-2 border-gray-300 rounded-lg cursor-pointer bg-gray-100 hover:bg-gray-200 transition"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-10 w-10 text-gray-600 mb-2"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5-5m0 0l5 5m-5-5v12"
          />
        </svg>
        <p className="text-gray-600 text-sm text-center px-2 truncate">
          {file ? file.name : "Upload your CV"}
        </p>
        <input
          id="cv-upload"
          type="file"
          className="hidden"
          onChange={handleFileChange}
        />
      </label>

      <button
        className="mt-4 w-full max-w-md bg-teal-600 text-white py-2 rounded-lg shadow hover:bg-teal-700 transition"
        onClick={uploadCV}
        disabled={uploading}
      >
        {uploading ? "Uploading..." : "+ Upload CV"}
      </button>

      {message && (
        <p className="mt-4 text-sm text-center text-gray-700">{message}</p>
      )}

      <div className="mt-6 w-full max-w-md">
        <h2 className="font-semibold text-gray-800 mb-3">Skills</h2>
        <div className="flex flex-wrap gap-3">
          {skills.length === 0 ? (
            <p className="text-gray-500 text-sm">No skills found</p>
          ) : (
            skills.map((skill, index) => (
              <span
                key={index}
                className="px-4 py-1 border rounded-full text-sm font-medium text-teal-700 border-teal-600 bg-teal-50"
              >
                {skill}
              </span>
            ))
          )}
        </div>
      </div>

      <button className="mt-8 text-teal-600 font-semibold hover:underline">
        Next →
      </button>
    </div>
  );
}
