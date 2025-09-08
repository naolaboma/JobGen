"use client";

import { FC, useState } from "react";
import { MapPin, Building2, Bookmark, Share2 } from "lucide-react";

interface JobCardProps {
  title: string;
  company: string;
  location: string;
  tags?: string[];
  salary: string;
  posted: string;
  deadline: string;
  applicants: number;
  description?: string;
  applyUrl?: string;
}

const cleanDescription = (html: string) => {
  return html
    .replace(/ð/g, "")
    .replace(/â€™/g, "'")
    .replace(/â€œ|â€�/g, '"')
    .replace(/â€/g, "-")
    .replace(/â€¢/g, "•")
    .replace(/â/g, "")
    .replace(/Ã©/g, "é")
    .replace(/Â/g, "")
    .replace(/°/g, "")
    .replace(/§»/g, "")
    .replace(/\\n/g, "\n");
};

const markdownBoldToHtml = (text: string) => {
  return text.replace(
    /\*\*(.*?)\*\*/g,
    (_, boldText) => `<strong>${boldText}</strong>`
  );
};

const convertHeadingParagraphsToList = (html: string) => {
  const parser = new DOMParser();
  const doc = parser.parseFromString(html, "text/html");

  const headings = Array.from(doc.querySelectorAll("h2, h3"));
  headings.forEach((heading) => {
    let next = heading.nextElementSibling;
    const paragraphs: Element[] = [];
    while (next && next.tagName.toLowerCase() === "p") {
      paragraphs.push(next);
      next = next.nextElementSibling;
    }
    if (paragraphs.length >= 2) {
      const ul = doc.createElement("ul");
      ul.className = "list-disc pl-6 space-y-2";
      paragraphs.forEach((p) => {
        const li = doc.createElement("li");
        li.innerHTML = p.innerHTML;
        ul.appendChild(li);
        p.remove();
      });
      heading.parentNode?.insertBefore(ul, next);
    }
  });
  return doc.body.innerHTML;
};

const JobCard: FC<JobCardProps> = ({
  title,
  company,
  location,
  tags = [],
  salary,
  posted,
  deadline,
  applicants,
  description,
  applyUrl,
}) => {
  const [expanded, setExpanded] = useState(false);

  function parseDescription(html?: string) {
    if (!html)
      return <p className="text-gray-600 italic">No description available.</p>;

    let cleanedHtml = cleanDescription(html);
    cleanedHtml = markdownBoldToHtml(cleanedHtml);
    cleanedHtml = convertHeadingParagraphsToList(cleanedHtml);

    const parts = cleanedHtml.split(/<hr\s*\/?>/i);

    return parts.map((part, i) => {
      if (typeof window === "undefined") return null;
      const parser = new DOMParser();
      const doc = parser.parseFromString(part, "text/html");
      const headingEl = doc.querySelector("h2, h3");
      const heading = headingEl?.textContent || `Section ${i + 1}`;
      if (headingEl) headingEl.remove();
      const contentHtml = doc.body.innerHTML;

      return (
        <section key={i} className="mb-8 last:mb-0">
          <h2 className="text-xl sm:text-3xl font-bold mb-4 text-teal-700 border-l-4 border-teal-500 pl-3">
            {heading}
          </h2>
          <div
            className="prose prose-teal max-w-none text-xs leading-relaxed prose-ul:list-disc prose-ul:pl-6 prose-li:mb-2 font-sans"
            dangerouslySetInnerHTML={{ __html: contentHtml }}
          />
        </section>
      );
    });
  }

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col items-center p-4 sm:p-6 text-black font-inter">
      {/* Header Row */}
      <div className="w-full max-w-5xl grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="p-6 md:col-span-2">
          <h1 className="text-3xl font-extrabold tracking-tight">{title}</h1>
          <div className="flex items-center gap-2 mt-3 text-lg text-gray-700">
            <Building2 size={20} />
            <span>{company}</span>
          </div>
          <div className="flex items-center gap-2 mt-1 text-lg text-gray-700">
            <MapPin size={20} />
            <span>{location}</span>
          </div>
          <div className="flex flex-wrap gap-2 mt-4">
            {tags.map((tag, i) => (
              <span
                key={i}
                className="px-3 py-1 text-sm font-medium bg-teal-100 text-teal-800 rounded-full"
              >
                {tag}
              </span>
            ))}
          </div>
        </div>

        <div className="bg-white rounded-2xl shadow px-6 flex flex-col gap-3">
          <h3 className="text-xl font-semibold text-teal-700 ">Job overview</h3>
          <p className="text-lg font-bold">{salary}</p>
          <p className="text-base text-gray-600">Full-time</p>
          <p className="text-base text-gray-600">Posted {posted}</p>
          <p className="text-base text-gray-600">Deadline: {deadline}</p>
          <p className="text-base text-gray-600">{applicants} applicants</p>
        </div>
      </div>

      {/* Body Row */}
      <div className="w-full max-w-5xl mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="col-span-2 bg-white rounded-lg border border-gray-300 p-6 relative">
          <h2 className="text-2xl font-bold mb-6 border-b border-gray-200 pb-3">
            Job Description
          </h2>

          <div
            className={`prose prose-teal max-w-none text-xs leading-relaxed prose-ul:list-disc prose-ul:pl-6 prose-li:mb-2 font-sans transition-all duration-300 ease-in-out ${
              expanded ? "max-h-[1000px]" : "max-h-[240px]"
            }`}
            style={{ overflow: "hidden", scrollbarWidth: "none", msOverflowStyle: "none" }}
          >
            {parseDescription(description)}
          </div>

          <button
            onClick={() => setExpanded((prev) => !prev)}
            className="mt-4 text-teal-600 hover:underline font-semibold text-xs"
            aria-expanded={expanded}
          >
            {expanded ? "Show Less" : "Read More"}
          </button>
        </div>

        <div className="space-y-8">
          <div className="bg-white rounded-2xl shadow p-6">
            <h3 className="text-xl font-semibold text-teal-700">Company</h3>
            <p className="mt-3 text-gray-700">{company}</p>
            <a
              href="#"
              className="text-teal-600 hover:underline text-base block mt-4"
            >
              Visit website
            </a>
            <a
              href="#"
              className="text-teal-600 hover:underline text-base mt-2 block"
            >
              Careers page
            </a>
          </div>

          {applyUrl && (
            <a
              href={applyUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="block w-full bg-teal-500 text-white py-3 rounded-2xl shadow hover:bg-teal-600 transition text-lg font-semibold text-center"
            >
              Apply Now
            </a>
          )}

          <div className="bg-white border border-gray-300 rounded-2xl shadow p-3 flex items-center justify-center gap-6">
            <button className="flex items-center gap-2 hover:text-teal-600 transition text-base font-medium">
              <Bookmark size={20} />
              <span>Save</span>
            </button>
            <button className="flex items-center gap-2 hover:text-teal-600 transition text-base font-medium">
              <Share2 size={20} />
              <span>Share</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default JobCard;
