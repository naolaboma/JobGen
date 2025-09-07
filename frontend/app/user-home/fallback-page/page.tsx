"use client";

import { useDispatch, useSelector } from "react-redux";
import { useFetchJobsQuery } from "@/lib/redux/api/JobApi";
import JobCard from "@/app/components/JobCard";
import SearchBar from "@/app/components/SearchBar";
import Filters from "@/app/components/Filters";
import Pagination from "@/app/components/Pagination";
import Link from "next/link";
import { setFilters, setPage, setSort } from "@/lib/redux/slices/jobSlice";
import { RootState } from "@/store/store";

export default function FallbackPage() {
  const dispatch = useDispatch();
  const filters = useSelector((state: RootState) => state.job);

  const { data, isLoading, error } = useFetchJobsQuery(filters);

  const jobs = data?.data?.items ?? [];
  const totalPages = data?.data?.total_pages ?? 1;

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
    const sortOrder =
      name === "sort_order"
        ? (value as "asc" | "desc")
        : filters.sort_order ?? "desc";
    const sortBy = name === "sort_by" ? value : filters.sort_by ?? "posted_at";
    dispatch(setSort({ sort_by: sortBy, sort_order: sortOrder }));
  };

  const handlePageChange = (newPage: number) => {
    dispatch(setPage(newPage));
  };

  if (isLoading)
    return (
      <div className="text-center py-20 text-lg font-medium">
        Loading jobs...
      </div>
    );

  if (error) {
    const errorMessage =
      "error" in error ? error.error : "Unknown error occurred";
    return (
      <div className="text-center py-20 text-red-600 text-lg">
        Error: {errorMessage}
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Sticky Top Header */}
      <div className="sticky top-0 z-10 bg-white shadow-sm px-6 py-4 border-b">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Job Listings</h1>
            <p className="text-sm text-gray-600 mt-1">
              Set up your profile for personalized recommendations!{" "}
              <Link
                href="/profile-setup"
                className="text-[#7BBFB3] font-medium hover:underline"
              >
                Set Up Profile
              </Link>
            </p>
          </div>

          {/* Sorting Controls */}
          <div className="flex gap-2">
            <select
              name="sort_by"
              value={filters.sort_by ?? "posted_at"}
              onChange={handleSortChange}
              className="p-2 border rounded-md shadow-sm bg-white text-sm"
            >
              <option value="posted_at">Posted Date</option>
              <option value="title">Title</option>
            </select>
            <select
              name="sort_order"
              value={filters.sort_order ?? "desc"}
              onChange={handleSortChange}
              className="p-2 border rounded-md shadow-sm bg-white text-sm"
            >
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </select>
          </div>
        </div>

        {/* Search + Filters Row */}
        <div className="mt-4 flex flex-col md:flex-row md:items-center gap-4">
          <SearchBar
            onSearch={handleSearch}
            initialQuery={filters.query ?? ""}
          />
          <Filters
            onFilterChange={handleFilterChange}
            initialSkills={filters.skills ?? ""}
            initialLocation={filters.location ?? ""}
            initialSponsorship={filters.sponsorship}
          />
        </div>
      </div>

      {/* Job Listings */}
      <div className="p-6">
        {jobs.length === 0 ? (
          <div className="text-center py-20 text-gray-600 text-lg">
            No jobs found. Try adjusting filters.
          </div>
        ) : (
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {jobs.map((job) => (
              <JobCard
                key={job.id}
                {...job}
                percentage={Math.floor(Math.random() * 100)} // placeholder
              />
            ))}
          </div>
        )}

        {/* Pagination */}
        <div className="mt-10 flex justify-center">
          <Pagination
            currentPage={filters.page ?? 1}
            totalPages={totalPages}
            onPageChange={handlePageChange}
          />
        </div>
      </div>
    </div>
  );
}
