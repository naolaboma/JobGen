// app/components/JobCard.tsx
"use client";

import { FC } from "react";
import { MapPin, Building2, Bookmark, Share2 } from "lucide-react";

interface JobCardProps {
  title: string;
  company: string;
  location: string;
  tags: string[];
  salary: string;
  posted: string;
  deadline: string;
  applicants: number;
}

const JobCard: FC<JobCardProps> = ({
  title,
  company,
  location,
  tags,
  salary,
  posted,
  deadline,
  applicants,
}) => {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center p-4 sm:p-6 text-black font-inter">
      {/* Header Row */}
      <div className="w-full max-w-5xl grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Left side: Job Title & Meta */}
        <div className="p-6 md:col-span-2">
          <h1 className="text-2xl sm:text-3xl font-bold">{title}</h1>
          <div className="flex items-center gap-2 mt-2 text-base">
            <Building2 size={18} />
            <span>{company}</span>
          </div>
          <div className="flex items-center gap-2 mt-1 text-base">
            <MapPin size={18} />
            <span>{location}</span>
          </div>

          <div className="flex flex-wrap gap-2 mt-3">
            {tags.map((tag, i) => (
              <span
                key={i}
                className="px-3 py-1 text-base bg-gray-100 border rounded-full"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>

        {/* Right side: Job overview */}
        <div className="bg-white rounded-2xl shadow p-6">
          <h3 className="text-lg sm:text-xl font-semibold">Job overview</h3>
          <p className="text-base font-bold mt-2">{salary}</p>
          <p className="text-base">Full-time</p>
          <p className="text-base mt-2">Posted {posted}</p>
          <p className="text-base">Deadline: {deadline}</p>
          <p className="text-base">{applicants} applicants</p>
        </div>
      </div>

      {/* Body Row */}
      <div className="w-full max-w-5xl mt-6 grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Job Description */}
        <div className="col-span-2 bg-white rounded-2xl shadow p-6">
          <h2 className="text-lg sm:text-2xl font-semibold mb-3">Job Description</h2>
          <p className="text-base mb-4">
            We are seeking for a <strong>Frontend Developer</strong> to build
            and maintain user interfaces using React, collaborate with designers
            and backend developers, and optimize applications for maximum speed
            and scalability.
          </p>

          <h3 className="text-lg sm:text-xl font-semibold">Responsibilities</h3>
          <ul className="list-disc list-inside mb-4 text-base">
            <li>Develop new user-facing features</li>
            <li>Build reusable code and libraries</li>
            <li>Ensure feasibility of UI/UX designs</li>
            <li>Optimize application for maximum speed</li>
          </ul>

          <h3 className="text-lg sm:text-xl font-semibold">Requirements</h3>
          <ul className="list-disc list-inside mb-4 text-base">
            <li>3+ years of experience in JavaScript</li>
            <li>Knowledge of React, Redux, REST APIs</li>
            <li>Strong HTML5, web markup skills</li>
            <li>Responsive & mobile design experience</li>
          </ul>

          <h3 className="text-lg sm:text-xl font-semibold">Nice-to-haves</h3>
          <ul className="list-disc list-inside mb-4 text-base">
            <li>Experience with TypeScript</li>
            <li>Familiarity with GraphQL</li>
          </ul>

          <h3 className="text-lg sm:text-xl font-semibold">Perks & benefits</h3>
          <ul className="list-disc list-inside text-base">
            <li>Health insurance</li>
            <li>Remote work options</li>
            <li>Flexible working hours</li>
            <li>Professional development opportunities</li>
          </ul>
        </div>

        {/* Sidebar (Company + Actions) */}
        <div className="space-y-6">
          {/* Company info */}
          <div className="bg-white rounded-2xl shadow p-6">
            <h3 className="text-lg sm:text-xl font-semibold">Company</h3>
            <p className="text-base mt-2">
              JobGen is an AI-powered platform that helps people find jobs,
              improve resumes, and get personalized job recommendations.
            </p>
            <a href="#" className="text-teal-600 text-base block mt-2">
              Visit website
            </a>
            <a href="#" className="text-teal-600 text-base">
              Careers page
            </a>
          </div>

          {/* Apply button */}
          <button className="w-full bg-teal-500 text-white py-3 rounded-2xl shadow hover:bg-teal-600 transition text-base">
            Apply Now
          </button>

          {/* Save & Share box */}
          <div className="bg-white border border-gray-300 rounded-2xl shadow p-3 flex items-center justify-center gap-6">
            <button className="flex items-center gap-1 hover:text-teal-600 transition text-base">
              <Bookmark size={18} />
              <span>Save</span>
            </button>
            <button className="flex items-center gap-1 hover:text-teal-600 transition text-base">
              <Share2 size={18} />
              <span>Share</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default JobCard;
