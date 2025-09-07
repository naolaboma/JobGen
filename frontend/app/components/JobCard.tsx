import React from "react";

const JobCard = ({
  title,
  description,
  company_name,
  location,
  apply_url,
  percentage,
}: JobProps) => {
  //   let colorClass = '';
  //   if (percentage <= 40) {
  //     colorClass = 'bg-red-200 text-red-800';
  //   } else if (percentage <= 79) {
  //     colorClass = 'bg-yellow-200 text-yellow-800';
  //   } else {
  //     colorClass = 'bg-green-200 text-green-800';
  //   }

  return (
    <div className="max-w-2xl mb-5 border rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
      <div className="p-4">
        <h2 className="text-xl font-semibold mb-2 text-gray-900">{title}</h2>
        <p className="text-gray-600 mb-2">{company_name}</p>
        <p className="text-gray-600 mb-2">{location}</p>
        <p className="text-gray-600 mb-4">{description.substring(0, 100)}...</p>
        <a
          href={apply_url}
          className="text-[#7BBFB3] hover:underline"
          target="_blank"
          rel="noopener noreferrer"
        >
          Apply Now
        </a>
        {/* {percentage !== undefined && (
          <div className={`inline-flex items-center justify-center px-3 py-1 text-sm font-bold rounded-full ${colorClass} mt-2`}>
            {percentage}%
          </div>
        )} */}
      </div>
    </div>
  );
};

export default JobCard;
