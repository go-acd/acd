package constants

import "errors"

var (
	// Response errors

	// ErrResponseUnknown is returned when the response status is not known.
	ErrResponseUnknown = errors.New("response returned an unknown status")
	// ErrResponseBadInput Bad input parameter. Error message should indicate
	// which one and why.
	ErrResponseBadInput = errors.New("response returned with status 400")
	// ErrResponseInvalidToken The client passed in the invalid Auth token.
	// Client should refresh the token and then try again.
	ErrResponseInvalidToken = errors.New("response returned with status 401")
	// ErrResponseForbidden Forbidden.
	ErrResponseForbidden = errors.New("response returned with status 403")
	// ErrResponseDuplicateExists Duplicate file exists.
	ErrResponseDuplicateExists = errors.New("response returned with status 409")
	// ErrResponseInternalServerError Servers are not working as expected. The
	// request is probably valid but needs to be requested again later.
	ErrResponseInternalServerError = errors.New("response returned with status 500")
	// ErrResponseUnavailable Service Unavailable.
	ErrResponseUnavailable = errors.New("response returned with status 503")
	// ErrJSONDecodingResponseBody is returned if there was an error decoding the
	// response body.
	ErrJSONDecodingResponseBody = errors.New("error while JSON-decoding the response body")
	// ErrReadingResponseBody is returned if ioutil.ReadAll() has failed.
	ErrReadingResponseBody = errors.New("error reading the entire response body")

	// Request errors

	// ErrCreatingHTTPRequest is returned if there was an error creating the HTTP
	// request.
	ErrCreatingHTTPRequest = errors.New("error creating HTTP request")
	// ErrDoingHTTPRequest is returned if there was an error doing the HTTP
	// request.
	ErrDoingHTTPRequest = errors.New("error doing the HTTP request")
	// ErrHTTPRequestTimeout is returned when an HTTP request has timed out.
	ErrHTTPRequestTimeout = errors.New("the request has timed out")

	// Downloading errors

	// ErrNodeDownload is returned if there was an error downloading the file.
	ErrNodeDownload = errors.New("error downloading the node")

	// Uploading errors

	// ErrFileExistsAndIsFolder is returned if attempting to upload a file but a
	// folder with the same path already exists.
	ErrFileExistsAndIsFolder = errors.New("the file exists and is a folder")
	// ErrFileExistsAndIsNotFolder is returned if attempting to create a folder
	// but a file with the same path already exists.
	ErrFileExistsAndIsNotFolder = errors.New("the file exists and is not a folder")
	// ErrFileExistsWithDifferentContents is returned if attempting to upload a
	// file but overwrite is disabled and the file already exists with different
	// contents.
	ErrFileExistsWithDifferentContents = errors.New("the files exists but is with different contents")
	// ErrWritingMetadata is returned if an error occurs whilst writing the metadata
	ErrWritingMetadata = errors.New("error writing the metadata")
	// ErrCreatingWriterFromFile is returned if an error happens when creating a writer from a file
	ErrCreatingWriterFromFile = errors.New("error creating a writer from a file")
	// ErrWritingFileContents is returned if an error happens when writing the file contents
	ErrWritingFileContents = errors.New("error writing the file contents")
	// ErrNoContentsToUpload is returned if the reader does not even have one byte.
	ErrNoContentsToUpload = errors.New("reader has not contents to upload")

	// JSON errors

	// ErrJSONEncoding is returned when an error occurs whilst encoding an object into JSON.
	ErrJSONEncoding = errors.New("error encoding an object to JSON")
	// ErrJSONDecoding is returned when an error occurs whilst decoding JSON to an object.
	ErrJSONDecoding = errors.New("error decoding JSON to an object")

	// GOB errors

	// ErrGOBEncoding is returned when an error occurs whilst encoding an object into GOB.
	ErrGOBEncoding = errors.New("error encoding an object to GOB")
	// ErrGOBDecoding is returned when an error occurs whilst decoding GOB to an object.
	ErrGOBDecoding = errors.New("error decoding GOB to an object")

	// Node errors

	// ErrNodeNotFound is returned when a node is not found.
	ErrNodeNotFound = errors.New("node not found")
	// ErrCannotCreateRootNode is returned if you attempt to create the root node
	ErrCannotCreateRootNode = errors.New("root node cannot be created")
	// ErrLoadingCache is returned when an error happens while loading from cacheFile
	ErrLoadingCache = errors.New("error loading from the cache file")
	// ErrMustFetchFresh is returned if the changes API requested a change.
	ErrMustFetchFresh = errors.New("must refresh the node tree")
	// ErrCannotCreateANodeUnderAFile is returned if you attempt to create a
	// folder/file under an existing file.
	ErrCannotCreateANodeUnderAFile = errors.New("cannot create a node under a file")

	// URL errors

	// ErrParsingURL is returned if an error occured whilst parsing a URL
	ErrParsingURL = errors.New("error parsing the URL")

	// File-related errors

	// ErrStatFile is returned if there was an error getting info about the file.
	ErrStatFile = errors.New("error stat() the file")
	// ErrOpenFile is returned if an error occurred while opening the file for reading
	ErrOpenFile = errors.New("error opening the file for reading")
	// ErrCreateFile is returned if an error is returned when trying to create a file
	ErrCreateFile = errors.New("error creating and/or truncating a file")
	// ErrCreateFolder is returned if an error occurred when trying to create a folder.
	ErrCreateFolder = errors.New("error creating a folder")
	// ErrFileExists is returned if the file already exists (on the server or locally).
	ErrFileExists = errors.New("the file already exists")
	// ErrFileNotFound is returned if no such file or directory.
	ErrFileNotFound = errors.New("no such file or directory")
	// ErrPathIsNotFolder is returned if the path is not a folder.
	ErrPathIsNotFolder = errors.New("path is not a folder")
	// ErrPathIsFolder is returned if the path is a folder.
	ErrPathIsFolder = errors.New("path is a folder")
	// ErrWrongPermissions is returned if the file has the wrong permissions.
	ErrWrongPermissions = errors.New("file has wrong permissions")
)
