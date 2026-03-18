// This help with building index file, as we have nested directory
// that can cause different redirect behavior to static file hosting
// - Directory -> 'docs/index.html' is forced redirected to
// - No directory -> 'docs.html' is used
export const trailingSlash = 'always';
