"use client";

import type { JobType, UserRole } from "../../types/ProfileData";
import { useState, useRef, useCallback, useEffect } from "react";
import {
  useGetProfileQuery,
  useGetMyProfilePictureQuery,
  useUpdateProfileMutation,
  useDeleteAccountMutation,
  useUploadProfilePictureMutation,
  useUploadDocumentMutation,
  useDeleteFileMutation,
} from "@/lib/redux/slices/ProfileApiSlice";
import Cropper from "react-easy-crop";
import type { Area } from "react-easy-crop";
import {
  ArrowLeft,
  Download,
  Upload,
  Trash2,
  Edit3,
  Menu,
  X,
} from "lucide-react";
import { Card, CardHeader, CardContent } from "../components/customCard";
import { Button } from "@/app/components/ui/button";
import Footer from "../components/footer";
import Image from "next/image";
import Link from "next/link";

function getCroppedImg(imageSrc: string, crop: Area): Promise<string> {
  return new Promise((resolve, reject) => {
    const image = new window.Image();
    image.crossOrigin = "anonymous";

    image.onload = () => {
      const w = Math.round(crop.width);
      const h = Math.round(crop.height);
      const sx = Math.round(crop.x);
      const sy = Math.round(crop.y);

      const canvas = document.createElement("canvas");
      canvas.width = w;
      canvas.height = h;

      const ctx = canvas.getContext("2d");
      if (!ctx) return reject(new Error("Failed to get canvas 2D context."));

      ctx.drawImage(image, sx, sy, w, h, 0, 0, w, h);

      if (canvas.toBlob) {
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              try {
                return resolve(canvas.toDataURL("image/jpeg", 0.92));
              } catch {
                return reject(new Error("Canvas is empty"));
              }
            }
            resolve(URL.createObjectURL(blob));
          },
          "image/jpeg",
          0.92
        );
      } else {
        try {
          resolve(canvas.toDataURL("image/jpeg", 0.92));
        } catch (err) {
          reject(err as Error);
        }
      }
    };

    image.onerror = (error) => reject(error as string);
    image.src = imageSrc;
  });
}

export default function ProfilePage() {
  // Demo data for fallback
  const demoData = {
    userName: "Random User",
    role: "Admin",
    fullName: "Random User",
    email: "RandomUser@gmail.com",
    phone: "",
    location: "",
    jobType: "Remote",
    country: "Ethiopia",
    city: "Addis Ababa",
    skills: [
      { name: "AI Modeling" },
      { name: "Data Analysis" },
      { name: "Full-Stack Dev" },
    ],
    bio: "This is a sample bio for the user. Update your profile to add more details.",
    profilePic: "/person.png",
  };

  const { data: profile } = useGetProfileQuery();
  const { data: profilePicBlob } = useGetMyProfilePictureQuery();

  const [profilePic, setProfilePic] = useState<string>(demoData.profilePic);
  useEffect(() => {
    if (profilePicBlob) {
      const url = URL.createObjectURL(profilePicBlob);
      setProfilePic(url);
      return () => URL.revokeObjectURL(url);
    } else {
      setProfilePic(demoData.profilePic);
    }
  }, [profilePicBlob, demoData.profilePic]);

  const [bio, setBio] = useState(demoData.bio);
  useEffect(() => {
    setBio(profile?.bio || demoData.bio);
  }, [profile?.bio, demoData.bio]);
  const BIO_WORD_LIMIT = 50;
  const getBioWordCount = (text: string) =>
    text.trim().split(/\s+/).filter(Boolean).length;
  const [isProfilePicHover, setIsProfilePicHover] = useState(false);
  const profilePicInputRef = useRef<HTMLInputElement>(null);
  const [cvFile, setCvFile] = useState<File | null>(null);
  const [cvUrl, setCvUrl] = useState<string>("");
  const [isCvHover, setIsCvHover] = useState(false);
  const cvInputRef = useRef<HTMLInputElement>(null);

  const [cropModalOpen, setCropModalOpen] = useState(false);
  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const [crop, setCrop] = useState({ x: 0, y: 0 });
  const [zoom, setZoom] = useState(1);
  const [croppedAreaPixels, setCroppedAreaPixels] = useState<Area | null>(null);

  const [updateProfile, { isLoading: isUpdatingProfile }] =
    useUpdateProfileMutation();
  const [deleteAccount, { isLoading: isDeletingAccount }] =
    useDeleteAccountMutation();
  const [uploadProfilePicture] = useUploadProfilePictureMutation();
  const [uploadDocument, { isLoading: isUploadingDocument }] =
    useUploadDocumentMutation();
  const [deleteFile, { isLoading: isDeletingFile }] = useDeleteFileMutation();

  const handleProfilePicChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];

    if (file) {
      const url = URL.createObjectURL(file);
      setSelectedImage(url);
      setCropModalOpen(true);
    }
    e.target.value = "";
  };

  const onCropComplete = useCallback((_: Area, croppedAreaPixels: Area) => {
    setCroppedAreaPixels(croppedAreaPixels);
  }, []);

  const handleCropSave = async () => {
    if (selectedImage && croppedAreaPixels) {
      const croppedUrl = await getCroppedImg(selectedImage, croppedAreaPixels);
      let blob: Blob | null = null;
      if (croppedUrl.startsWith("blob:")) {
        blob = await fetch(croppedUrl).then((res) => res.blob());
      } else if (croppedUrl.startsWith("data:")) {
        const arr = croppedUrl.split(",");
        const mime = arr[0].match(/:(.*?);/)?.[1] || "image/jpeg";
        const bstr = atob(arr[1]);
        let n = bstr.length;
        const u8arr = new Uint8Array(n);
        while (n--) {
          u8arr[n] = bstr.charCodeAt(n);
        }
        blob = new Blob([u8arr], { type: mime });
      }
      if (blob) {
        const formData = new FormData();
        formData.append("file", blob, "profile.jpg");
        try {
          await uploadProfilePicture(formData).unwrap();
        } catch {
          alert("Failed to upload profile picture.");
        }
      }
      setProfilePic(croppedUrl);
      setCropModalOpen(false);
      URL.revokeObjectURL(selectedImage);
      setSelectedImage(null);
    }
  };
  const handleCropCancel = () => {
    if (selectedImage) {
      URL.revokeObjectURL(selectedImage);
    }
    setSelectedImage(null);
    setCropModalOpen(false);
  };

  const handleCvUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file && file.type === "application/pdf") {
      setCvFile(file);
      setCvUrl(URL.createObjectURL(file));
      const formData = new FormData();
      formData.append("file", file, file.name);
      try {
        await uploadDocument(formData).unwrap();
      } catch {
        alert("Failed to upload CV.");
      }
    } else {
      alert("Please upload a PDF file.");
    }
    e.target.value = "";
  };

  const handleCvDelete = async () => {
    if (cvFile) {
      try {
        await deleteFile(cvFile.name).unwrap();
      } catch {
        alert("Failed to delete CV.");
      }
    }
    setCvFile(null);
    setCvUrl("");
  };
  const handleUpdateProfile = async () => {
    const payload = {
      username: formData.userName,
      role: formData.role as UserRole,
      full_name: formData.fullName,
      email: formData.email,
      phone_number: formData.phone,
      location: formData.location,
      job_type: formData.jobType as JobType,
      preferred_country: formData.country,
      city_region: formData.city,
      skills: skills.map((s) => s.name),
      bio: bio,
    };
    try {
      await updateProfile(payload).unwrap();
      alert("Profile updated successfully!");
    } catch {
      alert("Failed to update profile.");
    }
  };

  const handleDeleteAccount = async () => {
    try {
      await deleteAccount().unwrap();
      alert("Account deleted successfully!");
    } catch {
      alert("Failed to delete account.");
    }
  };

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  // Form data with fallback to demo data
  const [formData, setFormData] = useState({ ...demoData });
  useEffect(() => {
    setFormData({
      userName: profile?.username || demoData.userName,
      role: profile?.role || demoData.role,
      fullName: profile?.full_name || demoData.fullName,
      email: profile?.email || demoData.email,
      phone: profile?.phone_number || demoData.phone,
      location: profile?.location || demoData.location,
      jobType: profile?.job_type || demoData.jobType,
      country: profile?.preferred_country || demoData.country,
      city: profile?.city_region || demoData.city,
      skills:
        profile?.skills && profile.skills.length > 0
          ? profile.skills.map((s: string) => ({ name: s }))
          : demoData.skills,
      bio: profile?.bio || demoData.bio,
      profilePic: profilePic || demoData.profilePic,
    });
  }, [profile, profilePic]);

  const [skills, setSkills] = useState<{ name: string }[]>(demoData.skills);
  useEffect(() => {
    if (profile?.skills && profile.skills.length > 0) {
      setSkills(profile.skills.map((s: string) => ({ name: s })));
    } else {
      setSkills(demoData.skills);
    }
  }, [profile?.skills]);
  const [editingSkills, setEditingSkills] = useState(false);

  const cityOptions = {
    Ethiopia: [
      "Addis Ababa",
      "Dire Dawa",
      "Bahir Dar",
      "Hawassa",
      "Mekelle",
      "Gondar",
    ],
    USA: ["New York", "Los Angeles", "Washington D.C.", "Atlanta"],
  };

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  const closeMobileMenu = () => {
    setIsMobileMenuOpen(false);
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => {
      if (field === "country") {
        const newCountry = value as keyof typeof cityOptions;
        return {
          ...prev,
          [field]: value,
          city: cityOptions[newCountry][0],
        };
      }
      return {
        ...prev,
        [field]: value,
      };
    });
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col">
      {/* Header */}
      <header className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center">
                <div className="w-6 h-6 bg-[#44C3BB] rounded-full flex items-center justify-center">
                  <div className="flex gap-1">
                    <div className="w-1 h-1 bg-white rounded-full"></div>
                    <div className="w-1 h-1 bg-white rounded-full"></div>
                  </div>
                </div>
              </div>
              <span className="text-xl font-semibold text-white dark:text-gray-100">
                JobGen
              </span>
            </div>

            {/* Desktop Navigation */}
            <nav className="hidden md:flex items-center gap-8">
              <a href="#" className="hover:text-gray-200 transition-colors">
                Blog
              </a>
              <a href="#" className="hover:text-gray-200 transition-colors">
                About
              </a>
              <a href="#" className="hover:text-gray-200 transition-colors">
                Contact
              </a>
              <div className="w-8 h-8 rounded-full overflow-hidden cursor-pointer transition-all duration-200 ring-0 hover:ring-2 hover:ring-[#44C3BB]/60 hover:scale-150 z-50">
                <Image
                  src={profilePic}
                  alt="Profile"
                  className="w-full h-full object-cover"
                  width={48}
                  height={48}
                  unoptimized
                />
              </div>
            </nav>

            {/* Mobile Menu Button */}
            <button
              onClick={toggleMobileMenu}
              className="md:hidden p-2 rounded-md hover:bg-white/10 transition-colors cursor-pointer"
            >
              {isMobileMenuOpen ? (
                <X className="w-6 h-6" />
              ) : (
                <Menu className="w-6 h-6" />
              )}
            </button>
          </div>

          {/* Mobile Navigation */}
          {isMobileMenuOpen && (
            <div className="md:hidden py-4 border-t border-white/20">
              <div className="flex flex-col gap-4">
                <a
                  href="#"
                  onClick={closeMobileMenu}
                  className="hover:text-gray-200 transition-colors"
                >
                  Blog
                </a>
                <a
                  href="#"
                  onClick={closeMobileMenu}
                  className="hover:text-gray-200 transition-colors"
                >
                  About
                </a>
                <a
                  href="#"
                  onClick={closeMobileMenu}
                  className="hover:text-gray-200 transition-colors"
                >
                  Contact
                </a>
              </div>
            </div>
          )}
        </div>
      </header>

      {/* Back to Dashboard */}
      <div className="max-w-7xl mx-auto my-4 px-4 mt-8 sm:px-6 lg:px-8 py-6 w-full flex">
        <div className="w-full flex justify-start">
          <Link
            // href="/dashboard"
            href="/settings"
            className="flex items-center gap-2 text-gray-600 dark:text-gray-300 hover:text-gray-800 dark:hover:text-white transition-colors cursor-pointer"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Dashboard
          </Link>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-12 flex-1">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Profile Information */}
          <div className="lg:col-span-1">
            <Card className="shadow-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
              <CardContent className="p-6">
                {/* Profile Photo */}
                <div className="flex flex-col items-center mb-6">
                  <div
                    className="w-32 h-32 rounded-full overflow-hidden mb-4 relative group"
                    onMouseEnter={() => setIsProfilePicHover(true)}
                    onMouseLeave={() => setIsProfilePicHover(false)}
                  >
                    <Image
                      src={profilePic}
                      alt="Profile"
                      className="w-full h-full object-cover"
                      width={128}
                      height={128}
                      unoptimized
                    />
                    {isProfilePicHover && (
                      <div className="absolute inset-0 bg-black/50 flex items-center justify-center cursor-pointer transition-opacity z-10">
                        <Button
                          type="button"
                          className="z-20"
                          onClick={(e: React.MouseEvent) => {
                            e.stopPropagation();
                            profilePicInputRef.current?.click();
                          }}
                        >
                          Update
                        </Button>
                      </div>
                    )}
                    <input
                      ref={profilePicInputRef}
                      type="file"
                      accept="image/*"
                      className="hidden"
                      onChange={handleProfilePicChange}
                    />
                    <div
                      className="w-32 h-32 rounded-full overflow-hidden mb-4 relative group"
                      onMouseEnter={() => setIsProfilePicHover(true)}
                      onMouseLeave={() => setIsProfilePicHover(false)}
                      onClick={() => profilePicInputRef.current?.click()}
                    />
                    {/* Crop Modal */}
                    {cropModalOpen && selectedImage && (
                      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
                        <div className="bg-white dark:bg-gray-900 p-6 rounded-lg shadow-lg relative w-[90vw] max-w-md flex flex-col items-center">
                          <div className="w-64 h-64 relative bg-gray-200 rounded-full overflow-hidden">
                            <Cropper
                              image={selectedImage}
                              crop={crop}
                              zoom={zoom}
                              aspect={1}
                              cropShape="round"
                              showGrid={false}
                              onCropChange={setCrop}
                              onZoomChange={setZoom}
                              onCropComplete={onCropComplete}
                            />
                          </div>
                          <div className="flex gap-4 mt-4">
                            <Button onClick={handleCropSave}>Save</Button>
                            <Button
                              variant="destructive"
                              onClick={handleCropCancel}
                            >
                              Cancel
                            </Button>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                  <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
                    {formData.userName}
                  </h2>
                  <p className="text-gray-600 dark:text-gray-400">
                    {formData.role}
                  </p>
                </div>

                {/* Personal Information */}
                <div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
                    Personal Information
                  </h3>

                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Full Name
                      </label>
                      <input
                        type="text"
                        value={formData.fullName || demoData.fullName}
                        onChange={(e) =>
                          handleInputChange("fullName", e.target.value)
                        }
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#2C8C86] dark:focus:ring-[#2C8C86] focus:border-transparent"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Email Address
                      </label>
                      <input
                        type="email"
                        value={formData.email || demoData.email}
                        onChange={(e) =>
                          handleInputChange("email", e.target.value)
                        }
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#2C8C86] dark:focus:ring-[#2C8C86] focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Phone Number
                      </label>
                      <input
                        type="tel"
                        value={formData.phone || demoData.phone}
                        onChange={(e) =>
                          handleInputChange("phone", e.target.value)
                        }
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#2C8C86] dark:focus:ring-[#2C8C86] focus:border-transparent"
                        placeholder="e.g. +251 12 34 6789"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Location
                      </label>
                      <input
                        type="text"
                        value={formData.location || demoData.location}
                        onChange={(e) =>
                          handleInputChange("location", e.target.value)
                        }
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#2C8C86] dark:focus:ring-[#2C8C86] focus:border-transparent"
                        placeholder="e.g. Addis Ababa, Ethiopia"
                      />
                    </div>
                  </div>
                </div>
                <div className="flex gap-4 mt-6">
                  <Button
                    onClick={handleUpdateProfile}
                    disabled={isUpdatingProfile}
                  >
                    {isUpdatingProfile ? "Updating..." : "Update Account"}
                  </Button>
                  <Button
                    variant="destructive"
                    onClick={handleDeleteAccount}
                    disabled={isDeletingAccount}
                  >
                    {isDeletingAccount ? "Deleting..." : "Delete Account"}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Middle Column */}
          <div className="lg:col-span-1 flex flex-col h-full justify-between">
            {/* CV Section */}
            <Card className="shadow-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 mb-2">
              <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0">
                <div className="flex items-center gap-3 p-6">
                  <h3 className="text-lg font-semibold">User CV</h3>
                </div>
              </CardHeader>
              <CardContent className="p-6">
                <div className="flex flex-col items-center py-8">
                  <div
                    className="w-16 h-16 border-2 border-gray-300 rounded-lg flex items-center justify-center mb-4 relative group"
                    onMouseEnter={() => setIsCvHover(true)}
                    onMouseLeave={() => setIsCvHover(false)}
                  >
                    <Download className="w-8 h-8 text-gray-400" />
                    {isCvHover && cvFile && (
                      <div
                        className="absolute inset-0 flex items-center justify-center bg-black/60 text-white rounded-lg transition-opacity"
                        style={{ zIndex: 10 }}
                      >
                        <Button
                          type="button"
                          onClick={() => {
                            setTimeout(() => {
                              const link = document.createElement("a");
                              link.href = cvUrl;
                              link.download = cvFile.name;
                              document.body.appendChild(link);
                              link.click();
                              document.body.removeChild(link);
                            }, 500);
                          }}
                        >
                          Download
                        </Button>
                      </div>
                    )}
                  </div>
                  <p className="text-gray-900 dark:text-gray-100 font-medium mb-6">
                    {cvFile ? cvFile.name : "Professional CV"}
                  </p>
                  <div className="flex gap-3">
                    <Button
                      type="button"
                      onClick={() => cvInputRef.current?.click()}
                      disabled={isUploadingDocument}
                    >
                      <Upload className="w-4 h-4 mr-2" />
                      {isUploadingDocument ? "Uploading..." : "Upload New"}
                    </Button>
                    <input
                      ref={cvInputRef}
                      type="file"
                      accept="application/pdf"
                      className="hidden"
                      onChange={handleCvUpload}
                    />
                    <Button
                      type="button"
                      onClick={handleCvDelete}
                      disabled={!cvFile || isDeletingFile}
                    >
                      <Trash2 className="w-4 h-4 mr-2" />
                      {isDeletingFile ? "Deleting..." : "Delete"}
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Skills Section */}
            <Card className="shadow-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 min-h-[288.9px]">
              <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0">
                <div className="flex items-center gap-3 p-6">
                  <h3 className="text-lg font-semibold">Skill Set</h3>
                </div>
              </CardHeader>
              <CardContent className="p-6">
                {skills.length === 0 && (
                  <div className="mb-4 p-4 rounded bg-gray-100 dark:bg-gray-800 text-gray-500 dark:text-gray-400 text-center border border-dashed border-gray-300 dark:border-gray-700">
                    No skills added yet. Click the{" "}
                    <span className="font-bold text-green-600">+</span> button
                    to add your first skill!
                  </div>
                )}
                <div className="grid grid-cols-2 gap-4 mb-6">
                  {skills.map((skill, idx) => (
                    <div key={idx} className="flex items-center gap-2">
                      {editingSkills ? (
                        <>
                          <input
                            type="text"
                            value={skill.name}
                            onChange={(e) => {
                              const newSkills = [...skills];
                              newSkills[idx].name = e.target.value;
                              setSkills(newSkills);
                            }}
                            className="font-semibold text-gray-900 dark:text-gray-100 w-full mb-1 px-2 py-1 rounded border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-900"
                          />
                          <button
                            className="ml-1 text-red-500 hover:text-red-700 text-lg font-bold px-2 cursor-pointer"
                            onClick={() =>
                              setSkills(skills.filter((_, i) => i !== idx))
                            }
                            aria-label="Remove skill"
                          >
                            x
                          </button>
                        </>
                      ) : (
                        <div className="font-semibold text-gray-900 dark:text-gray-100">
                          {skill.name}
                        </div>
                      )}
                    </div>
                  ))}
                  {editingSkills && (
                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        className="text-green-600 hover:text-green-800 text-2xl font-bold px-2 cursor-pointer"
                        onClick={() => setSkills([...skills, { name: "" }])}
                        aria-label="Add skill"
                      >
                        +
                      </button>
                    </div>
                  )}
                </div>
                <div className="flex justify-center">
                  <Button onClick={() => setEditingSkills(!editingSkills)}>
                    <Edit3 className="w-4 h-4 mr-2" />
                    {editingSkills ? "Save Skills" : "Edit Skills"}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Right Column - Preference and Bio */}
          <div className="lg:col-span-1 flex flex-col h-full">
            <div className="flex flex-col gap-6 h-full justify-between">
              <Card className="shadow-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 pb-3">
                <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0">
                  <div className="flex items-center gap-3 p-6">
                    <h3 className="text-lg font-semibold">Preference</h3>
                  </div>
                </CardHeader>
                <CardContent className="p-6">
                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Job Type
                      </label>
                      <select
                        value={formData.jobType || demoData.jobType}
                        onChange={(e) =>
                          handleInputChange("jobType", e.target.value)
                        }
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent"
                      >
                        <option>Full Time</option>
                        <option>Part Time</option>
                        <option>Contract</option>
                        <option>Internship</option>
                        <option>Temporary</option>
                        <option>Remote</option>
                        <option>Hybrid</option>
                        <option>Free Lance</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Preferred Country
                      </label>
                      <select
                        value={formData.country || demoData.country}
                        onChange={(e) =>
                          handleInputChange("country", e.target.value)
                        }
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent"
                      >
                        <option value="Ethiopia">Ethiopia</option>
                        <option value="USA">USA</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        City/Region
                      </label>
                      <select
                        value={formData.city || demoData.city}
                        onChange={(e) =>
                          handleInputChange("city", e.target.value)
                        }
                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent"
                      >
                        {(
                          cityOptions[
                            (formData.country as keyof typeof cityOptions) ||
                              demoData.country
                          ] || []
                        ).map((city) => (
                          <option key={city} value={city}>
                            {city}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>
                </CardContent>
              </Card>
              {/* Bio Card */}
              <Card className="shadow-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
                <CardHeader className="bg-gradient-to-r from-[#44C3BB] to-[#3AB5AD] dark:from-gray-900 dark:to-gray-800 text-white dark:text-gray-100 p-0">
                  <div className="flex items-center gap-3 p-6">
                    <h3 className="text-lg font-semibold">Bio</h3>
                  </div>
                </CardHeader>
                <CardContent className="p-6">
                  <textarea
                    value={bio}
                    onChange={(e) => {
                      const words = e.target.value
                        .trim()
                        .split(/\s+/)
                        .filter(Boolean);
                      if (words.length <= BIO_WORD_LIMIT) {
                        setBio(e.target.value);
                      } else {
                        setBio(words.slice(0, BIO_WORD_LIMIT).join(" "));
                      }
                    }}
                    rows={5}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-700 rounded-md bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-[#44C3BB] focus:border-transparent resize-none"
                    placeholder={`Write a short bio about yourself... (max ${BIO_WORD_LIMIT} words)`}
                  />
                  <div className="text-right text-xs text-gray-500 mt-1">
                    {getBioWordCount(bio)} / {BIO_WORD_LIMIT} words
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
      {/* Informative text */}
      <div className="w-full flex justify-center mt-8 mb-2">
        <div className="max-w-md w-full text-center text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 rounded-lg px-4 py-3 border border-dashed border-gray-300 dark:border-gray-700 shadow-sm">
          <span className="font-medium text-gray-700 dark:text-gray-200">
            When you&apos;re done making changes, click the{" "}
            <span className="font-semibold text-[#44C3BB]">Update Account</span>{" "}
            button below your profile to save your changes.
          </span>
        </div>
      </div>
      <Footer />
    </div>
  );
}
