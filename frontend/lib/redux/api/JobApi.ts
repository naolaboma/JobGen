import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

interface ApiResponse {
  data: JobProps[] | string;
  error?: {
    code: string;
    details: string;
    message: string;
  };
  message: string;
  success: boolean;
}

export const JobApi = createApi({
  reducerPath: "jobApi",
  baseQuery: fetchBaseQuery({
    baseUrl: "http://localhost:8080/api/v1", // Update this URL
    prepareHeaders: (headers, { getState, endpoint }) => {
      // Apply Authorization header only for /jobs/matched endpoint
      if (endpoint === "fetchMatchedJobs") {
        const token = (getState() as any).auth?.token; // Adjust based on your auth slice
        if (token) {
          headers.set("Authorization", `Bearer ${token}`);
        }
      }
      return headers;
    },
  }),
  endpoints: (builder) => ({
    // Fetch all jobs (no authentication)
    fetchJobs: builder.query<JobProps[], JobQueryParams>({
      query: ({
        page = 1,
        limit = 2,
        query,
        skills,
        location,
        sponsorship,
        source,
        sort_by = "posted_at",
        sort_order = "desc",
      } = {}) => ({
        url: "/jobs",
        method: "GET",
        params: {
          page,
          limit,
          query,
          skills,
          location,
          sponsorship,
          source,
          sort_by,
          sort_order,
        },
      }),
    }),

    // Fetch matched jobs (requires authentication)
    fetchMatchedJobs: builder.query<
      JobProps[],
      { page?: number; limit?: number }
    >({
      query: (params) => ({
        url: "/jobs/matched",
        method: "GET",
        params: { page: params.page ?? 1, limit: params.limit ?? 2 },
      }),
      transformResponse: (response: ApiResponse) => {
        if (Array.isArray(response.data)) {
          return response.data;
        }
        return []; // Fallback if data is a string or unexpected format
      },
    }),
  }),
});

export const { useFetchJobsQuery, useFetchMatchedJobsQuery } = JobApi;

export default JobApi;
