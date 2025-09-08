"use client";

import React, { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useFetchJobsQuery } from "@/lib/redux/api/JobApi";
import JobSummaryCard from "@/app/components/JobSummaryCard";
import SearchBar from "@/app/components/SearchBar";
import Filters from "@/app/components/Filters";
import Link from "next/link";
import { setFilters, setPage, setSort } from "@/lib/redux/slices/jobSlice";
import type { RootState } from "@/store/store";
import Navbar from "@/app/components/header";

export default function Page() {
  const dispatch = useDispatch();
  const filters = useSelector((state: RootState) => state.job);

  const itemsPerPage = 3;
  const filtersWithLimit = { ...filters, limit: itemsPerPage };

  const { data, isLoading, error } = useFetchJobsQuery(filtersWithLimit);
  const jobs = data?.data?.items ?? [];
  const totalPages = data?.data?.total_pages ?? 1;

  const [showFilters, setShowFilters] = useState(false);

  const handleSearch = (query: string) => {
    dispatch(setFilters({ query }));
  };

  const handleFilterChange = (newFilters: {
    skills: string;
    location: string;
    sponsorship?: boolean;
  }) => {
    dispatch(setFilters(newFilters));
  };

  const handleSortChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const { name, value } = e.target;

    if (name === "sort_order") {
      dispatch(setSort({ sort_order: value as "asc" | "desc", sort_by: filters.sort_by ?? "posted_at" }));
    } else if (name === "sort_by") {
      dispatch(setSort({ sort_by: value, sort_order: filters.sort_order ?? "desc" }));
    }
  };

  const handlePageChange = (newPage: number) => {
    dispatch(setPage(newPage));
  };

  if (isLoading) {
    return <div className="text-center py-20 text-lg font-medium">Loading jobs...</div>;
  }

  if (error) {
    const errorMessage = "error" in (error as any) ? (error as any).error : "Unknown error occurred";
    return (
      <div className="text-center py-20 text-red-600 text-lg">
        Error: {errorMessage}
      </div>
    );
  }

  return (
    <>
      <Navbar />
      <div className="bg-gray-50 min-h-screen text-gray-900">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          {/* Header */}
          <div className="pt-8 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div>
              <h1 className="text-3xl font-bold text-gray-800">Find Your Next Job</h1>
              <p className="text-sm text-gray-600 mt-1">Curated opportunities that match your skills.</p>
            </div>
          </div>

          {/* Filters and chatbot */}
          <div className="mt-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
            <button
              onClick={() => setShowFilters(!showFilters)}
              className="px-5 py-2 rounded-lg bg-[#7BBFB3] text-white font-medium shadow-sm hover:brightness-105 active:brightness-95"
            >
              {showFilters ? "Hide Filters" : "Show Filters"}
            </button>
            <Link href="/chatbot" className="text-sm font-medium text-[#007B83] hover:underline">
              Refine with Chatbot
            </Link>
          </div>

          {/* Filters Panel */}
          {showFilters && (
            <div className="mt-4 rounded-xl border bg-white p-4 shadow-md">
              <div className="flex flex-col gap-4">
                <SearchBar onSearch={handleSearch} initialQuery={filters.query ?? ""} />
                <Filters
                  onFilterChange={handleFilterChange}
                  initialSkills={filters.skills ?? ""}
                  initialLocation={filters.location ?? ""}
                  initialSponsorship={filters.sponsorship}
                />
                <div className="flex flex-col sm:flex-row gap-3 pt-1">
                  <select
                    name="sort_by"
                    value={filters.sort_by ?? "posted_at"}
                    onChange={handleSortChange}
                    className="p-2 border rounded-md bg-white text-sm text-gray-800"
                  >
                    <option value="posted_at">Posted Date</option>
                    <option value="title">Title</option>
                  </select>
                  <select
                    name="sort_order"
                    value={filters.sort_order ?? "desc"}
                    onChange={handleSortChange}
                    className="p-2 border rounded-md bg-white text-sm text-gray-800"
                  >
                    <option value="desc">Descending</option>
                    <option value="asc">Ascending</option>
                  </select>
                </div>
              </div>
            </div>
          )}

          {/* Job listings */}
          <div className="py-6">
            {jobs.length === 0 ? (
              <div className="text-center py-20 text-lg text-gray-500">
                No jobs found. Try adjusting your filters.
              </div>
            ) : (
              <div className="space-y-4">
                {jobs.slice(0, itemsPerPage).map((job: any) => {
                  const id = job.id ?? job.job_id ?? job._id ?? job.uuid;
                  const title = job.title ?? job.job_title ?? "Untitled role";
                  const location = job.location ?? job.work_location ?? "Remote";
                  const jobType =
                    job.job_type ??
                    job.type ??
                    job.employment_type ??
                    job.tags?.find((t: string) =>
                      /full|part|contract|intern/i.test(t)
                    ) ??
                    "N/A";
                  const company = job.company ?? job.company_name ?? job.organization;

                  return (
                    <JobSummaryCard
                      key={String(id)}
                      id={String(id)}
                      title={title}
                      location={location}
                      jobType={jobType}
                      company={company}
                    />
                  );
                })}
              </div>
            )}

            {/* Pagination */}
            <div className="mt-10 flex justify-center">
              <PaginationSimple
                currentPage={filters.page ?? 1}
                totalPages={totalPages}
                onPageChange={handlePageChange}
              />
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

/* ---------------- Pagination with only Prev / Next ---------------- */
type PaginationProps = {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
};

function PaginationSimple({ currentPage, totalPages, onPageChange }: PaginationProps) {
  return (
    <div className="flex items-center gap-6">
      <button
        disabled={currentPage <= 1}
        onClick={() => onPageChange(currentPage - 1)}
        className="px-4 py-2 border rounded-md text-sm font-medium bg-[#7BBFB3] text-gray-800 hover:bg-gray-100 disabled:opacity-50"
      >
        ← Prev
      </button>

      <span className="text-sm text-gray-600">
        Page {currentPage} of {totalPages}
      </span>

      <button
        disabled={currentPage >= totalPages}
        onClick={() => onPageChange(currentPage + 1)}
        className="px-4 py-2 border rounded-md text-sm font-medium bg-[#7BBFB3] text-gray-800 hover:bg-gray-100 disabled:opacity-50"
      >
        Next →
      </button>
    </div>
  );
}
