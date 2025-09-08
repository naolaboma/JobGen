// app/components/JobSummaryCard.tsx
"use client";

import React, { FC } from "react";
import Link from "next/link";

interface JobSummaryProps {
  id: string;
  title: string;
  location?: string;
  jobType?: string;
  company?: string;
}

const JobSummaryCard: FC<JobSummaryProps> = ({
  id,
  title,
  location,
  jobType,
  company,
}) => {
  return (
    <Link href={`/job-details/${id}`} className="block">
      <div className="p-4 border rounded-lg hover:shadow-md transition-shadow duration-150 bg-white">
        <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
        {company && <p className="text-sm text-gray-600 mt-1">{company}</p>}
        <div className="mt-3 flex items-center gap-3 text-sm text-gray-700">
          <span>{location ?? "Remote"}</span>
          <span className="px-2 py-1 text-xs font-medium bg-teal-100 text-teal-800 rounded-full">
            {jobType ?? "N/A"}
          </span>
        </div>
      </div>
    </Link>
  );
};

export default JobSummaryCard;
