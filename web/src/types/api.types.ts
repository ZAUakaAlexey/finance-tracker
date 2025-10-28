export interface GeneralResponse<T> {
	data: T;
	message: string;
	errors: string[];
}

export interface NoPaginatedResponse<T> {
	data: {
		resource: T;
	};
	message: string | null;
	errors: Record<string, string[]>;
}

export interface IPaginatedResponse<T> {
	data: {
		resource: {
			items: T[];
			meta: PaginationMeta;
		};
	};
	message: string | null;
	errors: Record<string, string[]>;
}

export interface PaginationMeta {
	current_page: number;
	from: number;
	last_page: number;
	per_page: number;
	to: number;
	total: number;
}
