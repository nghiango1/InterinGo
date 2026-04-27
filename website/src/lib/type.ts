export interface DocInfo {
	slug: string, // Used for recovered
	label?: string, // Navigator Title
	session?: string, // Session title
	index?: number, // Session position

	// Will be build via slug
	href?: string, 
}

export interface NavigationRecord {
	label: string, 
	items: DocInfo[]
};
