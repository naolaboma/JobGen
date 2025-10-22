import {
  createApi,
  fetchBaseQuery,
  type FetchArgs,
  type FetchBaseQueryError,
  type FetchBaseQueryMeta,
} from "@reduxjs/toolkit/query/react";
import type { BaseQueryFn } from "@reduxjs/toolkit/query";
import type { ProfileData } from "../../../types/ProfileData";
import { getSession } from "next-auth/react";

type UploadResponse = {
  url?: string;
  fileName?: string;
  message?: string;
  [key: string]: unknown;
};

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
  if (typeof args === "string") {
    request = { url: args };
  } else {
    request = { ...args };
  }

  // Normalize headers and attach Authorization
  const headers = new Headers((request.headers as HeadersInit) || {});
  if (token) headers.set("Authorization", `Bearer ${token}`);
  request.headers = headers;

  const result = await rawBaseQuery(request, api, extraOptions);

  // Normalize error shape
  if (result && "error" in result && result.error) {
    const data = (result.error as FetchBaseQueryError).data as any;
    if (typeof data === "string") {
      try {
        (result.error as any).data = JSON.parse(data);
      } catch {
        (result.error as any).data = { message: String(data) };
      }
    } else if (!data || typeof data !== "object") {
      (result.error as any).data = { message: String(data) };
    }
  }

  return result;
};

export const profileApi = createApi({
  reducerPath: "profileApi",
  baseQuery,
  endpoints: (builder) => ({
    getProfile: builder.query<ProfileData, void>({
      query: () => "/users/profile",
    }),
    updateProfile: builder.mutation<ProfileData, Partial<ProfileData>>({
      query: (body) => ({
        url: "/users/profile",
        method: "PUT",
        body,
      }),
    }),
    deleteAccount: builder.mutation<{ message: string }, void>({
      query: () => ({
        url: "/users/account",
        method: "DELETE",
      }),
    }),

    // FILE ENDPOINTS
    getMyProfilePicture: builder.query<Blob, void>({
      query: () => ({
        url: "/files/profile-picture/me",
        responseHandler: (response: Response) => response.blob(),
      }),
    }),
    getProfilePicture: builder.query<Blob, string>({
      query: (id) => ({
        url: `/files/profile-picture/${id}`,
        responseHandler: (response: Response) => response.blob(),
      }),
    }),
    uploadDocument: builder.mutation<UploadResponse, FormData>({
      query: (formData) => ({
        url: "/files/upload/document",
        method: "POST",
        body: formData,
      }),
    }),
    uploadProfilePicture: builder.mutation<UploadResponse, FormData>({
      query: (formData) => ({
        url: "/files/upload/profile",
        method: "POST",
        body: formData,
      }),
    }),
    downloadFile: builder.query<Blob, string>({
      query: (id) => ({
        url: `/files/${id}`,
        responseHandler: (response: Response) => response.blob(),
      }),
    }),
    deleteFile: builder.mutation<{ message: string }, string>({
      query: (id) => ({
        url: `/files/${id}`,
        method: "DELETE",
      }),
    }),
  }),
});

export const {
  useGetProfileQuery,
  useUpdateProfileMutation,
  useDeleteAccountMutation,
  useGetMyProfilePictureQuery,
  useGetProfilePictureQuery,
  useUploadDocumentMutation,
  useUploadProfilePictureMutation,
  useDownloadFileQuery,
  useDeleteFileMutation,
} = profileApi;
