export enum JobType {
	FULL_TIME = "full-time",
	PART_TIME = "part-time",
	CONTRACT = "contract",
	INTERNSHIP = "internship",
	TEMPORARY = "temporary",
	REMOTE = "remote",
	HYBRID = "hybrid",
	FREELANCE = "freelance",
}

export enum UserRole {
	USER = "user",
	ADMIN = "admin",
}

export interface ProfileData {
	id: string;
	email: string;
	username: string;
	full_name: string;
	phone_number: string;
	location: string;
	skills: string[];
	experience_years: number;
	bio: string;
	profile_picture: string;
	job_type: JobType;
	preferred_country: string;
	city_region: string;
	role: UserRole;
	is_verified: boolean;
	is_active: boolean;
	last_login_at: string;
	created_at: string;
	updated_at: string;
}
