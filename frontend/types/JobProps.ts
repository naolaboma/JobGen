interface JobProps {
    id: string;
    title: string;
    company_name: string;
    location: string;
    description: string;
    apply_url: string;
    source: string;
    posted_at: string;
    is_sponsorship_available: boolean;
    extracted_skills: string[];
    created_at: string;
    updated_at: string;
    remote_ok_id: string;
    salary: string;
    tags: string[];
    percentage?: number; // Optional percentage field
}