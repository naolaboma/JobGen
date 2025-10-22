import {
  createApi,
  fetchBaseQuery,
  type FetchArgs,
  type FetchBaseQueryError,
  type FetchBaseQueryMeta,
} from "@reduxjs/toolkit/query/react";
import type { BaseQueryFn } from "@reduxjs/toolkit/query";
import { getSession } from "next-auth/react";

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

const rawBaseQuery = fetchBaseQuery({
  baseUrl: `${
    process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"
  }/api/v1`,
});

const baseQuery: BaseQueryFn<
  string | FetchArgs,
  unknown,
  FetchBaseQueryError,
  object,
  FetchBaseQueryMeta
> = async (args, api, extraOptions) => {
  const session = await getSession();
  const token = (session as any)?.accessToken as string | undefined;

  let request: FetchArgs;
  if (typeof args === "string") request = { url: args };
  else request = { ...args };

  const headers = new Headers((request.headers as HeadersInit) || {});
  if (token) headers.set("Authorization", `Bearer ${token}`);
  request.headers = headers;

  return rawBaseQuery(request, api, extraOptions);
};

export const JobApi = createApi({
  reducerPath: "jobApi",
  baseQuery,
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
        return [];
      },
    }),
  }),
});

export const { useFetchJobsQuery, useFetchMatchedJobsQuery } = JobApi;

export default JobApi;
