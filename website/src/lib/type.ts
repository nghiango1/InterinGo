export interface DocInfo {
	slug: string, // Used for recovered
	title?: string, // Navigator Title
	session?: string, // Session title
	index?: number, // Session position
}

export interface NavigationRecord {
	name: string, 
	docs: DocInfo[]
};
