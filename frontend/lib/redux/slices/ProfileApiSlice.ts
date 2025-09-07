import { createApi, fetchBaseQuery, type FetchArgs, type FetchBaseQueryError, type FetchBaseQueryMeta } from "@reduxjs/toolkit/query/react";
import type { BaseQueryFn } from "@reduxjs/toolkit/query";
import type { ProfileData } from "../../../types/ProfileData";

type UploadResponse = {
  url?: string;
  fileName?: string;
  message?: string;
  [key: string]: unknown;
};

const rawBaseQuery = fetchBaseQuery({
  baseUrl: "http://localhost:8080",
  prepareHeaders: (headers) => {
    return headers;
  },
});

const baseQuery: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError, object, FetchBaseQueryMeta> = async (args, api, extraOptions) => {
  const result = await rawBaseQuery(args, api, extraOptions);
  if (result && "error" in result && result.error) {
    const data = result.error.data;
    if (typeof data === "string") {
      try {
        result.error.data = JSON.parse(data);
      } catch {
        result.error.data = { message: String(data) };
      }
    } else if (!data || typeof data !== "object") {
      result.error.data = { message: String(data) };
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
